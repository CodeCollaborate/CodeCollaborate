package dbfs

import (
	"time"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datastore"
	"github.com/CodeCollaborate/Server/modules/datastore/documentstore"
	"github.com/CodeCollaborate/Server/modules/patching"
	"github.com/couchbase/gocb"
)

type couchbaseConn struct {
	config                 config.ConnCfg
	couchbaseDocumentStore datastore.DocumentStore

	bucket                *gocb.Bucket
	scrunchingLocksBucket *gocb.Bucket
}

type cbFile struct {
	FileID           int64
	Version          int64
	Changes          []string
	TempChanges      []string
	RemainingChanges []string
	UseTemp          bool
	PullSwp          bool
}

func (di *DatabaseImpl) openCouchBase() (*couchbaseConn, error) {
	if di.couchbaseDB == nil || di.couchbaseDB.config == (config.ConnCfg{}) {
		di.couchbaseDB = new(couchbaseConn)
		configMap := config.GetConfig()
		di.couchbaseDB.config = configMap.ConnectionConfig["Couchbase"]
	}

	if di.couchbaseDB.couchbaseDocumentStore == nil {
		di.couchbaseDB.couchbaseDocumentStore = documentstore.NewCouchbaseDocumentStore(&di.couchbaseDB.config)
		di.couchbaseDB.couchbaseDocumentStore.Connect()

		couchbaseDocumentStore := di.couchbaseDB.couchbaseDocumentStore.(*documentstore.CouchbaseDocumentStore)
		di.couchbaseDB.bucket, di.couchbaseDB.scrunchingLocksBucket = couchbaseDocumentStore.GetBuckets()
	}

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
	di.openCouchBase()

	data := datastore.FileData{
		FileID:              file.FileID,
		LastModifiedDate:    time.Now().UnixNano(),
		Version:             file.Version,
		TempChanges:         file.TempChanges,
		UseTemp:             file.UseTemp,
		RemainingChanges:    file.RemainingChanges,
		PullSwp:             file.PullSwp,
		ScrunchedPatchCount: 0,
		Changes:             file.Changes,
	}

	return di.couchbaseDB.couchbaseDocumentStore.AddFileData(&data)
}

// CBInsertNewFile inserts a new document with the given arguments
func (di *DatabaseImpl) CBInsertNewFile(fileID int64, version int64, changes []string) error {
	di.openCouchBase()

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
	di.openCouchBase()

	return di.couchbaseDB.couchbaseDocumentStore.DeleteFileData(fileID)
}

// CBGetFileVersion returns the current version of the file for the given FileID
func (di *DatabaseImpl) CBGetFileVersion(fileID int64) (int64, error) {
	di.openCouchBase()

	fileData, err := di.couchbaseDB.couchbaseDocumentStore.GetFileData(fileID)

	if err != nil {
		return -1, err
	}

	return fileData.Version, nil
}

// CBAppendFileChange mutates the file document with the new change and sets the new version number
// Returns the new version number, the missing patches, the total count of patches tracked, and an error, if any.
func (di *DatabaseImpl) CBAppendFileChange(fileMeta FileMeta, patchStr string) (string, int64, []string, int, error) {
	di.openCouchBase()

	patch, err := patching.NewPatchFromString(patchStr)
	if err != nil {
		return "", -1, nil, 0, err
	}
	newFileData, missingChanges, err := di.couchbaseDB.couchbaseDocumentStore.AppendPatch(fileMeta.FileID, patch)
	return newFileData.AggregatedChanges()[len(newFileData.AggregatedChanges())-1], newFileData.Version, missingChanges, len(newFileData.AggregatedChanges()), err
}
