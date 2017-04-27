package relationalstore

import (
	"database/sql"
	"fmt"
	"time"

	"errors"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql" // required to load into local namespace to
	// initialize sql driver mapping in sql.Open("mysql", ...)
	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datastore"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/go-sql-driver/mysql"
)

func init() {
	datastore.RegisterRelationalStore("mysql", NewMySQLRelationalStore)
}

// MySQLRelationalStore is a concrete implementation of the RelationalStore, using MySQL as the backing relational store of choice
type MySQLRelationalStore struct {
	cfg *config.ConnCfg
	db  *sql.DB
}

// NewMySQLRelationalStore creates a new instance of the MySQLRelationalStore, setting the configuration
func NewMySQLRelationalStore(cfg *config.ConnCfg) datastore.RelationalStore {
	return datastore.RelationalStore(&MySQLRelationalStore{
		cfg: cfg,
	})
}

// Connect starts this relationalStore's connection to the server
// If connect fails, it will throw a fatal error
func (store *MySQLRelationalStore) Connect() {
	if store.db != nil {
		err := store.db.Ping()
		if err != nil {
			store.db = nil
		} else {
			return
		}
	}

	if store.cfg.Schema == "" {
		utils.LogFatal("No MySQL schema found in config", datastore.ErrFatalServerErr, nil)
	}

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&parseTime=true",
		store.cfg.Username,
		store.cfg.Password,
		store.cfg.Host,
		store.cfg.Port,
		store.cfg.Schema,
		store.cfg.Timeout)
	db, err := sql.Open("mysql", connString)
	if err == nil {
		for i := 0; i < store.cfg.NumRetries; i++ {
			if err = db.Ping(); err != nil {
				time.Sleep(time.Duration(store.cfg.Timeout) * time.Second)
			} else {
				store.db = db
				break
			}
		}
	}

	if err != nil {
		store.db = nil
		utils.LogFatal("Unable to connect to MySQL", err, utils.LogFields{
			"Host":   store.cfg.Host,
			"Port":   store.cfg.Port,
			"Schema": store.cfg.Schema,
		})
	}
}

// Shutdown terminates this relationalStore's connection to the server
func (store *MySQLRelationalStore) Shutdown() {
	if store.db != nil {
		err := store.db.Close()
		store.db = nil
		utils.LogError("Failed to close MySQLRelationalStore database connection", err, utils.LogFields{
			"Host":   store.cfg.Host,
			"Port":   store.cfg.Port,
			"Schema": store.cfg.Schema,
		})
	} else {
		utils.LogError("Close called on uninitialized MySQLRelationalStore", datastore.ErrInternalServerErr, nil)
	}
}

func checkResult(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(cols) <= 0 {
		utils.LogError("MySQLRelationalStore: Returned result had no columns", errors.New("No columns returned"), nil)
		panic(errors.New("No columns returned"))
	}

	if cols[0] == "ERROR_CODE" || cols[0] == "ERROR_MSG" {
		var errCode int
		var errMsg string

		for rows.Next() {
			err = rows.Scan(&errCode, &errMsg)
			if err != nil {
				return err
			}
		}

		if errCode != 0 && errMsg != "" {
			utils.LogInfo("MySQLRelationalStore: MySQL returned an error", utils.LogFields{
				"ErrorCode": errCode,
				"ErrorMsg":  errMsg,
			})

			return errors.New(errMsg)
		}
	}

	return nil
}

// UserRegister registers a new user
func (store *MySQLRelationalStore) UserRegister(data *datastore.UserData) error {
	store.Connect()

	rows, err := store.db.Query("CALL user_register(?,?,?,?,?)", data.Username, data.Password, data.Email, data.FirstName, data.LastName)
	if err != nil {
		return err
	}

	err = checkResult(rows)
	if err != nil {
		return err
	}

	return nil
}

// UserLookup returns the user's information (except password)
func (store *MySQLRelationalStore) UserLookup(username string) (*datastore.UserData, error) {
	store.Connect()

	rows, err := store.db.Query("CALL user_lookup(?)", username)
	if err != nil {
		return nil, err
	}

	err = checkResult(rows)
	if err != nil {
		return nil, err
	}

	userData := new(datastore.UserData)
	for rows.Next() {
		err = rows.Scan(&userData.FirstName, &userData.LastName, &userData.Email, &userData.Username)
		if err != nil {
			return nil, err
		}
	}

	return userData, nil
}

