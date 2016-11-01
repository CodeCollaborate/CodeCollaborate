package dbfs

import (
	"strconv"
	"strings"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/couchbase/gocb"
)

type couchbaseConn struct {
	config config.ConnCfg
	bucket *gocb.Bucket
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
		return di.couchbaseDB, err
	}

	if di.couchbaseDB.config.Schema == "" {
		di.couchbaseDB.config.Schema = "documents"
	}

	myBucket, err := documentsCluster.OpenBucket(di.couchbaseDB.config.Schema, di.couchbaseDB.config.Password)
	if err != nil {
		return di.couchbaseDB, err
	}
	di.couchbaseDB.bucket = myBucket

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
func (di *DatabaseImpl) CBAppendFileChange(fileID int64, baseVersion int64, changes []string) (int64, error) {
	// TODO (non-immediate/required): verify changes are valid changes
	cb, err := di.openCouchBase()
	if err != nil {
		return -1, err
	}
	key := strconv.FormatInt(fileID, 10)

	// optimistic locking operation
	// check the version is accurate and get the object's cas,
	// then use it in the MutateIn call to verify the document hasn't updated underneath us
	frag, err := cb.bucket.LookupIn(key).Get("version").Get("usetemp").Execute()
	if err != nil {
		return -1, ErrResourceNotFound
	}

	cas := frag.Cas()
	var version int64
	var useTemp bool
	err = frag.Content("version", &version)
	if err != nil {
		return -1, ErrNoData
	}
	err = frag.Content("usetemp", &useTemp)
	if err != nil {
		return -1, ErrNoData
	}

	// check to make sure the patch is being applied to the most recent revision
	if baseVersion != version {
		return -1, ErrVersionOutOfDate
	}

	// use the cas to make sure the document hasn't changed
	builder := cb.bucket.MutateIn(key, cas, 0)

	if !useTemp {
		builder.ArrayAppendMulti("changes", changes, false)
	} else {
		builder.ArrayAppendMulti("tempchanges", changes, false)
	}

	builder = builder.Counter("version", 1, false)

	_, err = builder.Execute()
	if err != nil {
		return version, err
	}
	return version + 1, err
}
