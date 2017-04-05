package documentstore

import (
	"strconv"
	"strings"

	"fmt"
	"time"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datastore"
	"github.com/CodeCollaborate/Server/modules/patching"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/couchbase/gocb"
	"github.com/davecgh/go-spew/spew"
)

// CouchbaseDocumentStore is the concrete implementation of the DocumentStore, using Couchbase as the backing document store of choice
type CouchbaseDocumentStore struct {
	cfg    *config.ConnCfg
	bucket *gocb.Bucket
}

// NewCouchbaseDocumentStore creates a new instance of the CouchbaseDocumentStore, setting the configuration
func NewCouchbaseDocumentStore(cfg *config.ConnCfg) *CouchbaseDocumentStore {
	return &CouchbaseDocumentStore{
		cfg: cfg,
	}
}

// Connect starts this CouchbaseDocumentStore's connection to the server
func (store *CouchbaseDocumentStore) Connect() {
	if store.bucket != nil {
		return
	}

	var documentsCluster *gocb.Cluster
	var err error
	for i := 0; i < store.cfg.NumRetries; i++ {
		if strings.HasPrefix(store.cfg.Host, "couchbase://") {
			documentsCluster, err = gocb.Connect(fmt.Sprintf("%s:%d", store.cfg.Host, store.cfg.Port))
		} else {
			documentsCluster, err = gocb.Connect(fmt.Sprintf("couchbase://%s:%d", store.cfg.Host, store.cfg.Port))
		}

		// If there was an error, try again after a few seconds
		if err == nil {
			break
		} else {
			utils.LogError("Couchbase: could not connect to Couchbase instance", err, utils.LogFields{
				"Host":    store.cfg.Host,
				"Port":    store.cfg.Port,
				"Attempt": i + 1,
			})
			time.Sleep(time.Duration(store.cfg.Timeout) * time.Second)
		}
	}
	if err != nil {
		utils.LogFatal("Couchbase: failed to open bucket", err, utils.LogFields{
			"Host":     store.cfg.Host,
			"Port":     store.cfg.Port,
			"Schema":   store.cfg.Schema,
			"Attempts": store.cfg.NumRetries,
		})
	}

	// Set a default schema if does not exist
	if store.cfg.Schema == "" {
		store.cfg.Schema = "documents"
	}

	var schemaBucket *gocb.Bucket
	for i := 0; i < store.cfg.NumRetries; i++ {
		schemaBucket, err = documentsCluster.OpenBucket(store.cfg.Schema, store.cfg.Password)

		// If there was an error, try again after a few seconds
		if err == nil {
			store.bucket = schemaBucket
			break
		} else {
			utils.LogError("Couchbase: failed to open bucket", err, utils.LogFields{
				"Host":    store.cfg.Host,
				"Port":    store.cfg.Port,
				"Schema":  store.cfg.Schema,
				"Attempt": i + 1,
			})
			time.Sleep(time.Duration(store.cfg.Timeout) * time.Second)
		}
	}
	if err != nil {
		utils.LogFatal("Couchbase: failed to open bucket", err, utils.LogFields{
			"Host":     store.cfg.Host,
			"Port":     store.cfg.Port,
			"Schema":   store.cfg.Schema,
			"Attempts": store.cfg.NumRetries,
		})
	}
}

// Shutdown terminates this CouchbaseDocumentStore's connection to the server
func (store *CouchbaseDocumentStore) Shutdown() {
	store.bucket.Close()
}

// AddFileData stores the given FileData using the internal FileID
func (store *CouchbaseDocumentStore) AddFileData(data *datastore.FileData) error {
	if store.bucket == nil {
		store.Connect()
	}

	_, err := store.bucket.Insert(strconv.FormatInt(data.FileID, 10), data, 0)

	if err != nil {
		switch err {
		case gocb.ErrKeyExists:
			return datastore.ErrFileAlreadyExists
		default:
			utils.LogError("Couchbase: Failed to add file data", err, utils.LogFields{
				"FileData": data,
			})
			return datastore.ErrInternalServerErr
		}
	}

	return nil
}

