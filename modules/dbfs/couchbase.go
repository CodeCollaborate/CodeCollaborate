package dbfs

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/patching"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/couchbase/gocb"
)

type couchbaseConn struct {
	config                config.ConnCfg
	bucket                *gocb.Bucket
	scrunchingLocksBucket *gocb.Bucket
}

type cbFile struct {
	FileID           int64    `json:"-"`
	Version          int64    `json:"version"`
	Changes          []string `json:"changes"`
	TempChanges      []string `json:"tempchanges"`
	RemainingChanges []string `json:"remaining_changes"`
	UseTemp          bool     `json:"usetemp"`
	PullSwp          bool     `json:"pullswp"`
}

func (di *DatabaseImpl) openCouchBase() (*couchbaseConn, error) {
	if di.couchbaseDB != nil && di.couchbaseDB.bucket != nil {
		return di.couchbaseDB, nil
	}

	if di.couchbaseDB == nil || di.couchbaseDB.config == (config.ConnCfg{}) {
		di.couchbaseDB = new(couchbaseConn)
		configMap := config.GetConfig()
		di.couchbaseDB.config = configMap.ConnectionConfig["Couchbase"]
	}

	var documentsCluster *gocb.Cluster
	var err error

	if strings.HasPrefix(di.couchbaseDB.config.Host, "couchbase://") {
		documentsCluster, err = gocb.Connect(di.couchbaseDB.config.Host + ":" + strconv.Itoa(int(di.couchbaseDB.config.Port)))
	} else {
		documentsCluster, err = gocb.Connect("couchbase://" + di.couchbaseDB.config.Host + ":" + strconv.Itoa(int(di.couchbaseDB.config.Port)))
	}

	if err != nil {
		utils.LogError("Couchbase: could not connect to couchbase", err, utils.LogFields{
			"Host": di.couchbaseDB.config.Host,
		})
		return di.couchbaseDB, err
	}

	if di.couchbaseDB.config.Schema == "" {
		di.couchbaseDB.config.Schema = "documents"
	}

	schemaBucket, err := documentsCluster.OpenBucket(di.couchbaseDB.config.Schema, di.couchbaseDB.config.Password)
	if err != nil {
		utils.LogError("Couchbase: could not open bucket", err, utils.LogFields{
			"Host":   di.couchbaseDB.config.Host,
			"Bucket": di.couchbaseDB.config.Schema,
		})
		return di.couchbaseDB, err
	}
	di.couchbaseDB.bucket = schemaBucket

	// need to use 2nd bucket b/c couchbase has document expiry, not key expiry
	locksBucketName := di.couchbaseDB.config.Schema + "_scrunching_locks"
	slBucket, err := documentsCluster.OpenBucket(locksBucketName, di.couchbaseDB.config.Password)
	if err != nil {
		utils.LogError("Couchbase: could not open bucket", err, utils.LogFields{
			"Host":   di.couchbaseDB.config.Host,
			"Bucket": locksBucketName,
		})
		return di.couchbaseDB, err
	}
	di.couchbaseDB.scrunchingLocksBucket = slBucket

	return di.couchbaseDB, nil
}

// CloseCouchbase closes the CouchBase db connection
// YOU PROBABLY DON'T NEED TO RUN THIS EVER
func (di *DatabaseImpl) CloseCouchbase() error {
	if di.couchbaseDB != nil && di.couchbaseDB.bucket != nil {
		di.couchbaseDB.bucket.Close()
		di.couchbaseDB = nil
	} else {
		return ErrDbNotInitialized
	}

	return nil
}

// CBInsertNewFile inserts a new document into couchbase with CBFile.FileID == fileID
func (di *DatabaseImpl) cbInsertNewFile(file cbFile) error {
	cb, err := di.openCouchBase()

	if err != nil {
		return err
	}

	_, err = cb.bucket.Insert(strconv.FormatInt(file.FileID, 10), file, 0)
	return err
}

