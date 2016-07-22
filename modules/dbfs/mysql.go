package dbfs

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // required to load into local namespace to
	// initialize sql driver mapping in sql.Open("mysql", ...)

	"time"

	"path/filepath"
	"strings"

	"fmt"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/utils"
)

var mysqldb *mysqlConn

type mysqlConn struct {
	config config.ConnCfg
	db     *sql.DB
}

func getMySQLConn() (*mysqlConn, error) {
	if mysqldb != nil && mysqldb.db != nil {
		err := mysqldb.db.Ping()
		if err == nil {
			return mysqldb, nil
		}
	}

	if mysqldb == nil || mysqldb.config == (config.ConnCfg{}) {
		mysqldb = new(mysqlConn)
		configMap := config.GetConfig()
		mysqldb.config = configMap.ConnectionConfig["MySQL"]
	}

	if mysqldb.config.Schema == "" {
		mysqldb.config.Schema = "cc"
	}

	connString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?timeout=%vs&parseTime=true", mysqldb.config.Username, mysqldb.config.Password, mysqldb.config.Host, int(mysqldb.config.Port), mysqldb.config.Schema, int(mysqldb.config.Timeout))
	db, err := sql.Open("mysql", connString)
	if err != nil {
		utils.LogOnError(err, "Unable to connect to MySQL")
		return mysqldb, err
	}

	mysqldb.db = db
	return mysqldb, nil
}

// CloseMySQL closes the MySQL db connection
// YOU PROBABLY DON'T NEED TO RUN THIS EVER
func CloseMySQL() error {
	if mysqldb != nil && mysqldb.db != nil {
		err := mysqldb.db.Close()
		mysqldb = nil
		return err
	}
	return ErrDbNotInitialized
}

/**
STORED PROCEDURES
*/

// MySQLUserRegister registers a new user in MySQL
func MySQLUserRegister(user UserMeta) error {
	mysql, err := getMySQLConn()
	if err != nil {
		return err
	}

	result, err := mysql.db.Exec("CALL user_register(?,?,?,?,?)", user.Username, user.Password, user.Email, user.FirstName, user.LastName)
	if err != nil {
		return err
	}
	numrows, err := result.RowsAffected()

	if err != nil || numrows == 0 {
		return ErrNoDbChange
	}

	return nil
}

// MySQLUserGetPass is used to get the key and hash of a stored password to verify that a value is correct
func MySQLUserGetPass(username string) (password string, err error) {
	mysql, err := getMySQLConn()
	if err != nil {
		return "", err
	}

	rows, err := mysql.db.Query("CALL user_get_password(?)", username)
	if err != nil {
		return "", err
	}

	for rows.Next() {
		err = rows.Scan(&password)
		if err != nil {
			return "", err
		}
	}

	return password, nil
}

// MySQLUserDelete deletes a user from MySQL
// unexported b/c not part of the official API
func mySQLUserDelete(username string, pass string) error {
	mysql, err := getMySQLConn()
	if err != nil {
		return err
	}

	// FIXME (optional): use MySQLUserGetPass to verify the password is correct before deleting

	result, err := mysql.db.Exec("CALL user_delete(?)", username)
	if err != nil {
		return err
	}
	numrows, err := result.RowsAffected()

	if err != nil || numrows == 0 {
		return ErrNoDbChange
	}

	return nil
}