// SetFileData stores the given FileData using the internal FileID
func (store *CouchbaseDocumentStore) SetFileData(data *datastore.FileData) error {
	if store.bucket == nil {
		store.Connect()
	}

	_, err := store.bucket.Upsert(strconv.FormatInt(data.FileID, 10), data, 0)

	if err != nil {
		switch err {
		default:
			utils.LogError("Couchbase: Failed to set file data", err, utils.LogFields{
				"FileData": data,
			})
			return datastore.ErrInternalServerErr
		}
	}

	return nil
}

// GetFileData retrieves the FileData for the given fileID
func (store *CouchbaseDocumentStore) GetFileData(fileID int64) (*datastore.FileData, error) {
	if store.bucket == nil {
		store.Connect()
	}

	fileData, _, err := store.getFileDataWithCas(fileID)

	return fileData, err
}

func (store *CouchbaseDocumentStore) getFileDataWithCas(fileID int64) (*datastore.FileData, gocb.Cas, error) {
	if store.bucket == nil {
		store.Connect()
	}

	fileData := &datastore.FileData{}

	cas, err := store.bucket.Get(strconv.FormatInt(fileID, 10), fileData)
	if err != nil {
		switch err {
		case gocb.ErrKeyNotFound:
			return nil, 0, datastore.ErrFileDoesNotExist
		default:
			utils.LogError("Couchbase: Failed to get file data", err, utils.LogFields{
				"FileID": fileID,
			})
			return nil, 0, datastore.ErrInternalServerErr
		}
	}

	return fileData, cas, nil

}

// DeleteFileData deletes the FileData for the given fileID
func (store *CouchbaseDocumentStore) DeleteFileData(fileID int64) error {
	if store.bucket == nil {
		store.Connect()
	}

	_, err := store.bucket.Remove(strconv.FormatInt(fileID, 10), 0)
	if err != nil {
		switch err {
		case gocb.ErrKeyNotFound:
			return datastore.ErrFileDoesNotExist
		default:
			utils.LogError("Couchbase: Failed to delete file data", err, utils.LogFields{
				"FileID": fileID,
			})
			return datastore.ErrInternalServerErr
		}
	}

	return nil
}

