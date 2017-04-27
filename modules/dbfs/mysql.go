package dbfs

import (
	"errors"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datastore"
	"github.com/CodeCollaborate/Server/modules/datastore/relationalstore"
)

type mysqlConn struct {
	config               config.ConnCfg
	mysqlRelationalStore datastore.RelationalStore
}

func (di *DatabaseImpl) openMySQL() {
	if di.mysqldb == nil || di.mysqldb.config == (config.ConnCfg{}) {
		di.mysqldb = new(mysqlConn)
		configMap := config.GetConfig()
		di.mysqldb.config = *configMap.DataStoreConfig.RelationalStoreCfg
	}

	if di.mysqldb.mysqlRelationalStore == nil {
		di.mysqldb.mysqlRelationalStore = relationalstore.NewMySQLRelationalStore(&di.mysqldb.config)
		di.mysqldb.mysqlRelationalStore.Connect()
	}

}

// CloseMySQL closes the MySQL db connection
// YOU PROBABLY DON'T NEED TO RUN THIS EVER
func (di *DatabaseImpl) CloseMySQL() error {
	di.mysqldb.mysqlRelationalStore.Shutdown()

	return nil
}

/**
STORED PROCEDURES
*/

// MySQLUserRegister registers a new user in MySQL
func (di *DatabaseImpl) MySQLUserRegister(user UserMeta) error {
	di.openMySQL()

	userData := &datastore.UserData{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}

	return di.mysqldb.mysqlRelationalStore.UserRegister(userData)
}

// MySQLUserGetPass is used to get the key and hash of a stored password to verify that a value is correct
func (di *DatabaseImpl) MySQLUserGetPass(username string) (password string, err error) {
	di.openMySQL()

	return di.mysqldb.mysqlRelationalStore.UserGetPassword(username)
}

// MySQLUserDelete deletes a user from MySQL
func (di *DatabaseImpl) MySQLUserDelete(username string) ([]int64, error) {
	di.openMySQL()

	projectIDs, err := di.mysqldb.mysqlRelationalStore.UserGetOwnedProjectIDs(username)
	if err != nil {
		return nil, err
	}

	err = di.mysqldb.mysqlRelationalStore.UserDelete(username)
	if err != nil {
		return nil, err
	}

	return projectIDs, nil
}

// MySQLUserLookup returns user information about a user with the username 'username'
func (di *DatabaseImpl) MySQLUserLookup(username string) (user UserMeta, err error) {
	di.openMySQL()

	userData, err := di.mysqldb.mysqlRelationalStore.UserLookup(username)
	if err != nil {
		return user, err
	}

	user = UserMeta{
		Username:  userData.Username,
		Email:     userData.Email,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Password:  userData.Password,
	}

	return user, nil
}

// MySQLUserProjects returns the projectID, the project name, and the permission level the user `username` has on that project
func (di *DatabaseImpl) MySQLUserProjects(username string) ([]ProjectMeta, error) {
	di.openMySQL()

	projMetas, err := di.mysqldb.mysqlRelationalStore.UserGetProjects(username)
	if err != nil {
		return nil, err
	}

	result := []ProjectMeta{}
	for _, projMeta := range projMetas {
		permissionLevel := 10
		if projMeta.ProjectPermissions[username] != nil {
			permissionLevel = projMeta.ProjectPermissions[username].PermissionLevel
		} else {
			return nil, errors.New("No entry in permissions map for owner")
		}

		newMeta := ProjectMeta{
			ProjectID:       projMeta.ProjectID,
			Name:            projMeta.Name,
			PermissionLevel: permissionLevel,
		}

		result = append(result, newMeta)
	}
	return result, nil
}

// MySQLProjectCreate create a new project in MySQL
func (di *DatabaseImpl) MySQLProjectCreate(username string, projectName string) (projectID int64, err error) {
	di.openMySQL()

	return di.mysqldb.mysqlRelationalStore.ProjectCreate(username, projectName)
}

// MySQLProjectDelete deletes a project from MySQL
func (di *DatabaseImpl) MySQLProjectDelete(projectID int64, senderID string) error {
	di.openMySQL()

	return di.mysqldb.mysqlRelationalStore.ProjectDelete(projectID)
}