// MySQLUserLookup returns user information about a user with the username 'username'
func MySQLUserLookup(username string) (user UserMeta, err error) {
	mysql, err := getMySQLConn()
	if err != nil {
		return user, err
	}

	rows, err := mysql.db.Query("CALL user_lookup(?)", username)
	if err != nil {
		return user, err
	}

	for rows.Next() {
		err = rows.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Username)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

// MySQLUserProjects returns the projectID, the project name, and the permission level the user `username` has on that project
func MySQLUserProjects(username string) (projects []ProjectMeta, err error) {
	mysql, err := getMySQLConn()
	if err != nil {
		return nil, err
	}

	rows, err := mysql.db.Query("CALL user_projects(?)", username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		project := ProjectMeta{}
		err = rows.Scan(&project.ProjectID, &project.ProjectName, &project.PermissionLevel)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// MySQLProjectCreate create a new project in MySQL
func MySQLProjectCreate(username string, projectName string) (projectID int64, err error) {
	mysql, err := getMySQLConn()
	if err != nil {
		return -1, err
	}

	rows, err := mysql.db.Query("CALL project_create(?,?)", projectName, username)
	if err != nil {
		return -1, err
	}
	for rows.Next() {
		err = rows.Scan(&projectID)
		if err != nil {
			return -1, err
		}
	}

	return projectID, nil
}

// MySQLProjectDelete deletes a project from MySQL
func MySQLProjectDelete(projectID int64, senderID string) error {
	mysql, err := getMySQLConn()
	if err != nil {
		return err
	}

	result, err := mysql.db.Exec("CALL project_delete(?,?)", projectID, senderID)
	if err != nil {
		return err
	}
	numrows, err := result.RowsAffected()

	if err != nil || numrows == 0 {
		return ErrNoDbChange
	}
	return nil
}

// MySQLProjectGetFiles returns the Files from the project with projectID = projectID
func MySQLProjectGetFiles(projectID int64) (files []FileMeta, err error) {
	mysql, err := getMySQLConn()
	if err != nil {
		return nil, err
	}

	rows, err := mysql.db.Query("CALL project_get_files(?)", projectID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		file := FileMeta{}
		err = rows.Scan(&file.FileID, &file.Creator, &file.CreationDate, &file.RelativePath, &file.ProjectID, &file.Filename)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

// MySQLProjectGrantPermission gives the user `grantUsername` the permision `permissionLevel` on project `projectID`
func MySQLProjectGrantPermission(projectID int64, grantUsername string, permissionLevel int, grantedByUsername string) error {
	mysql, err := getMySQLConn()
	if err != nil {
		return err
	}

	result, err := mysql.db.Exec("CALL project_grant_permissions(?, ?, ?, ?)", projectID, grantUsername, permissionLevel, grantedByUsername)
	if err != nil {
		return err
	}
	numrows, err := result.RowsAffected()

	if err != nil || numrows == 0 {
		return ErrNoDbChange
	}
	return nil
}

// MySQLProjectRevokePermission removes revokeUsername's permissions from the project
// DOES NOT WORK FOR OWNER
func MySQLProjectRevokePermission(projectID int64, revokeUsername string, revokedByUsername string) error {
	mysql, err := getMySQLConn()
	if err != nil {
		return err
	}

	result, err := mysql.db.Exec("CALL project_revoke_permissions(?, ?)", projectID, revokeUsername)
	if err != nil {
		return err
	}
	numrows, err := result.RowsAffected()

	if err != nil || numrows == 0 {
		return ErrNoDbChange
	}
	return nil
}

// MySQLProjectRename allows for you to rename projects
func MySQLProjectRename(projectID int64, newName string) error {
	mysql, err := getMySQLConn()
	if err != nil {
		return err
	}

	result, err := mysql.db.Exec("CALL project_rename(?, ?)", projectID, newName)
	if err != nil {
		return err
	}
	numrows, err := result.RowsAffected()

	if err != nil || numrows == 0 {
		return ErrNoDbChange
	}
	return nil
}

// MySQLProjectLookup returns the project name and permissions for a project with ProjectID = 'projectID'
//
// TODO (non-immediate/required): decide on change to MySQLProjectLookup stored proc
// aka: decide if looking them up 1 at a time is good or not
// Looking them up 1 at a time may seem worse, however we're looking up rows based on their primary key
// so we get the speed benefits of it having a unique index on it
// Thoughts:
// 		FIND_IN_SET doesn't use any indices at all,
// 		both IN and FIND_IN_SET have issues with integers
// 		more issues when there are a variable number of ID's because MySQL doesn't have arrays
//
// http://stackoverflow.com/a/8150183 <- preferred if we switch b/c FIND_IN_SET doesn't use indexes
func MySQLProjectLookup(projectID int64, username string) (name string, permissions map[string]ProjectPermission, err error) {
	permissions = make(map[string](ProjectPermission))
	mysql, err := getMySQLConn()
	if err != nil {
		return "", permissions, err
	}

	// TODO (optional): un-hardcode '10' as the owner constant in the MySQL ProjectLookup stored proc

	rows, err := mysql.db.Query("CALL project_lookup(?)", projectID)
	if err != nil {
		return "", permissions, err
	}

	var hasAccess = false
	for rows.Next() {
		perm := ProjectPermission{}
		var hora string
		err = rows.Scan(&name, &perm.Username, &perm.PermissionLevel, &perm.GrantedBy, &hora)
		perm.GrantedDate, _ = time.Parse("2006-01-02 15:04:05", hora)
		if err != nil {
			return "", permissions, err
		}
		if !hasAccess && perm.PermissionLevel > 0 && perm.Username == username {
			hasAccess = true
		}
		permissions[perm.Username] = perm
	}

	// verify user has access to view this info
	if hasAccess {
		return name, permissions, err
	}
	return "", make(map[string](ProjectPermission)), err
}

// MySQLFileCreate create a new file in MySQL
func MySQLFileCreate(username string, filename string, relativePath string, projectID int64) (fileID int64, err error) {
	if strings.Contains(filename, filePathSeparator) {
		return -1, ErrMaliciousRequest
	}

	filepathClean := filepath.Clean(relativePath)
	if strings.HasPrefix(filepathClean, "..") {
		return -1, ErrMaliciousRequest
	}

	mysql, err := getMySQLConn()
	if err != nil {
		return -1, err
	}

	rows, err := mysql.db.Query("CALL file_create(?,?,?,?)", username, filename, filepathClean, projectID)
	if err != nil {
		return -1, err
	}
	for rows.Next() {
		err = rows.Scan(&fileID)
		if err != nil {
			return -1, err
		}
	}

	return fileID, nil
}

// MySQLFileDelete deletes a file from the MySQL database
// this does not delete the actual file
func MySQLFileDelete(fileID int64) error {
	mysql, err := getMySQLConn()
	if err != nil {
		return err
	}

	result, err := mysql.db.Exec("CALL file_delete(?)", fileID)
	if err != nil {
		return err
	}
	numrows, err := result.RowsAffected()

	if err != nil || numrows == 0 {
		return ErrNoDbChange
	}
	return nil
}

// MySQLFileMove updates MySQL with the  new path of the file with FileID == 'fileID'
func MySQLFileMove(fileID int64, newPath string) error {
	newPathClean := filepath.Clean(newPath)
	if strings.HasPrefix(newPathClean, "..") {
		return ErrMaliciousRequest
	}

	mysql, err := getMySQLConn()
	if err != nil {
		return err
	}

	result, err := mysql.db.Exec("CALL file_move(?, ?)", fileID, newPathClean)
	if err != nil {
		return err
	}
	numrows, err := result.RowsAffected()

	if err != nil || numrows == 0 {
		return ErrNoDbChange
	}
	return nil
}

// MySQLFileRename updates MySQL with the new name of the file with FileID == 'fileID'
func MySQLFileRename(fileID int64, newName string) error {
	if strings.Contains(newName, filePathSeparator) {
		return ErrMaliciousRequest
	}

	mysql, err := getMySQLConn()
	if err != nil {
		return err
	}

	result, err := mysql.db.Exec("CALL file_rename(?, ?)", fileID, newName)
	if err != nil {
		return err
	}
	numrows, err := result.RowsAffected()

	if err != nil || numrows == 0 {
		return ErrNoDbChange
	}
	return nil
}

// MySQLFileGetInfo returns the meta data about the given file
func MySQLFileGetInfo(fileID int64) (FileMeta, error) {
	file := FileMeta{}
	mysql, err := getMySQLConn()
	if err != nil {
		return file, err
	}

	rows, err := mysql.db.Query("CALL file_get_info(?)", fileID)
	if err != nil {
		return file, err
	}

	file.FileID = fileID
	for rows.Next() {
		err = rows.Scan(&file.Creator, &file.CreationDate, &file.RelativePath, &file.ProjectID, &file.Filename)
		if err != nil {
			return file, err
		}
	}

	return file, nil
}