// CBInsertNewFile inserts a new document with the given arguments
func (di *DatabaseImpl) CBInsertNewFile(fileID int64, version int64, changes []string) error {
	return di.cbInsertNewFile(cbFile{
		FileID:           fileID,
		Version:          version,
		Changes:          changes,
		UseTemp:          false,
		TempChanges:      []string{},
		PullSwp:          false,
		RemainingChanges: []string{},
	})
}

// CBDeleteFile deletes the document with FileID == fileID from couchbase
func (di *DatabaseImpl) CBDeleteFile(fileID int64) error {
	cb, err := di.openCouchBase()
	if err != nil {
		return err
	}
	_, err = cb.bucket.Remove(strconv.FormatInt(fileID, 10), 0)
	return err
}

// CBGetFileVersion returns the current version of the file for the given FileID
func (di *DatabaseImpl) CBGetFileVersion(fileID int64) (int64, error) {
	cb, err := di.openCouchBase()
	if err != nil {
		return -1, err
	}

	frag, err := cb.bucket.LookupIn(strconv.FormatInt(fileID, 10)).Get("version").Execute()
	if err != nil {
		return -1, err
	}

	var version int64
	err = frag.Content("version", &version)
	if err != nil {
		return -1, ErrResourceNotFound
	}

	return version, err
}