// AppendPatch appends the patch to the document with the given fileID, and returns the resultant FileData and missing patches if successful
func (store *CouchbaseDocumentStore) AppendPatch(fileID int64, patch *patching.Patch) (*datastore.FileData, []string, error) {
	if store.bucket == nil {
		store.Connect()
	}

	// Retrieve file data and CAS value; will be used later in optimistic write
	fileData, cas, err := store.getFileDataWithCas(fileID)
	if err != nil {
		return nil, nil, err
	}

	prevPatches, err := patching.GetPatches(fileData.Changes)
	if err != nil {
		utils.LogError("Failed to parse previous patches into patch objects", err, utils.LogFields{
			"PrevPatches": fileData.Changes,
		})
		return nil, nil, datastore.ErrInternalServerErr
	}

	if cas == 0 {
		utils.LogWarn("Couchbase returned a CAS value of 0, optimistic locking unavailable", utils.LogFields{
			"CAS":    cas,
			"FileID": fileID,
		})
	}

	// Find the minimum version that we have in the document store
	minVersion := fileData.Version
	if len(prevPatches) > 0 {
		// Allow transform-patches to start on the same base version as the head (after linearization, we have all the necessary patches)
		// The head here is defined as the first patch that has not been scrunched
		minVersion = prevPatches[fileData.ScrunchedPatchCount].BaseVersion
	}

	// Transform changes against newer changes.
	// Calculate the patches that it does not have.
	toApply := []*patching.Patch{}

	if patch.BaseVersion > fileData.Version {
		// If base version is the same as the current version, we are up to date; no need to transform

	} else if patch.BaseVersion > fileData.Version {
		// check to make sure the patch is being applied to the most recent revision
		utils.LogError("BaseVersion too high", datastore.ErrFileBaseVersionTooHigh, nil)
		return nil, nil, datastore.ErrFileBaseVersionTooHigh

	} else if patch.BaseVersion < minVersion {
		// if it's less than the minVersion, we've scrunched.
		utils.LogError("BaseVersion less than minVersion", datastore.ErrFileBaseVersionTooLow, nil)
		return nil, nil, datastore.ErrFileBaseVersionTooLow

	} else if patch.BaseVersion == minVersion {
		// If it's equal to the minVersion, we use the entire array
		toApply = prevPatches

	} else {
		// Otherwise, find the right starting point
		// There is exactly one patch per version incremented (for non-scrunched patches)
		startIndex := len(prevPatches) - int(fileData.Version-patch.BaseVersion)
		if startIndex < 0 {
			utils.LogError("StartIndex was negative", datastore.ErrInternalServerErr, nil)
			return nil, nil, datastore.ErrInternalServerErr
		}
		toApply = prevPatches[startIndex:]
	}

	utils.LogDebug("Transforming incoming patch against missing patches", utils.LogFields{
		"PatchesToApply": toApply,
		"Change":         strings.Replace(patch.String(), "\n", "\\n", -1),
	})

	transformedPatch := patch
	if len(toApply) > 0 {
		consolidatedMissingPatches, err := patching.ConsolidatePatches(toApply)
		if err != nil {
			utils.LogError("Failed to consolidate missing patches", err, utils.LogFields{
				"Patch":       strings.Replace(patch.String(), "\n", "\\n", -1),
				"prevPatches": strings.Replace(spew.Sprint(prevPatches), "\n", "\\n", -1),
			})
			return nil, nil, datastore.ErrInternalServerErr
		}

		transformResults, err := patching.TransformPatches(patch, consolidatedMissingPatches)
		if err != nil {
			utils.LogError("Failed to transform patch against missing patches", err, utils.LogFields{
				"Patch": strings.Replace(patch.String(), "\n", "\\n", -1),
				"consolidatedMissingPatches": strings.Replace(consolidatedMissingPatches.String(), "\n", "\\n", -1),
			})
			return nil, nil, datastore.ErrInternalServerErr
		}
		transformedPatch = transformResults.PatchXPrime

		// TODO(wongb): Is this necessary anymore?
		//transformedPatch.BaseVersion = version
	}

	// use the cas to make sure the document hasn't changed
	builder := store.bucket.MutateIn(strconv.FormatInt(fileData.FileID, 10), gocb.Cas(cas), 0)
	builder.ArrayAppendMulti("Changes", []string{transformedPatch.String()}, false)
	builder.Upsert("LastModifiedTime", time.Now().UnixNano()/int64(time.Millisecond), false)
	builder = builder.Counter("Version", 1, false)

	_, err = builder.Execute()
	if err != nil {
		switch err {
		// TODO(wongb): Add CAS failure case
		case gocb.ErrKeyNotFound:
			return nil, nil, datastore.ErrFileDoesNotExist
		default:
			utils.LogError("Couchbase: Failed to delete file data", err, utils.LogFields{
				"FileID": fileID,
			})
			return nil, nil, datastore.ErrInternalServerErr
		}
	}

	toApplyStr := fileData.Changes[len(fileData.Changes)-len(toApply):]
	fileData.Changes = append(fileData.Changes, transformedPatch.String())

	return fileData, toApplyStr, err
}

// ScrunchChanges appends the patch to the document with the given fileID, and returns the resultant FileData if successful
func (store *CouchbaseDocumentStore) ScrunchChanges(fileID int64) error {
	if store.bucket == nil {
		store.Connect()
	}

	return datastore.ErrNotYetImplemented
}
