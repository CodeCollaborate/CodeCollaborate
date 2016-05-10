package dbfs

import (
	"strconv"
	"strings"

	"github.com/CodeCollaborate/Server/modules/config"
	"gopkg.in/couchbase/gocb.v1"
)

var couchbasedb *couchbaseConn

type couchbaseConn struct {
	config config.ConnCfg
	bucket *gocb.Bucket
}

type cbFile struct {
	FileID  int64    `json:"-"`
	Version int64    `json:"version"`
	Changes []string `json:"changes"`
}

func openCouchBase() (*couchbaseConn, error) {
	if couchbasedb != nil && couchbasedb.bucket != nil {
		return couchbasedb, nil
	}

	if couchbasedb == nil || couchbasedb.config == (config.ConnCfg{}) {
		couchbasedb = new(couchbaseConn)
		configMap := config.GetConfig()
		couchbasedb.config = configMap.ConnectionConfig["Couchbase"]
	}

	var documentsCluster *gocb.Cluster
	var err error

	if strings.HasPrefix(couchbasedb.config.Host, "couchbase://") {
		documentsCluster, err = gocb.Connect(couchbasedb.config.Host + ":" + strconv.Itoa(int(couchbasedb.config.Port)))
	} else {
		documentsCluster, err = gocb.Connect("couchbase://" + couchbasedb.config.Host + ":" + strconv.Itoa(int(couchbasedb.config.Port)))
	}

	if err != nil {
		return couchbasedb, err
	}

	if couchbasedb.config.Schema == "" {
		couchbasedb.config.Schema = "documents"
	}

	myBucket, err := documentsCluster.OpenBucket(couchbasedb.config.Schema, couchbasedb.config.Password)
	if err != nil {
		return couchbasedb, err
	}
	couchbasedb.bucket = myBucket

	return couchbasedb, nil
}

// CloseCouchbase closes the CouchBase db connection
// YOU PROBABLY DON'T NEED TO RUN THIS EVER
func CloseCouchbase() error {
	if couchbasedb != nil && couchbasedb.bucket != nil {
		couchbasedb.bucket.Close()
		couchbasedb = nil
	} else {
		return ErrDbNotInitialized
	}

	return nil
}

// CBInsertNewFile inserts a new document into couchbase with CBFile.FileID == fileID
func cbInsertNewFile(file cbFile) error {
	cb, err := openCouchBase()
	//defer cb.close()

	if err != nil {
		return err
	}

	_, err = cb.bucket.Insert(strconv.FormatInt(file.FileID, 10), file, 0)
	return err
}

// CBInsertNewFile inserts a new document with the given arguments
func CBInsertNewFile(fileID int64, version int64, changes []string) error {
	return cbInsertNewFile(cbFile{FileID: fileID, Version: version, Changes: changes})
}

// CBDeleteFile deletes the document with FileID == fileID from couchbase
func CBDeleteFile(fileID int64) error {
	cb, err := openCouchBase()
	if err != nil {
		return err
	}
	_, err = cb.bucket.Remove(strconv.FormatInt(fileID, 10), 0)
	return err
}

// CBGetFileVersion returns the current version of the file for the given FileID
func CBGetFileVersion(fileID int64) (int64, error) {
	cb, err := openCouchBase()
	if err != nil {
		return -1, err
	}

	frag, err := cb.bucket.LookupIn(strconv.FormatInt(fileID, 10)).Get("version").Execute()
	if err != nil {
		return -1, err
	}

	var version int64
	frag.Content("version", &version)
	return version, err
}

// CBGetFileChanges returns the array of file changes for the given fileID
func CBGetFileChanges(fileID int64) ([]string, error) {
	cb, err := openCouchBase()
	if err != nil {
		return []string{}, err
	}

	frag, err := cb.bucket.LookupIn(strconv.FormatInt(fileID, 10)).Get("changes").Execute()
	if err != nil {
		return []string{}, err
	}

	var changes []string
	frag.Content("changes", &changes)

	return changes, err
}

// CBAppendFileChange mutates the file document with the new change and sets the new version number
func CBAppendFileChange(fileID int64, version int64, change string) error {
	cb, err := openCouchBase()
	//defer cb.close()

	if err != nil {
		return err
	}

	_, err = cb.bucket.MutateIn(strconv.FormatInt(fileID, 10), 0, 0).PushBack("changes", change, false).Replace("version", version).Execute()
	return err
}
