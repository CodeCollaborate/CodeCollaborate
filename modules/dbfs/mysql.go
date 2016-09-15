package dbfs

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // required to load into local namespace to
	// initialize sql driver mapping in sql.Open("mysql", ...)
	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/utils"
)

type mysqlConn struct {
	config config.ConnCfg
	db     *sql.DB
}

func (di *DatabaseImpl) getMySQLConn() (*mysqlConn, error) {
	if di.mysqldb != nil && di.mysqldb.db != nil {
		err := di.mysqldb.db.Ping()
		if err == nil {
			return di.mysqldb, nil
		}
	}

	if di.mysqldb == nil || di.mysqldb.config == (config.ConnCfg{}) {
		di.mysqldb = new(mysqlConn)
		configMap := config.GetConfig()
		di.mysqldb.config = configMap.ConnectionConfig["MySQL"]
	}

	if di.mysqldb.config.Schema == "" {
		panic("No MySQL schema found in config")
	}

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&parseTime=true",
		di.mysqldb.config.Username,
		di.mysqldb.config.Password,
		di.mysqldb.config.Host,
		di.mysqldb.config.Port,
		di.mysqldb.config.Schema,
		di.mysqldb.config.Timeout)
	db, err := sql.Open("mysql", connString)
	if err == nil {
		if err = db.Ping(); err != nil {
			di.mysqldb = nil
			err = ErrDbNotInitialized
		} else {
			di.mysqldb.db = db
		}
	}

	utils.LogOnError(err, "Unable to connect to MySQL")
	return di.mysqldb, err
}

// CloseMySQL closes the MySQL db connection
// YOU PROBABLY DON'T NEED TO RUN THIS EVER
func (di *DatabaseImpl) CloseMySQL() error {
	if di.mysqldb != nil && di.mysqldb.db != nil {
		err := di.mysqldb.db.Close()
		di.mysqldb = nil
		return err
	}
	return ErrDbNotInitialized
}

/**
STORED PROCEDURES
*/

// MySQLUserRegister registers a new user in MySQL
func (di *DatabaseImpl) MySQLUserRegister(user UserMeta) error {
	mysql, err := di.getMySQLConn()
	if err != nil {
		return err
	}

	result, err := mysql.db.Exec("CALL user_register(?,?,?,?,?)", user.Username, user.Password, user.Email, user.FirstName, user.LastName)
	if err != nil {
		return err
	}
	numRows, err := result.RowsAffected()

	if err != nil || numRows == 0 {
		return ErrNoDbChange
	}

	return nil
}

