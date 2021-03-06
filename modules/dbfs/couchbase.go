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
	"github.com/davecgh/go-spew/spew"
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
// Returns the new version number, the missing patches, the total count of patches tracked, and an error, if any.
func (di *DatabaseImpl) CBAppendFileChange(fileMeta FileMeta, patchStr string) (string, int64, []string, int, error) {
	cb, err := di.openCouchBase()
	if err != nil {
		return "", -1, nil, 0, err
	}

	// optimistic locking operation
	// check the version is accurate and get the object's cas,
	// then use it in the MutateIn call to verify the document hasn't updated underneath us
	prevChangeStrs, cas, version, useTemp, err := di.PullChanges(fileMeta)
	if err != nil {
		return "", -1, nil, 0, err
	}

	prevChanges, err := patching.GetPatches(prevChangeStrs)
	if err != nil {
		utils.LogError("Failed to parse previous changes into patch objects", err, utils.LogFields{
			"PrevChanges": prevChangeStrs,
		})
		return "", -1, nil, 0, err
	}

	if cas == uint64(0) {
		utils.LogWarn("Couchbase returned a CAS value of 0, optimistic locking is unavailable", utils.LogFields{
			"cas":  cas,
			"File": fileMeta,
		})
	}

	minVersion := version
	if len(prevChangeStrs) > 0 {
		startPatch, err := patching.NewPatchFromString(prevChangeStrs[0])
		if err != nil {
			utils.LogError("Failed to parse first patch", err, utils.LogFields{
				"PatchStr": prevChangeStrs[0],
			})
			return "", -1, nil, 0, ErrInternalServerError
		}

		// Allow transform-patches to start on the same base version as the head (after linearization, we have all the necessary patches)
		minVersion = startPatch.BaseVersion
	}
	minStartIndex := int64(math.MaxInt64)
	prevChangesCopy := make([]string, len(prevChangeStrs))
	copy(prevChangesCopy, prevChangeStrs)

	// Build patch, transform changes against newer changes.
	change, err := patching.NewPatchFromString(patchStr)
	if err != nil {
		return "", -1, nil, 0, errors.New("Failed to parse patch")
	}

	// For every patch, calculate the patches that it does not have.
	utils.LogDebug("CHANGES VERSIONS", utils.LogFields{
		"Version":     version,
		"BaseVersion": change.BaseVersion,
		"Diff":        int(version - change.BaseVersion),
		"Len":         len(prevChangeStrs),
		"ChangeStr":   patchStr,
		"minVersion":  minVersion,
	})

	//startIndex := len(prevChangeStrs) - int(version-change.BaseVersion)
	//if startIndex < 0 {
	//	utils.LogError("StartIndex is negative", ErrVersionOutOfDate, nil)
	//	return nil, -1, nil, ErrVersionOutOfDate
	//}

	startIndex := int64(len(prevChangeStrs) - 1)

	if change.BaseVersion > version {
		// check to make sure the patch is being applied to the most recent revision
		utils.LogError("BaseVersion too high", ErrVersionOutOfDate, nil)
		return "", -1, nil, 0, ErrVersionOutOfDate
	} else if change.BaseVersion == version {
		// If we are building on the server's base version, don't need to transform.
		startIndex = int64(len(prevChangeStrs))
	} else if change.BaseVersion < minVersion {
		// if it's less than the minVersion, we've scrunched.
		utils.LogError("BaseVersion less than minVersion", ErrVersionOutOfDate, nil)
		return "", -1, nil, 0, ErrVersionOutOfDate
	} else if change.BaseVersion == minVersion {
		// If it's equal to the minVersion, we use the entire array
		startIndex = int64(0)
	} else {
		// Otherwise, find the right starting point
		startIndex = int64(len(prevChangeStrs)) - (version - change.BaseVersion)
		for startIndex >= 0 && startIndex < int64(len(prevChangeStrs)) {
			otherPatch, err := patching.NewPatchFromString(prevChangeStrs[startIndex])
			if err != nil {
				utils.LogError("Failed to parse patch", err, utils.LogFields{
					"PatchStr":   strings.Replace(prevChangeStrs[startIndex], "\n", "\\n", -1),
					"StartIndex": startIndex,
				})
				return "", -1, nil, 0, ErrInternalServerError
			}

			if change.BaseVersion > otherPatch.BaseVersion {
				break
			} else {
				startIndex--
			}
		}
		startIndex++ // go back to the actual base version
	}

	// If it's negative at this point, it means we started off with an index that was less than -1.
	// In other words, we've probably scrunched the changes we're looking for.
	if startIndex < 0 {
		utils.LogError("StartIndex was negative", ErrVersionOutOfDate, nil)
		return "", -1, nil, 0, ErrVersionOutOfDate
	}

	if startIndex < minStartIndex {
		minStartIndex = startIndex
	}

	utils.LogDebug("FINISHED CHECKING", utils.LogFields{
		"Change":     patchStr,
		"StartIndex": startIndex,
		"Len":        len(prevChangeStrs),
	})

	// Apply patches from the change's baseVersion onwards
	toApply := prevChangeStrs[startIndex:]

	utils.LogDebug("TRANSFORMING", utils.LogFields{
		"PatchesToApply": toApply,
		"Change":         patchStr,
		"StartIndex":     startIndex,
		"Len":            len(prevChangeStrs),
	})

	transformedPatch := change
	if startIndex != int64(len(prevChangeStrs)) {
		consolidatedPatch, err := patching.ConsolidatePatches(prevChanges[startIndex:])
		if err != nil {
			utils.LogError("Failed to consolidate patches", err, utils.LogFields{
				"Patch":       strings.Replace(change.String(), "\n", "\\n", -1),
				"prevChanges": strings.Replace(spew.Sprint(prevChanges), "\n", "\\n", -1),
			})
		}

		transformResults, err := patching.TransformPatches(change, consolidatedPatch)
		if err != nil {
			utils.LogError("Failed to transform patch", err, utils.LogFields{
				"Patch":             strings.Replace(change.String(), "\n", "\\n", -1),
				"consolidatedPatch": strings.Replace(consolidatedPatch.String(), "\n", "\\n", -1),
			})
			return "", -1, nil, 0, err
		}

		transformedPatch = transformResults.PatchXPrime
		transformedPatch.BaseVersion = version
	}

	// use the cas to make sure the document hasn't changed
	builder := cb.bucket.MutateIn(strconv.FormatInt(fileMeta.FileID, 10), gocb.Cas(cas), 0)

	if !useTemp {
		builder.ArrayAppendMulti("changes", []string{transformedPatch.String()}, false)
	} else {
		builder.ArrayAppendMulti("tempchanges", []string{transformedPatch.String()}, false)
	}

	builder = builder.Counter("version", 1, false)

	_, err = builder.Execute()
	if err != nil {
		return "", -1, nil, 0, err
	}

	// TODO: Evaluate whether prevChangesCopy is the correct item to send back
	// use prevChangesCopy, so we don't send back the transformed patch set
	return transformedPatch.String(), version + 1, prevChangesCopy[minStartIndex:], len(prevChangeStrs) + 1, err
}