// UserGetOwnedProjectIDs returns the ProjectIDs for all project the given user owns
func (store *MySQLRelationalStore) UserGetOwnedProjectIDs(username string) ([]int64, error) {
	store.Connect()

	rows, err := store.db.Query("CALL user_get_projectids(?)", username)
	if err != nil {
		return nil, err
	}

	err = checkResult(rows)
	if err != nil {
		return nil, err
	}

	projectIDs := []int64{}
	for rows.Next() {
		var projectID int64
		err = rows.Scan(&projectID)
		if err != nil {
			return nil, err
		}

		projectIDs = append(projectIDs, projectID)
	}

	return projectIDs, nil
}

// UserGetProjects returns the Project Metadata for all projects the user has permissions for
func (store *MySQLRelationalStore) UserGetProjects(username string) ([]datastore.ProjectMetadata, error) {
	store.Connect()

	rows, err := store.db.Query("CALL user_projects(?)", username)
	if err != nil {
		return nil, err
	}

	err = checkResult(rows)
	if err != nil {
		return nil, err
	}

	projectsMap := map[int64]*datastore.ProjectMetadata{}

	for rows.Next() {
		var projID int64
		var projName string
		var projOwner string
		var permUsername sql.NullString
		var permLevel sql.NullInt64
		var permGrantedBy sql.NullString
		var permGrantedDate mysql.NullTime

		err = rows.Scan(&projID, &projName, &projOwner, &permUsername, &permLevel, &permGrantedBy, &permGrantedDate)
		if err != nil {
			return nil, err
		}

		if projectsMap[projID] != nil {
			project := projectsMap[projID]

			perm := &datastore.ProjectPermission{
				Username:        permUsername.String,
				PermissionLevel: int(permLevel.Int64),
				GrantedBy:       permGrantedBy.String,
				GrantedDate:     permGrantedDate.Time,
			}

			project.ProjectPermissions[permUsername.String] = perm
		} else {
			project := &datastore.ProjectMetadata{
				ProjectID: projID,
				Name:      projName,
				ProjectPermissions: map[string]*datastore.ProjectPermission{
					projOwner: {
						Username:        projOwner,
						PermissionLevel: 10,
						GrantedBy:       projOwner,
						GrantedDate:     time.Unix(0, 0),
					},
				},
			}

			if permUsername.Valid && permLevel.Valid && permGrantedBy.Valid && permGrantedDate.Valid {
				project.ProjectPermissions[permUsername.String] = &datastore.ProjectPermission{
					Username:        permUsername.String,
					PermissionLevel: int(permLevel.Int64),
					GrantedBy:       permGrantedBy.String,
					GrantedDate:     permGrantedDate.Time,
				}
			}

			projectsMap[projID] = project
		}
	}

	projects := []datastore.ProjectMetadata{}

	for _, val := range projectsMap {
		projects = append(projects, *val)
	}

	return projects, nil

}

// UserGetProjectPermissions returns the permissions that a user has for a given project
func (store *MySQLRelationalStore) UserGetProjectPermissions(username string, projectID int64) (int, error) {
	store.Connect()

	rows, err := store.db.Query("CALL user_project_permission(?, ?)", username, projectID)
	if err != nil {
		return 0, err
	}

	err = checkResult(rows)
	if err != nil {
		return 0, err
	}

	var permission int
	for rows.Next() {
		err = rows.Scan(&permission)
		if err != nil {
			return 0, err
		}
	}

	return permission, nil
}

// UserGetPassword retrieves the hash of the user's password
func (store *MySQLRelationalStore) UserGetPassword(username string) (string, error) {
	store.Connect()

	rows, err := store.db.Query("CALL user_get_password(?)", username)
	if err != nil {
		return "", err
	}

	err = checkResult(rows)
	if err != nil {
		return "", err
	}

	var password string
	for rows.Next() {
		err = rows.Scan(&password)
		if err != nil {
			return "", err
		}
	}

	return password, nil
}

