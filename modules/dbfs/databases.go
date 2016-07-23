package dbfs

// Dbfs is the globally used dbfs object for the server
var Dbfs DBFS

// DBFS is the interface which maps all of the necessary database and file system functions
type DBFS interface {
	// couchbase

	CloseCouchbase() error
	CBInsertNewFile(fileID int64, version int64, changes []string) error
	CBDeleteFile(fileID int64) error
	CBGetFileVersion(fileID int64) (int64, error)
	CBGetFileChanges(fileID int64) ([]string, error)
	CBAppendFileChange(fileID int64, baseVersion int64, changes []string) (int64, error)

	// mysql

	CloseMySQL() error
	MySQLUserRegister(user UserMeta) error
	MySQLUserGetPass(username string) (password string, err error)
	MySQLUserDelete(username string, pass string) error
	MySQLUserLookup(username string) (user UserMeta, err error)
	MySQLUserProjects(username string) (projects []ProjectMeta, err error)
	MySQLProjectCreate(username string, projectName string) (projectID int64, err error)
	MySQLProjectDelete(projectID int64, senderID string) error
	MySQLProjectGetFiles(projectID int64) (files []FileMeta, err error)
	MySQLProjectGrantPermission(projectID int64, grantUsername string, permissionLevel int, grantedByUsername string) error
	MySQLProjectRevokePermission(projectID int64, revokeUsername string, revokedByUsername string) error
	MySQLProjectRename(projectID int64, newName string) error
	MySQLProjectLookup(projectID int64, username string) (name string, permissions map[string]ProjectPermission, err error)
	MySQLFileCreate(username string, filename string, relativePath string, projectID int64) (fileID int64, err error)
	MySQLFileDelete(fileID int64) error
	MySQLFileMove(fileID int64, newPath string) error
	MySQLFileRename(fileID int64, newName string) error
	MySQLFileGetInfo(fileID int64) (FileMeta, error)

	// filesystem

	FileWrite(relpath string, filename string, projectID int64, raw []byte) (string, error)
	FileDelete(relpath string, filename string, projectID int64) error
	FileRead(relpath string, filename string, projectID int64) (*[]byte, error)
}

// DatabaseImpl is the concrete implementation of the DBFS interface
type DatabaseImpl struct {
	couchbaseDB *couchbaseConn
	mysqldb     *mysqlConn
}
