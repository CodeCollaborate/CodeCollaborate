package dbfs

import (
	"fmt"
	"strconv"
	"time"

	"github.com/CodeCollaborate/Server/modules/patching"
	"github.com/CodeCollaborate/Server/utils"
)

// TODO (wongb): Move this to config

// MinBufferLength specifies the minimum number of patches left in the database after scrunching
var MinBufferLength = 50

// MaxBufferLength specifies the maximum number of patches left in the database before scrunching
var MaxBufferLength = MinBufferLength * 10

// ScrunchingExpiryLength specifies the maximum time we will allow scrunching to occur before we
// consider it failed and could retry
var ScrunchingExpiryLength = uint32((5 * time.Minute).Seconds())

// ScrunchFile scrunches all but the last minBufferLength items into the file on disk
// It then removes the changes from Couchbase
func (di *DatabaseImpl) ScrunchFile(meta FileMeta) error {
	utils.LogDebug("Scrunching: Starting", utils.LogFields{
		"FileID": meta.FileID,
	})

	start := time.Now()

	changes, baseFile, err := di.getForScrunching(meta, MinBufferLength)
	if err != nil {
		return fmt.Errorf("Scrunching - Failed to retrieve patches and file for scrunching: %v", err)
	}
	if len(changes) == 0 {
		// nothing to do, return silently
		return nil
	}

	utils.LogDebug("Scrunching: Got patches", utils.LogFields{
		"FileID":     meta.FileID,
		"Changes":    changes,
		"NumChanges": len(changes),
	})

	result, err := patching.PatchTextFromString(string(baseFile), changes)
	if err != nil {
		return fmt.Errorf("Scrunching - Failed to scrunch file: %v", err)
	}

	utils.LogDebug("Scrunching: Patched text", utils.LogFields{
		"FileID":     meta.FileID,
		"NumChanges": len(changes),
	})

	if err := di.FileWriteToSwap(meta, []byte(result)); err != nil {
		return fmt.Errorf("Scrunching - Failed to write to swap file: %v", err)
	}

	utils.LogDebug("Scrunching: Wrote to swap file", utils.LogFields{
		"FileID": meta.FileID,
	})

	if err := di.deleteForScrunching(meta, len(changes)); err != nil {
		return fmt.Errorf("Scrunching - Failed to removed scrunched changes: %v", err)
	}

	elapsed := time.Since(start)

	utils.LogDebug("Scrunching: Done", utils.LogFields{
		"FileID":         meta.FileID,
		"Execution Time": elapsed.Seconds(),
	})

	return nil
}

// GetForScrunching gets all but the remainder entries for a file and creates a temp swp file
// returns the changes for scrunching, the swap file contents, and any errors
func (di *DatabaseImpl) getForScrunching(fileMeta FileMeta, remainder int) ([]string, []byte, error) {
	cb, err := di.openCouchBase()
	if err != nil {
		return []string{}, []byte{}, err
	}

	fileData, err := cb.couchbaseDocumentStore.GetFileData(fileMeta.FileID)
	if err != nil {
		return []string{}, []byte{}, ErrResourceNotFound
	}

	changes := fileData.Changes

	if len(changes)-(remainder+1) < 0 {
		return []string{}, []byte{}, ErrNoDbChange
	}

	err = di.scrunchingAddLock(strconv.FormatInt(fileMeta.FileID, 10))
	if err != nil {
		// If it finds a document, we're already scrunching and it will fail (because insert, not upsert).
		// Unfortunately, couchbase doesn't have any better way to tell if a key exists,
		// so we can't do any better than doing this and just eating the error *grumble*
		utils.LogDebug("Scrunching: Scrunching (probably) already in progress, aborting", utils.LogFields{
			"FileID":            fileMeta.FileID,
			"Couchbase Message": err,
		})
		return []string{}, []byte{}, nil
	}

	swp, err := di.makeSwp(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)

	return changes[0 : len(changes)-remainder], swp, err
}

