package dbfs

// DatabaseImpl is the concrete implementation of the DBFS interface
type DatabaseImpl struct {
	couchbaseDB *couchbaseConn
	mysqldb     *mysqlConn
}
