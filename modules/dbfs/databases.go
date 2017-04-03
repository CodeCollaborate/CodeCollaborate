package dbfs

// DBFS is the interface which maps all of the necessary database and file system functions
type DBFS interface {
	// multi

	// ScrunchFile scrunches the file for the given metadata. All new changes called while scrunching is
	// in progress are redirected, and merged back when done.
	ScrunchFile(meta FileMeta) error

	// getForScrunching gets all but the remainder entries for a file and creates a temp swp file.
	// Returns the changes for scrunching, the swap file contents, and any errors
	getForScrunching(fileMeta FileMeta, remainder int) ([]string, []byte, error)

	// deleteForScrunching deletes `num` elements from the front of `changes` for file with `fileID` and deletes the
	// swp file
	deleteForScrunching(fileMeta FileMeta, num int) error

	// PullFile pulls the changes and the file bytes from the databases
	PullFile(meta FileMeta) (*[]byte, []string, error)

	// PullChanges pulls the changes from the databases and returns them along with the temporary lock value,
	// the file version, and the useTemp flag
	PullChanges(meta FileMeta) ([]string, uint64, int64, bool, error)

	// Couchbase

	// CloseCouchbase closes the CouchBase db connection
	// YOU PROBABLY DON'T NEED TO RUN THIS EVER
	CloseCouchbase() error

	// CBInsertNewFile inserts a new document with the given arguments
	CBInsertNewFile(fileID int64, version int64, changes []string) error

	// CBDeleteFile deletes the document with FileID == fileID from couchbase
	CBDeleteFile(fileID int64) error

	// CBGetFileVersion returns the current version of the file for the given FileID
	CBGetFileVersion(fileID int64) (int64, error)

	// CBAppendFileChange mutates the file document with the new change and sets the new version number
	// Returns the new version number, the missing patches, the total count of patches tracked, and an error, if any.
	CBAppendFileChange(file FileMeta, patches string) (string, int64, []string, int, error)

	// MySQL

	// CloseMySQL closes the MySQL db connection
	// YOU PROBABLY DON'T NEED TO RUN THIS EVER
	CloseMySQL() error

	// MySQLUserRegister registers a new user in MySQL
	MySQLUserRegister(user UserMeta) error

	// MySQLUserGetPass is used to get the key and hash of a stored password to verify that a value is correct
	MySQLUserGetPass(username string) (password string, err error)

	// MySQLUserDelete deletes a user from MySQL
	MySQLUserDelete(username string) ([]int64, error)

	// MySQLUserLookup returns user information about a user with the username 'username'
	MySQLUserLookup(username string) (user UserMeta, err error)

	// MySQLUserProjects returns the projectID, the project name, and the permission level the user `username` has on that project
	MySQLUserProjects(username string) (projects []ProjectMeta, err error)

	// MySQLProjectCreate create a new project in MySQL
	MySQLProjectCreate(username string, projectName string) (projectID int64, err error)

	// MySQLProjectDelete deletes a project from MySQL
	MySQLProjectDelete(projectID int64, senderID string) error

	// MySQLProjectGetFiles returns the Files from the project with projectID = projectID
	MySQLProjectGetFiles(projectID int64) (files []FileMeta, err error)

	// MySQLProjectGrantPermission gives the user `grantUsername` the permission `permissionLevel` on project `projectID`
	MySQLProjectGrantPermission(projectID int64, grantUsername string, permissionLevel int8, grantedByUsername string) error

	// MySQLProjectRevokePermission removes revokeUsername's permissions from the project
	// DOES NOT WORK FOR OWNER (which is kinda a good thing)
	MySQLProjectRevokePermission(projectID int64, revokeUsername string, revokedByUsername string) error

	// MySQLUserProjectPermissionLookup returns the permission level of `username` on the project with the given projectID
	MySQLUserProjectPermissionLookup(projectID int64, username string) (int8, error)

	// MySQLProjectRename allows for you to rename projects
	MySQLProjectRename(projectID int64, newName string) error

	// MySQLProjectLookup returns the project name and permissions for a project with ProjectID = 'projectID'
	// NOTE: There's an important to do on the DatabaseImpl version of this
	MySQLProjectLookup(projectID int64, username string) (name string, permissions map[string]ProjectPermission, err error)

	// MySQLFileCreate create a new file in MySQL
	MySQLFileCreate(username string, filename string, relativePath string, projectID int64) (fileID int64, err error)

	// MySQLFileDelete deletes a file from the MySQL database
	// this does not delete the actual file
	MySQLFileDelete(fileID int64) error

	// MySQLFileMove updates MySQL with the  new path of the file with FileID == 'fileID'
	MySQLFileMove(fileID int64, newPath string) error

	// MySQLFileRename updates MySQL with the new name of the file with FileID == 'fileID'
	MySQLFileRename(fileID int64, newName string) error

	// MySQLFileGetInfo returns the meta data about the given file
	MySQLFileGetInfo(fileID int64) (FileMeta, error)

	// filesystem

	// FileWrite writes the file with the given bytes to a calculated path, and
	// returns that path so it can be put in MySQL
	FileWrite(relpath string, filename string, projectID int64, raw []byte) (string, error)

	// FileDelete deletes the file with the given metadata from the file system
	// Couple this with dbfs.MySQLFileDelete and dbfs.CBDeleteFile
	FileDelete(relpath string, filename string, projectID int64) error

	// FileMove moves a file form the starting path to the end path
	FileMove(startRelpath string, startFilename string, endRelpath string, endFilename string, projectID int64) error

	// FileWriteToSwap writes the swapfile for the file with the given info
	FileWriteToSwap(meta FileMeta, raw []byte) error
}