// DeleteForScrunching deletes `num` elements from the front of `changes` for file with `fileID` and deletes the
// swp file
func (di *DatabaseImpl) deleteForScrunching(fileMeta FileMeta, num int) error {
	cb, err := di.openCouchBase()
	if err != nil {
		return err
	}
	// NOTE: the test for this in multi_test.go walks through this logic, ensuring pull works throughout
	//		 therefore, any changes made here need to be reflected there as well

	fileKey := strconv.FormatInt(fileMeta.FileID, 10)

	// turn on writing to TempChanges
	builder := cb.bucket.MutateIn(fileKey, 0, 0)
	builder = builder.Upsert("TempChanges", []string{}, false)
	builder = builder.Upsert("UseTemp", true, false)
	_, err = builder.Execute()
	if err != nil {
		return err
	}

	// get changes in normal changes
	fileData, err := cb.couchbaseDocumentStore.GetFileData(fileMeta.FileID)
	if err != nil {
		return err
	}

	changes := fileData.Changes

	if len(changes) <= num {
		// somehow something scrunched this file at the same time
		utils.LogWarn("Scrunching: possible concurrent scrunching of the same file. "+
			"Maybe `ScrunchingExpiryLength` isn't long enough?", utils.LogFields{
			"FileID": fileMeta.FileID,
		})
		return nil
	}

	// turn off writing to TempChanges & reset normal changes
	builder = cb.bucket.MutateIn(fileKey, 0, 0)
	builder = builder.Upsert("RemainingChanges", changes[num:], false)
	builder = builder.Upsert("Changes", []string{}, false)
	builder = builder.Upsert("UseTemp", false, false)
	builder = builder.Upsert("PullSwp", true, false)
	_, err = builder.Execute()
	if err != nil {
		return err
	}

	// get changes in temp changes
	fileData, err = cb.couchbaseDocumentStore.GetFileData(fileMeta.FileID)
	if err != nil {
		return err
	}

	tempChanges := fileData.TempChanges

	err = di.swapSwp(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)
	if err != nil {
		utils.LogError("error replacing file with scrunched swap file", err, utils.LogFields{
			"Filename":    fileMeta.Filename,
			"ProjectID":   fileMeta.ProjectID,
			"File relath": fileMeta.RelativePath,
		})
		// undo everything
		builder = cb.bucket.MutateIn(fileKey, 0, 0)
		builder = builder.ArrayPrependMulti("Changes", append(changes, tempChanges...), false)
		builder = builder.Upsert("RemainingChanges", []string{}, false)
		builder = builder.Upsert("TempChanges", []string{}, false)
		builder = builder.Upsert("PullSwp", false, false)
		builder.Execute()
		di.deleteSwp(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)
		return err
	}

	// prepend changes and reset temporarily stored changes
	builder = cb.bucket.MutateIn(fileKey, 0, 0)
	builder = builder.ArrayPrependMulti("Changes", append(changes[num:], tempChanges...), false)
	builder = builder.Upsert("RemainingChanges", []string{}, false)
	builder = builder.Upsert("TempChanges", []string{}, false)
	builder = builder.Upsert("PullSwp", false, false)
	_, err = builder.Execute()

	err = di.deleteSwp(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)
	if err != nil {
		utils.LogError("error deleting swap file", err, utils.LogFields{
			"Filename":    fileMeta.Filename,
			"ProjectID":   fileMeta.ProjectID,
			"File relath": fileMeta.RelativePath,
		})
	}

	err = di.scrunchingRemoveLock(fileKey)
	if err != nil {
		utils.LogDebug("Scrunching: took longer than allocated scrunching time", utils.LogFields{
			"FileID":       fileMeta.FileID,
			"Allowed Time": ScrunchingExpiryLength,
		})
	}

	return err
}

// scrunchingAddLock hints to the server that the file with key `key` is currently being scrunched
func (di *DatabaseImpl) scrunchingAddLock(key string) error {
	cb, err := di.openCouchBase()
	if err != nil {
		return err
	}

	empty := true
	_, err = cb.scrunchingLocksBucket.Insert(key, &empty, ScrunchingExpiryLength)
	return err
}

// scrunchingRemoveLock removes the scrunching lock on the file with key `key` so that it can be scrunched later
func (di *DatabaseImpl) scrunchingRemoveLock(key string) error {
	cb, err := di.openCouchBase()
	if err != nil {
		return err
	}

	_, err = cb.scrunchingLocksBucket.Remove(key, 0)
	return err
}

// PullFile pulls the changes and the file bytes from the databases
func (di *DatabaseImpl) PullFile(meta FileMeta) (*[]byte, []string, error) {
	cb, err := di.openCouchBase()
	if err != nil {
		return new([]byte), []string{}, err
	}

	fileData, err := cb.couchbaseDocumentStore.GetFileData(meta.FileID)
	if err != nil {
		return new([]byte), []string{}, err
	}

	changes := fileData.AggregatedChanges()
	if fileData.PullSwp {
		bytes, err := di.swapRead(meta.RelativePath, meta.Filename, meta.ProjectID)
		if err != nil {
			return new([]byte), []string{}, err
		}
		return bytes, changes, nil
	}

	bytes, err := di.FileRead(meta.RelativePath, meta.Filename, meta.ProjectID)
	if err != nil {
		return new([]byte), []string{}, err
	}
	return bytes, changes, err
}