// UserDelete deletes the given user
func (store *MySQLRelationalStore) UserDelete(username string) error {
	store.Connect()

	// TODO(wongb): Move this to datahandler
	//rows, err := store.db.Query("Call user_get_projectids(?)", username)
	//
	//var projectIDs []int64
	//for rows.Next() {
	//	projectID := int64(-1)
	//	err = rows.Scan(&projectID)
	//	if err != nil {
	//		return []int64{}, err
	//	}
	//	if projectID == -1 {
	//		return []int64{}, ErrNoData
	//	}
	//	projectIDs = append(projectIDs, projectID)
	//}

	rows, err := store.db.Query("CALL user_delete(?)", username)
	if err != nil {
		return err
	}

	err = checkResult(rows)
	if err != nil {
		return err
	}

	return nil
}

// ProjectCreate creates a new project, assigning it a projectID
func (store *MySQLRelationalStore) ProjectCreate(username string, projectName string) (int64, error) {
	store.Connect()

	rows, err := store.db.Query("CALL project_create(?,?)", projectName, username)
	if err != nil {
		return -1, err
	}

	err = checkResult(rows)
	if err != nil {
		return -1, err
	}

	var projectID int64
	for rows.Next() {
		err = rows.Scan(&projectID)
		if err != nil {
			return -1, err
		}
	}

	return projectID, nil
}

// ProjectLookup returns the project metadata, including permissions
func (store *MySQLRelationalStore) ProjectLookup(projectID int64) (*datastore.ProjectMetadata, error) {
	store.Connect()

	// TODO (optional): un-hardcode '10' as the owner constant in the MySQL ProjectLookup stored proc
	rows, err := store.db.Query("CALL project_lookup(?)", projectID)
	if err != nil {
		return nil, err
	}

	err = checkResult(rows)
	if err != nil {
		return nil, err
	}

	projMeta := &datastore.ProjectMetadata{}
	projMeta.ProjectID = projectID
	for rows.Next() {
		perm := &datastore.ProjectPermission{}
		err = rows.Scan(&projMeta.Name, &perm.Username, &perm.PermissionLevel, &perm.GrantedBy, &perm.GrantedDate)
		if err != nil {
			return nil, err
		}

		if projMeta.ProjectPermissions == nil {
			projMeta.ProjectPermissions = map[string]*datastore.ProjectPermission{}
		}
		projMeta.ProjectPermissions[perm.Username] = perm
	}

	return projMeta, err
}

