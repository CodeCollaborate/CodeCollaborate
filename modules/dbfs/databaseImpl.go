package dbfs

import (
	"github.com/CodeCollaborate/Server/modules/datastore"
)

// DatabaseImpl is the concrete implementation of the DBFS interface
type DatabaseImpl struct {
	couchbaseDB *couchbaseConn
	mysqldb     *mysqlConn
	bucketStore datastore.BucketStore
}