// MySQLUserGetPass is used to get the key and hash of a stored password to verify that a value is correct
func (di *DatabaseImpl) MySQLUserGetPass(username string) (password string, err error) {
	mysql, err := di.getMySQLConn()
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
// technically not part of the official API
func (di *DatabaseImpl) MySQLUserDelete(username string, pass string) error {
	mysql, err := di.getMySQLConn()
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
func (di *DatabaseImpl) MySQLUserLookup(username string) (user UserMeta, err error) {
	mysql, err := di.getMySQLConn()
	if err != nil {
		return user, err
	}

	rows, err := mysql.db.Query("CALL user_lookup(?)", username)
	if err != nil {
		return user, err
	}

	result := false
	for rows.Next() {
		err = rows.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Username)
		if err != nil {
			return user, err
		}
		result = true
	}
	if !result {
		return user, ErrNoData
	}
	return user, nil
}

// MySQLUserProjects returns the projectID, the project name, and the permission level the user `username` has on that project
func (di *DatabaseImpl) MySQLUserProjects(username string) ([]ProjectMeta, error) {
	mysql, err := di.getMySQLConn()
	if err != nil {
		return nil, err
	}

	rows, err := mysql.db.Query("CALL user_projects(?)", username)
	if err != nil {
		return nil, err
	}

	projects := []ProjectMeta{}

	for rows.Next() {
		project := ProjectMeta{}
		err = rows.Scan(&project.ProjectID, &project.Name, &project.PermissionLevel)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// MySQLProjectCreate create a new project in MySQL
func (di *DatabaseImpl) MySQLProjectCreate(username string, projectName string) (projectID int64, err error) {
	mysql, err := di.getMySQLConn()
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
func (di *DatabaseImpl) MySQLProjectDelete(projectID int64, senderID string) error {
	mysql, err := di.getMySQLConn()
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
func (di *DatabaseImpl) MySQLProjectGetFiles(projectID int64) (files []FileMeta, err error) {
	mysql, err := di.getMySQLConn()
	if err != nil {
		return nil, err
	}

	rows, err := mysql.db.Query("CALL project_get_files(?)", projectID)
	if err != nil {
		return nil, err
	}

	files = []FileMeta{}

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

// MySQLProjectGrantPermission gives the user `grantUsername` the permission `permissionLevel` on project `projectID`
func (di *DatabaseImpl) MySQLProjectGrantPermission(projectID int64, grantUsername string, permissionLevel int, grantedByUsername string) error {
	mysql, err := di.getMySQLConn()
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
// DOES NOT WORK FOR OWNER (which is kinda a good thing)
func (di *DatabaseImpl) MySQLProjectRevokePermission(projectID int64, revokeUsername string, revokedByUsername string) error {
	mysql, err := di.getMySQLConn()
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
func (di *DatabaseImpl) MySQLProjectRename(projectID int64, newName string) error {
	mysql, err := di.getMySQLConn()
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
func (di *DatabaseImpl) MySQLProjectLookup(projectID int64, username string) (name string, permissions map[string]ProjectPermission, err error) {
	permissions = make(map[string](ProjectPermission))
	mysql, err := di.getMySQLConn()
	if err != nil {
		return "", permissions, err
	}

	// TODO (optional): un-hardcode '10' as the owner constant in the MySQL ProjectLookup stored proc

	rows, err := mysql.db.Query("CALL project_lookup(?)", projectID)
	if err != nil {
		return "", permissions, err
	}

	result := false
	var hasAccess = false
	for rows.Next() {
		perm := ProjectPermission{}
		var timeVal string
		err = rows.Scan(&name, &perm.Username, &perm.PermissionLevel, &perm.GrantedBy, &timeVal)
		perm.GrantedDate, _ = time.Parse("2006-01-02 15:04:05", timeVal)
		if err != nil {
			return "", permissions, err
		}
		if !hasAccess && perm.PermissionLevel > 0 && perm.Username == username {
			hasAccess = true
		}
		permissions[perm.Username] = perm
		result = true
	}

	// verify user has access to view this info
	if !result || !hasAccess {
		return "", make(map[string](ProjectPermission)), ErrNoData
	}
	return name, permissions, err
}

// MySQLFileCreate create a new file in MySQL
func (di *DatabaseImpl) MySQLFileCreate(username string, filename string, relativePath string, projectID int64) (fileID int64, err error) {
	if strings.Contains(filename, filePathSeparator) {
		return -1, ErrMaliciousRequest
	}

	filepathClean := filepath.Clean(relativePath)
	if strings.HasPrefix(filepathClean, "..") {
		return -1, ErrMaliciousRequest
	}

	mysql, err := di.getMySQLConn()
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
func (di *DatabaseImpl) MySQLFileDelete(fileID int64) error {
	mysql, err := di.getMySQLConn()
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
func (di *DatabaseImpl) MySQLFileMove(fileID int64, newPath string) error {
	newPathClean := filepath.Clean(newPath)
	if strings.HasPrefix(newPathClean, "..") {
		return ErrMaliciousRequest
	}

	mysql, err := di.getMySQLConn()
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
func (di *DatabaseImpl) MySQLFileRename(fileID int64, newName string) error {
	if strings.Contains(newName, filePathSeparator) {
		return ErrMaliciousRequest
	}

	mysql, err := di.getMySQLConn()
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
func (di *DatabaseImpl) MySQLFileGetInfo(fileID int64) (FileMeta, error) {
	file := FileMeta{}
	mysql, err := di.getMySQLConn()
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