// ProjectGetFiles returns the list of files that the project contains
// TODO(wongb): MySQL stored proc uses SELECT *; that is bad practice, change to use all columns.
func (store *MySQLRelationalStore) ProjectGetFiles(projectID int64) ([]datastore.FileMetadata, error) {
	store.Connect()

	rows, err := store.db.Query("CALL project_get_files(?)", projectID)
	if err != nil {
		return nil, err
	}

	err = checkResult(rows)
	if err != nil {
		return nil, err
	}

	files := []datastore.FileMetadata{}

	for rows.Next() {
		file := datastore.FileMetadata{}
		err = rows.Scan(&file.FileID, &file.Creator, &file.CreationDate, &file.RelativePath, &file.ProjectID, &file.Filename)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

// ProjectGrantPermissions grants the given permissions to the user with provided username
// TODO(wongb): CHANGE TO ALLOW BULK UPDATES
func (store *MySQLRelationalStore) ProjectGrantPermissions(projectID int64, grantUsername string, permissionLevel int, grantedByUsername string) error {
	store.Connect()

	rows, err := store.db.Query("CALL project_grant_permissions(?, ?, ?, ?)", projectID, grantUsername, permissionLevel, grantedByUsername)
	if err != nil {
		return err
	}

	err = checkResult(rows)
	if err != nil {
		return err
	}

	return nil
}

// ProjectRevokePermissions revokes all permissions from the user with provided username
// TODO(wongb): CHANGE TO ALLOW BULK UPDATES
func (store *MySQLRelationalStore) ProjectRevokePermissions(projectID int64, revokeUsername string) error {
	store.Connect()

	rows, err := store.db.Query("CALL project_revoke_permissions(?, ?)", projectID, revokeUsername)
	if err != nil {
		return err
	}

	err = checkResult(rows)
	if err != nil {
		return err
	}

	return nil
}

// ProjectRename renames the project
func (store *MySQLRelationalStore) ProjectRename(projectID int64, newName string) error {
	store.Connect()

	rows, err := store.db.Query("CALL project_rename(?, ?)", projectID, newName)
	if err != nil {
		return err
	}

	err = checkResult(rows)
	if err != nil {
		return err
	}

	return nil
}

// ProjectDelete deletes a project from MySQL
func (store *MySQLRelationalStore) ProjectDelete(projectID int64) error {
	store.Connect()

	// TODO(wongb): Ensure project_delete checks the user's permission levels before calling this
	rows, err := store.db.Query("CALL project_delete(?)", projectID)
	if err != nil {
		return err
	}

	err = checkResult(rows)
	if err != nil {
		return err
	}

	return nil
}

// FileCreate creates a new file, assigning it a new fileID
func (store *MySQLRelationalStore) FileCreate(username string, projectID int64, filename string, relativePath string) (int64, error) {
	store.Connect()

	filename = filepath.Clean(filename)
	if strings.Contains(filename, datastore.FilePathSeparator) || strings.Contains(filename, "..") {
		return -1, datastore.ErrInvalidFileName
	}

	relativePath = filepath.Clean(relativePath)
	if strings.HasPrefix(relativePath, "..") {
		return -1, datastore.ErrInvalidFilePath
	}

	rows, err := store.db.Query("CALL file_create(?,?,?,?)", username, filename, relativePath, projectID)
	if err != nil {
		return -1, err
	}

	err = checkResult(rows)
	if err != nil {
		return -1, err
	}

	var fileID int64
	for rows.Next() {
		err = rows.Scan(&fileID)
		if err != nil {
			return -1, err
		}
	}

	return fileID, nil
}

// FileGet returns the metadata for the given fileID
func (store *MySQLRelationalStore) FileGet(fileID int64) (*datastore.FileMetadata, error) {
	store.Connect()

	rows, err := store.db.Query("CALL file_get_info(?)", fileID)
	if err != nil {
		return nil, err
	}

	err = checkResult(rows)
	if err != nil {
		return nil, err
	}

	fileMeta := &datastore.FileMetadata{
		FileID: fileID,
	}
	for rows.Next() {
		err = rows.Scan(&fileMeta.Creator, &fileMeta.CreationDate, &fileMeta.RelativePath, &fileMeta.ProjectID, &fileMeta.Filename)
		if err != nil {
			return nil, err
		}
	}

	return fileMeta, nil
}

// FileMove updates the filepath for the given fileID
func (store *MySQLRelationalStore) FileMove(fileID int64, newRelativePath string) error {
	store.Connect()

	newPathClean := filepath.Clean(newRelativePath)
	if strings.HasPrefix(newPathClean, "..") {
		return datastore.ErrInvalidFilePath
	}

	rows, err := store.db.Query("CALL file_move(?, ?)", fileID, newPathClean)
	if err != nil {
		return err
	}

	err = checkResult(rows)
	if err != nil {
		return err
	}

	return nil
}

// FileRename updates the filename for the given fileID
func (store *MySQLRelationalStore) FileRename(fileID int64, newName string) error {
	store.Connect()

	newNameClean := filepath.Clean(newName)
	if strings.Contains(newNameClean, datastore.FilePathSeparator) || strings.Contains(newNameClean, "..") {
		return datastore.ErrInvalidFileName
	}

	rows, err := store.db.Query("CALL file_rename(?, ?)", fileID, newNameClean)
	if err != nil {
		return err
	}

	err = checkResult(rows)
	if err != nil {
		return err
	}

	return nil
}

// FileDelete deletes the stored metadata for the given fileID
func (store *MySQLRelationalStore) FileDelete(fileID int64) error {
	store.Connect()

	rows, err := store.db.Query("CALL file_delete(?)", fileID)
	if err != nil {
		return err
	}

	err = checkResult(rows)
	if err != nil {
		return err
	}

	return nil
}