// CBAppendFileChange mutates the file document with the new change and sets the new version number
// Returns the new version number, the missing patches, and an error, if any.
func (di *DatabaseImpl) CBAppendFileChange(fileID int64, patches, prevChanges []string) ([]string, int64, []string, error) {
	cb, err := di.openCouchBase()
	if err != nil {
		return nil, -1, nil, err
	}
	key := strconv.FormatInt(fileID, 10)

	// optimistic locking operation
	// check the version is accurate and get the object's cas,
	// then use it in the MutateIn call to verify the document hasn't updated underneath us
	frag, err := cb.bucket.LookupIn(key).Get("version").Get("usetemp").Execute()
	if err != nil {
		return nil, -1, nil, ErrResourceNotFound
	}

	cas := frag.Cas()
	var version int64
	var useTemp bool
	err = frag.Content("version", &version)
	if err != nil {
		return nil, -1, nil, ErrNoData
	}
	err = frag.Content("usetemp", &useTemp)
	if err != nil {
		return nil, -1, nil, ErrNoData
	}

	minVersion := version
	if len(prevChanges) > 0 {
		startPatch, err := patching.NewPatchFromString(prevChanges[0])
		if err != nil {
			return nil, -1, nil, ErrInternalServerError
		}

		minVersion = startPatch.BaseVersion + 1 // The patch with minVersion would have generated version minVersion + 1
	}
	minStartIndex := int64(math.MaxInt64)
	transformedPatches := []string{}

	// Build patch, transform changes against newer changes.
	for _, changeStr := range patches {
		change, err := patching.NewPatchFromString(changeStr)
		if err != nil {
			return nil, -1, nil, errors.New("Failed to parse patch")
		}

		// check to make sure the patch is being applied to the most recent revision
		if change.BaseVersion > version {
			utils.LogError("BaseVersion too high", ErrVersionOutOfDate, nil)
			return nil, -1, nil, ErrVersionOutOfDate
		}

		// For every patch, calculate the patches that it does not have.
		utils.LogDebug("CHANGES VERSIONS", utils.LogFields{
			"Version":     version,
			"BaseVersion": change.BaseVersion,
			"Diff":        int(version - change.BaseVersion),
			"Len":         len(prevChanges),
			"ChangeStr":   changeStr,
			"PrevChanges": prevChanges,
			"minVersion":  minVersion,
		})

		//startIndex := len(prevChanges) - int(version-change.BaseVersion)
		//if startIndex < 0 {
		//	utils.LogError("StartIndex is negative", ErrVersionOutOfDate, nil)
		//	return nil, -1, nil, ErrVersionOutOfDate
		//}

		startIndex := int64(len(prevChanges) - 1)

		if change.BaseVersion == version {
			// If we are building on the server's base version, don't need to transform.
			startIndex = int64(len(prevChanges))
		} else if change.BaseVersion == minVersion {
			// If it's equal to the minVersion, we use the entire array
			startIndex = int64(0)
		} else if change.BaseVersion < minVersion {
			// if it's less than the minVersion, we've scrunched.
			utils.LogError("BaseVersion less than minVersion", ErrVersionOutOfDate, nil)
			return nil, -1, nil, ErrVersionOutOfDate
		}
		// Otherwise, find the right starting point
		startIndex = int64(len(prevChanges)) - (version - change.BaseVersion)
		for startIndex >= 0 && startIndex < int64(len(prevChanges)) {
			otherPatch, err := patching.NewPatchFromString(prevChanges[startIndex])

			if err != nil {
				return nil, -1, nil, ErrInternalServerError
			}

			if change.BaseVersion > otherPatch.BaseVersion {
				startIndex++ // go back to the actual base version
				break
			} else {
				startIndex--
			}
		}

		if startIndex < 0 {
			utils.LogError("BaseVersion too low", ErrVersionOutOfDate, nil)
			return nil, -1, nil, ErrVersionOutOfDate
		}

		if startIndex < minStartIndex {
			minStartIndex = startIndex
		}

		utils.LogDebug("FINISHED CHECKING", utils.LogFields{
			"Change":     changeStr,
			"StartIndex": startIndex,
			"Len":        len(prevChanges),
		})

		// Apply patches from the change's baseVersion onwards
		toApply := prevChanges[startIndex:]

		utils.LogDebug("TRANSFORMING", utils.LogFields{
			"PatchesToApply": toApply,
			"Change":         changeStr,
			"StartIndex":     startIndex,
			"Len":            len(prevChanges),
		})

		transformedPatch, err := change.TransformFromString(toApply) // rewrite change with transformed patch
		if err != nil {
			return nil, -1, nil, ErrInternalServerError // Could not parse one of the old patches - should never happen.
		}

		// Update the BaseVersion to be be the previous change
		//transformedPatch.BaseVersion++
		transformedPatches = append(transformedPatches, transformedPatch.String())
	}

	/*
		// THIS BLOCK OF CODE HAS BEEN DISABLED BECAUSE PATCH CONSOLIDATION DOES NOT WORK AS EXPECTED
		// For this to correctly work, a new OT-like algorithm will need to be implemented.

		var consolidatedPatch *patching.Patch
		for _, transformed := range transformedPatches {
			if consolidatedPatch == nil {
				consolidatedPatch = transformed
			} else {
				hoistedPatch := transformed.Transform([]*patching.Patch{consolidatedPatch.Undo()})
				newChanges := append(consolidatedPatch.Changes, hoistedPatch.Changes...)
				sort.Sort(newChanges)
				consolidatedPatch.Changes = newChanges
			}
		}

		// use the cas to make sure the document hasn't changed
		builder := cb.bucket.MutateIn(key, cas, 0)

		if !useTemp {
			builder.ArrayAppend("changes", consolidatedPatch.String(), false)
		} else {
			builder.ArrayAppend("tempchanges", consolidatedPatch.String(), false)
		}
	*/

	// use the cas to make sure the document hasn't changed
	builder := cb.bucket.MutateIn(key, cas, 0)

	if !useTemp {
		builder.ArrayAppendMulti("changes", transformedPatches, false)
	} else {
		builder.ArrayAppendMulti("tempchanges", transformedPatches, false)
	}

	builder = builder.Counter("version", 1, false)

	_, err = builder.Execute()
	if err != nil {
		return nil, -1, nil, err
	}

	return transformedPatches, version + 1, prevChanges[minStartIndex:], err
}