// MySQLProjectGetFiles returns the Files from the project with projectID = projectID
func (di *DatabaseImpl) MySQLProjectGetFiles(projectID int64) (files []FileMeta, err error) {
	di.openMySQL()

	filesMetadata, err := di.mysqldb.mysqlRelationalStore.ProjectGetFiles(projectID)
	if err != nil {
		return nil, err
	}

	result := []FileMeta{}
	for _, fileMeta := range filesMetadata {
		newMeta := FileMeta{
			ProjectID:    fileMeta.ProjectID,
			Creator:      fileMeta.Creator,
			RelativePath: fileMeta.RelativePath,
			Filename:     fileMeta.Filename,
			FileID:       fileMeta.FileID,
			CreationDate: fileMeta.CreationDate,
		}

		result = append(result, newMeta)
	}
	return result, nil
}

// MySQLProjectGrantPermission gives the user `grantUsername` the permission `permissionLevel` on project `projectID`
func (di *DatabaseImpl) MySQLProjectGrantPermission(projectID int64, grantUsername string, permissionLevel int, grantedByUsername string) error {
	di.openMySQL()

	return di.mysqldb.mysqlRelationalStore.ProjectGrantPermissions(projectID, grantUsername, permissionLevel, grantedByUsername)
}

// MySQLProjectRevokePermission removes revokeUsername's permissions from the project
// DOES NOT WORK FOR OWNER (which is kinda a good thing)
func (di *DatabaseImpl) MySQLProjectRevokePermission(projectID int64, revokeUsername string, revokedByUsername string) error {
	di.openMySQL()

	return di.mysqldb.mysqlRelationalStore.ProjectRevokePermissions(projectID, revokeUsername)
}

// MySQLUserProjectPermissionLookup returns the permission level of `username` on the project with the given projectID
func (di *DatabaseImpl) MySQLUserProjectPermissionLookup(projectID int64, username string) (int, error) {
	di.openMySQL()

	return di.mysqldb.mysqlRelationalStore.UserGetProjectPermissions(username, projectID)
}

// MySQLProjectRename allows for you to rename projects
func (di *DatabaseImpl) MySQLProjectRename(projectID int64, newName string) error {
	di.openMySQL()

	return di.mysqldb.mysqlRelationalStore.ProjectRename(projectID, newName)
}

// MySQLProjectLookup returns the project name and permissions for a project with ProjectID = 'projectID'
func (di *DatabaseImpl) MySQLProjectLookup(projectID int64, username string) (name string, permissions map[string]ProjectPermission, err error) {
	di.openMySQL()

	projMeta, err := di.mysqldb.mysqlRelationalStore.ProjectLookup(projectID)
	if err != nil {
		return "", nil, err
	}

	permMap := map[string]ProjectPermission{}

	for key, value := range projMeta.ProjectPermissions {
		permMap[key] = ProjectPermission{
			Username:        value.Username,
			PermissionLevel: value.PermissionLevel,
			GrantedBy:       value.GrantedBy,
			GrantedDate:     value.GrantedDate,
		}
	}

	return projMeta.Name, permMap, nil
}

// MySQLFileCreate create a new file in MySQL
func (di *DatabaseImpl) MySQLFileCreate(username string, filename string, relativePath string, projectID int64) (int64, error) {
	di.openMySQL()

	return di.mysqldb.mysqlRelationalStore.FileCreate(username, projectID, filename, relativePath)
}

// MySQLFileDelete deletes a file from the MySQL database
// this does not delete the actual file
func (di *DatabaseImpl) MySQLFileDelete(fileID int64) error {
	di.openMySQL()

	return di.mysqldb.mysqlRelationalStore.FileDelete(fileID)
}

// MySQLFileMove updates MySQL with the  new path of the file with FileID == 'fileID'
func (di *DatabaseImpl) MySQLFileMove(fileID int64, newPath string) error {
	di.openMySQL()

	return di.mysqldb.mysqlRelationalStore.FileMove(fileID, newPath)
}

// MySQLFileRename updates MySQL with the new name of the file with FileID == 'fileID'
func (di *DatabaseImpl) MySQLFileRename(fileID int64, newName string) error {
	di.openMySQL()

	return di.mysqldb.mysqlRelationalStore.FileRename(fileID, newName)
}

// MySQLFileGetInfo returns the meta data about the given file
func (di *DatabaseImpl) MySQLFileGetInfo(fileID int64) (FileMeta, error) {
	di.openMySQL()

	fileData, err := di.mysqldb.mysqlRelationalStore.FileGet(fileID)
	if err != nil {
		return FileMeta{}, err
	}

	fileMeta := FileMeta{
		FileID:       fileData.FileID,
		Creator:      fileData.Creator,
		CreationDate: fileData.CreationDate,
		RelativePath: fileData.RelativePath,
		Filename:     fileData.Filename,
		ProjectID:    fileData.ProjectID,
	}

	return fileMeta, nil
}
