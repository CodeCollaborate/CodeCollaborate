package dbfs

import (
	"strconv"

	"github.com/CodeCollaborate/Server/utils"
)

// GetForScrunching gets all but the remainder entries for a file and creates a temp swp file
// returns the changes for scrunching, the swap file contents, and any errors
func (di *DatabaseImpl) GetForScrunching(fileID int64, remainder int) ([]string, []byte, error) {
	cb, err := di.openCouchBase()
	if err != nil {
		return []string{}, []byte{}, err
	}

	fileMeta, err := di.MySQLFileGetInfo(fileID)
	if err != nil {
		return []string{}, []byte{}, err
	}

	frag, err := cb.bucket.LookupIn(strconv.FormatInt(fileID, 10)).Get("changes").Execute()
	if err != nil {
		return []string{}, []byte{}, ErrResourceNotFound
	}

	changes := []string{}
	err = frag.Content("changes", &changes)
	if err != nil {
		return []string{}, []byte{}, ErrResourceNotFound
	}

	if len(changes)-remainder+1 > 0 {
		return []string{}, []byte{}, ErrNoDbChange
	}

	swp, err := di.makeSwp(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)

	return changes[0 : len(changes)-remainder], swp, err
}

// DeleteForScrunching deletes `num` elements from the front of `changes` for file with `fileID` and deletes the
// swp file
func (di *DatabaseImpl) DeleteForScrunching(fileID int64, num int) error {
	cb, err := di.openCouchBase()
	if err != nil {
		return err
	}

	fileMeta, err := di.MySQLFileGetInfo(fileID)
	if err != nil {
		return err
	}

	key := strconv.FormatInt(fileID, 10)

	// turn on writing to TempChanges
	builder := cb.bucket.MutateIn(key, 0, 0)
	builder = builder.Upsert("tempchanges", []string{}, false)
	builder = builder.Upsert("usetemp", true, false)
	_, err = builder.Execute()
	if err != nil {
		return err
	}

	// get changes in normal changes
	frag, err := cb.bucket.LookupIn(key).Get("changes").Execute()
	if err != nil {
		return err
	}

	changes := []string{}
	err = frag.Content("changes", &changes)
	if err != nil {
		return ErrResourceNotFound
	}

	// turn off writing to TempChanges & reset normal changes
	builder = cb.bucket.MutateIn(key, 0, 0)
	builder = builder.Upsert("remaining_changes", changes[num:], false)
	builder = builder.Upsert("changes", []string{}, false)
	builder = builder.Upsert("usetemp", false, false)
	builder = builder.Upsert("pullswp", true, false)
	_, err = builder.Execute()
	if err != nil {
		return err
	}

	// get changes in TempChanges
	frag, err = cb.bucket.LookupIn(key).Get("tempchanges").Execute()
	if err != nil {
		return err
	}

	tempChanges := []string{}
	err = frag.Content("tempchanges", &tempChanges)
	if err != nil {
		return ErrResourceNotFound
	}

	err = di.swapSwp(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)
	if err != nil {
		utils.LogError("error replacing file with scrunched swap file", err, utils.LogFields{
			"Filename":    fileMeta.Filename,
			"ProjectID":   fileMeta.ProjectID,
			"File relath": fileMeta.RelativePath,
		})
		// undo everything
		builder = cb.bucket.MutateIn(key, 0, 0)
		builder = builder.ArrayPrependMulti("changes", append(changes, tempChanges...), false)
		builder = builder.Upsert("remaining_changes", []string{}, false)
		builder = builder.Upsert("tempchanges", []string{}, false)
		builder = builder.Upsert("pullswp", false, false)
		builder.Execute()
		di.deleteSwp(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)
		return err
	}
	err = di.deleteSwp(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)
	if err != nil {
		utils.LogError("error deleting swap file", err, utils.LogFields{
			"Filename":    fileMeta.Filename,
			"ProjectID":   fileMeta.ProjectID,
			"File relath": fileMeta.RelativePath,
		})
	}

	// prepend changes and reset temporarily stored changes
	builder = cb.bucket.MutateIn(key, 0, 0)
	builder = builder.ArrayPrependMulti("changes", append(changes[num:], tempChanges...), false)
	builder = builder.Upsert("remaining_changes", []string{}, false)
	builder = builder.Upsert("tempchanges", []string{}, false)
	builder = builder.Upsert("pullswp", false, false)
	_, err = builder.Execute()

	return err
}

// PullFile pulls the changes and the file bytes from the databases
func (di *DatabaseImpl) PullFile(meta FileMeta) (*[]byte, []string, error) {
	cb, err := di.openCouchBase()
	if err != nil {
		return new([]byte), []string{}, err
	}

	file := cbFile{}
	_, err = cb.bucket.Get(strconv.FormatInt(meta.FileID, 10), &file)
	if err != nil {
		return new([]byte), []string{}, err
	}
	var changes []string

	if file.PullSwp {
		changes = append(file.RemainingChanges, file.TempChanges...)
		changes = append(changes, file.Changes...)

		bytes, err := di.swapRead(meta.RelativePath, meta.Filename, meta.ProjectID)
		if err != nil {
			return new([]byte), []string{}, err
		}
		return bytes, changes, nil
	} else if file.UseTemp {
		changes = append(file.Changes, file.TempChanges...)
	} else {
		changes = file.Changes
	}

	bytes, err := di.FileRead(meta.RelativePath, meta.Filename, meta.ProjectID)
	if err != nil {
		return new([]byte), []string{}, err
	}
	return bytes, changes, err
}
