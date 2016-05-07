package dbfs

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // required to load into local namespace to
	// initialize sql driver mapping in sql.Open("mysql", ...)
	"strconv"

	"time"

	"os"
	"path/filepath"
	"strings"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/utils"
)

var mysqldb *mysqlConn

var mysqldbName = "cc"

type mysqlConn struct {
	config config.ConnCfg
	db     *sql.DB
	dbname string
}

func openMySQLConn(dbName string) (*mysqlConn, error) {
	if mysqldb != nil && mysqldb.db != nil {
		err := mysqldb.db.Ping()
		if err == nil && mysqldb.dbname == dbName {
			return mysqldb, nil
		}
	}

	if mysqldb == nil || mysqldb.config == (config.ConnCfg{}) {
		mysqldb = new(mysqlConn)
		configMap := config.GetConfig()
		mysqldb.config = configMap.ConnectionConfig["MySQL"]
	}

	db, err := sql.Open("mysql", mysqldb.config.Username+":"+mysqldb.config.Password+"@tcp("+mysqldb.config.Host+":"+strconv.Itoa(int(mysqldb.config.Port))+")/"+dbName+"?timeout="+strconv.Itoa(int(mysqldb.config.Timeout))+"s"+"&parseTime=true")
	if err != nil {
		utils.LogOnError(err, "Unable to connect to MySQL")
		return mysqldb, err
	}

	mysqldb.db = db
	mysqldb.dbname = dbName
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
func MySQLUserRegister(username string, pass string, email string, firstName string, lastName string) error {
	mysql, err := openMySQLConn(mysqldbName)
	if err != nil {
		return err
	}

	result, err := mysql.db.Exec("CALL user_register(?,?,?,?,?)", username, pass, email, firstName, lastName)
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
	mysql, err := openMySQLConn(mysqldbName)
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
	mysql, err := openMySQLConn(mysqldbName)
	if err != nil {
		return err
	}

	// FIXME: use MySQLUserGetPass to verify the password is correct before deleting

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
func MySQLUserLookup(username string) (firstname string, lastname string, email string, err error) {
	mysql, err := openMySQLConn(mysqldbName)
	if err != nil {
		return "", "", "", err
	}

	rows, err := mysql.db.Query("CALL user_lookup(?)", username)
	if err != nil {
		return "", "", "", err
	}

	for rows.Next() {
		var username string
		err = rows.Scan(&firstname, &lastname, &email, &username)
		if err != nil {
			return "", "", "", err
		}
	}

	return firstname, lastname, email, nil
}

// MySQLUserProjects returns the projectID, the project name, and the permission level the user `username` has on that project
func MySQLUserProjects(username string) (projects []Project, err error) {
	mysql, err := openMySQLConn(mysqldbName)
	if err != nil {
		return nil, err
	}

	rows, err := mysql.db.Query("CALL user_projects(?)", username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		project := Project{}
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
	mysql, err := openMySQLConn(mysqldbName)
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
func MySQLProjectDelete(projectID int64) error {
	mysql, err := openMySQLConn(mysqldbName)
	if err != nil {
		return err
	}

	result, err := mysql.db.Exec("CALL project_delete(?)", projectID)
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
	mysql, err := openMySQLConn(mysqldbName)
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
	mysql, err := openMySQLConn(mysqldbName)
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
func MySQLProjectRevokePermission(projectID int64, revokeUsername string) error {
	mysql, err := openMySQLConn(mysqldbName)
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
	mysql, err := openMySQLConn(mysqldbName)
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
func MySQLProjectLookup(projectID int64) (name string, permisions map[string]ProjectPermission, err error) {
	permisions = make(map[string](ProjectPermission))
	mysql, err := openMySQLConn(mysqldbName)
	if err != nil {
		return "", permisions, err
	}

	rows, err := mysql.db.Query("CALL project_lookup(?)", projectID)
	if err != nil {
		return "", permisions, err
	}

	for rows.Next() {
		perm := ProjectPermission{}
		var tiempo string
		err = rows.Scan(&name, &perm.Username, &perm.PermissionLevel, &perm.GrantedBy, &tiempo)
		perm.GrantedDate, _ = time.Parse("2006-01-02 15:04:05", tiempo)
		if err != nil {
			return "", permisions, err
		}
		permisions[perm.Username] = perm
	}

	return name, permisions, err
}

// MySQLFileCreate create a new file in MySQL
func MySQLFileCreate(username string, filename string, relativePath string, projectID int64) (fileID int64, err error) {
	pathsep := strconv.QuoteRune(os.PathSeparator)
	if strings.Contains(filename, pathsep[1:len(pathsep)-1]) {
		return -1, ErrMalliciousRequest
	}

	filepathClean := filepath.Clean(relativePath)
	if strings.HasPrefix(filepathClean, "..") {
		return -1, ErrMalliciousRequest
	}

	mysql, err := openMySQLConn(mysqldbName)
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
	mysql, err := openMySQLConn(mysqldbName)
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
		return ErrMalliciousRequest
	}

	mysql, err := openMySQLConn(mysqldbName)
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
	pathsep := strconv.QuoteRune(os.PathSeparator)
	if strings.Contains(newName, pathsep[1:len(pathsep)-1]) {
		return ErrMalliciousRequest
	}

	mysql, err := openMySQLConn(mysqldbName)
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
