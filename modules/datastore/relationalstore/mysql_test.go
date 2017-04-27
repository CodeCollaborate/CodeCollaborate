package relationalstore

import (
	"reflect"
	"testing"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var userOne = &datastore.UserData{
	Username:  "_test_user1",
	Password:  "_test_password1",
	Email:     "_test_email1@codecollab.cc",
	FirstName: "_thing",
	LastName:  "_one",
}
var projectOneName = "testProjectOne"
var fileOneName = "testFileOne"

var userTwo = &datastore.UserData{
	Username:  "_test_user2",
	Password:  "_test_password2",
	Email:     "_test_email2@codecollab.cc",
	FirstName: "_thing",
	LastName:  "_two",
}

func TestMySQLRelationalStore_RegisterSelf(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	documentStore := datastore.InitRelationalStore("mysql", config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	require.True(t, reflect.TypeOf(documentStore).String() == reflect.TypeOf(&MySQLRelationalStore{}).String(), "relationalStore initalized wrong type")
}

func TestMySQLRelationalStore_Connect(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	cb.UserDelete(userOne.Username)
	cb.UserDelete(userTwo.Username)
}

func TestMySQLRelationalStore_Shutdown(t *testing.T) {
	// Nothing to test
}

func TestMySQLRelationalStore_UserRegister(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	err = cb.UserRegister(userOne)
	require.NotNil(t, err, "Failed to throw error registering duplicate user")
}

func TestMySQLRelationalStore_UserLookup(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	userRet, err := cb.UserLookup(userOne.Username)
	require.Nil(t, err, "MySQL threw an error when looking up user")

	if userRet.FirstName != userOne.FirstName || userRet.LastName != userOne.LastName || userRet.Email != userOne.Email {
		t.Fatalf("Wrong return, got: %v %v, email: %v", userRet.FirstName, userRet.LastName, userRet.Email)
	}

	_, err = cb.UserLookup("unknownUsername")
	require.EqualError(t, err, "No such username found", "MySQL did not throw expected error for unknown username")
}

func TestMySQLRelationalStore_UserGetProjects(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	projects, err := cb.UserGetProjects(userOne.Username)
	require.Nil(t, err, "MySQL threw an error when retrieving user projects")
	require.Len(t, projects, 1, "MySQL returned an incorrect number of projects")
	require.Equal(t, projectID, projects[0].ProjectID, "MySQL returned project with incorrect projectID")
	require.Equal(t, projectOneName, projects[0].Name, "MySQL returned project with incorrect name")
}

func TestMySQLRelationalStore_UserGetOwnedProjectIDs(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	projectIDs, err := cb.UserGetOwnedProjectIDs(userOne.Username)
	require.Nil(t, err, "MySQL threw an error when retrieving user projects")
	require.Len(t, projectIDs, 1, "MySQL returned an incorrect number of projects")
	require.Equal(t, projectID, projectIDs[0], "MySQL returned project with incorrect projectID")
}

func TestMySQLRelationalStore_UserGetProjectPermissions(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	defer func() {
		cb.UserDelete(userOne.Username)
		cb.UserDelete(userTwo.Username)
	}()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering userOne")

	err = cb.UserRegister(userTwo)
	require.Nil(t, err, "MySQL threw an error when registering userTwo")

	projectID, _ := cb.ProjectCreate(userOne.Username, projectOneName)

	permLevel, err := cb.UserGetProjectPermissions(userOne.Username, projectID)
	assert.Nil(t, err, "MySQL threw an error while looking for Project Permissions")
	ownerPerm, _ := config.PermissionByLabel("owner")
	assert.Equal(t, ownerPerm.Level, permLevel, "Expected userOne to be owner")

	permLevel, err = cb.UserGetProjectPermissions(userTwo.Username, projectID)
	assert.Nil(t, err, "MySQL threw an error while looking for Project Permissions")
	assert.Equal(t, 0, permLevel, "UserTwo had permissions before granting")

	readPerm, _ := config.PermissionByLabel("read")
	err = cb.ProjectGrantPermissions(projectID, userTwo.Username, readPerm.Level, userOne.Username)
	assert.Nil(t, err)

	permLevel, err = cb.UserGetProjectPermissions(userTwo.Username, projectID)
	assert.Nil(t, err, "MySQL threw an error while looking for Project Permissions")
	assert.Equal(t, readPerm.Level, permLevel, "UserTwo did not have expected permissions")
}

func TestMySQLRelationalStore_UserGetPassword(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	pass, err := cb.UserGetPassword(userOne.Username)
	require.Nil(t, err, "MySQL threw error while attempting to get password")

	require.Equal(t, pass, userOne.Password, "Wrong password returned")
}

func TestMySQLRelationalStore_UserDelete(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user one")
	defer cb.UserDelete(userOne.Username)

	err = cb.UserRegister(userTwo)
	require.Nil(t, err, "MySQL threw an error when registering user two")
	defer cb.UserDelete(userTwo.Username)

	err = cb.UserDelete(userOne.Username)
	assert.Nil(t, err, "Error encountered while deleting user one")

	// check user actually deleted
	_, err = cb.UserLookup(userOne.Username)
	assert.EqualError(t, err, "No such username found", "expected no user to be returned")

	// test with projects for notifications
	err = cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user one")

	projectID1, err := cb.ProjectCreate(userOne.Username, "_test_project_1")
	projectID2, err := cb.ProjectCreate(userOne.Username, "_test_project_2")

	writePerm, err := config.PermissionByLabel("write")
	assert.Nil(t, err, "api permissions error")

	err = cb.ProjectGrantPermissions(projectID1, userTwo.Username, writePerm.Level, userOne.Username)
	assert.NoError(t, err, "project grant permission error")
	err = cb.ProjectGrantPermissions(projectID2, userTwo.Username, writePerm.Level, userOne.Username)
	assert.NoError(t, err, "project grant permission error")

	err = cb.UserDelete(userOne.Username)
	assert.Nil(t, err, "Error encountered while deleting user")

	// check user actually deleted
	_, err = cb.UserLookup(userOne.Username)
	assert.EqualError(t, err, "No such username found", "expected no user to be returned")

	// check projects actually deleted
	_, err = cb.ProjectLookup(projectID1)
	assert.EqualError(t, err, "No such projectID found", "expected project1 to not exist")

	_, err = cb.ProjectLookup(projectID2)
	assert.EqualError(t, err, "No such projectID found", "expected project2 to not exist")
}

func TestMySQLRelationalStore_ProjectCreate(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	_, err = cb.ProjectCreate(userOne.Username, projectOneName)
	require.EqualError(t, err, "Owner already has a project with the given name", "MySQL threw incorrect error on duplicate owner/project")
}

func TestMySQLRelationalStore_ProjectLookup(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user one")
	defer cb.UserDelete(userOne.Username)

	err = cb.UserRegister(userTwo)
	require.Nil(t, err, "MySQL threw an error when registering user two")
	defer cb.UserDelete(userTwo.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	projMeta, err := cb.ProjectLookup(projectID)
	require.Nil(t, err, "MySQL threw an error when looking up a project")
	require.Equal(t, projectID, projMeta.ProjectID, "MySQL returned an incorrect ID for the given project")
	require.Equal(t, projectOneName, projMeta.Name, "MySQL returned an incorrect Name for the given project")
	require.Len(t, projMeta.ProjectPermissions, 1, "MySQL returned a permissions array of incorrect size for the given project")
	require.Equal(t, userOne.Username, projMeta.ProjectPermissions[userOne.Username].Username, "MySQL returned permission entry with incorrect username for userOne")
	require.Equal(t, 10, projMeta.ProjectPermissions[userOne.Username].PermissionLevel, "MySQL returned permission entry with incorrect username for userOne")
	require.Equal(t, userOne.Username, projMeta.ProjectPermissions[userOne.Username].GrantedBy, "MySQL returned permission entry with incorrect username for userOne")

	err = cb.ProjectGrantPermissions(projectID, userTwo.Username, 5, userOne.Username)
	require.Nil(t, err, "MySQL threw an error when granting permissions to user two")

	projMeta, err = cb.ProjectLookup(projectID)
	require.Nil(t, err, "MySQL threw an error when looking up a project")
	require.Len(t, projMeta.ProjectPermissions, 2, "MySQL returned permissions map for project with incorrect length")
	require.Equal(t, userOne.Username, projMeta.ProjectPermissions[userOne.Username].Username, "MySQL returned permission entry with incorrect username for userOne")
	require.Equal(t, 10, projMeta.ProjectPermissions[userOne.Username].PermissionLevel, "MySQL returned permission entry with incorrect username for userOne")
	require.Equal(t, userOne.Username, projMeta.ProjectPermissions[userOne.Username].GrantedBy, "MySQL returned permission entry with incorrect username for userOne")
	require.Equal(t, userTwo.Username, projMeta.ProjectPermissions[userTwo.Username].Username, "MySQL returned permission entry with incorrect username for userTwo")
	require.Equal(t, 5, projMeta.ProjectPermissions[userTwo.Username].PermissionLevel, "MySQL returned permission entry with incorrect username for userTwo")
	require.Equal(t, userOne.Username, projMeta.ProjectPermissions[userTwo.Username].GrantedBy, "MySQL returned permission entry with incorrect username for userTwo")
}

func TestMySQLRelationalStore_ProjectGetFiles(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	fileID, err := cb.FileCreate(userOne.Username, projectID, fileOneName, ".")
	require.Nil(t, err, "MySQL threw an error when creating a new file")

	files, err := cb.ProjectGetFiles(projectID)
	require.Nil(t, err, "MySQL threw an error when attempting to get files in project")
	require.Len(t, files, 1, "MySQL returned files array of incorrect length")

	require.Equal(t, fileID, files[0].FileID, "MySQL returned file[0] with incorrect FileID")
	require.Equal(t, fileOneName, files[0].Filename, "MySQL returned file[0] with incorrect Filename")
	require.Equal(t, ".", files[0].RelativePath, "MySQL returned file[0] with incorrect RelativePath")
	require.Equal(t, projectID, files[0].ProjectID, "MySQL returned file[0] with incorrect ProjectID")
	require.Equal(t, userOne.Username, files[0].Creator, "MySQL returned file[0] with incorrect Creator")
}

func TestMySQLRelationalStore_ProjectGrantPermissions(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user one")
	defer cb.UserDelete(userOne.Username)

	err = cb.UserRegister(userTwo)
	require.Nil(t, err, "MySQL threw an error when registering user two")
	defer cb.UserDelete(userTwo.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	err = cb.ProjectGrantPermissions(projectID, userTwo.Username, 5, userOne.Username)
	require.Nil(t, err, "MySQL threw an error when granting permissions to user two")

	projects, err := cb.UserGetProjects(userTwo.Username)
	require.Nil(t, err, "MySQL threw an error when retrieving projects for user two")

	require.Len(t, projects, 1, "MySQL failed to grant permissions, returned an incorrect number of projects for user two")
	require.Equal(t, projectID, projects[0].ProjectID, "MySQL returned project with incorrect projectID")
	require.Equal(t, projectOneName, projects[0].Name, "MySQL returned project with incorrect name")

	require.Len(t, projects[0].ProjectPermissions, 2, "MySQL returned permissions map for project with incorrect length")
	require.Equal(t, userOne.Username, projects[0].ProjectPermissions[userOne.Username].Username, "MySQL returned permission entry with incorrect username for userOne")
	require.Equal(t, 10, projects[0].ProjectPermissions[userOne.Username].PermissionLevel, "MySQL returned permission entry with incorrect username for userOne")
	require.Equal(t, userOne.Username, projects[0].ProjectPermissions[userOne.Username].GrantedBy, "MySQL returned permission entry with incorrect username for userOne")
	require.Equal(t, userTwo.Username, projects[0].ProjectPermissions[userTwo.Username].Username, "MySQL returned permission entry with incorrect username for userTwo")
	require.Equal(t, 5, projects[0].ProjectPermissions[userTwo.Username].PermissionLevel, "MySQL returned permission entry with incorrect username for userTwo")
	require.Equal(t, userOne.Username, projects[0].ProjectPermissions[userTwo.Username].GrantedBy, "MySQL returned permission entry with incorrect username for userTwo")
}

func TestMySQLRelationalStore_ProjectRevokePermissions(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user one")
	defer cb.UserDelete(userOne.Username)

	err = cb.UserRegister(userTwo)
	require.Nil(t, err, "MySQL threw an error when registering user two")
	defer cb.UserDelete(userTwo.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	err = cb.ProjectGrantPermissions(projectID, userTwo.Username, 5, userOne.Username)
	require.Nil(t, err, "MySQL threw an error when granting permissions to user two")

	projects, err := cb.UserGetProjects(userTwo.Username)
	require.Nil(t, err, "MySQL threw an error when retrieving projects for user two")
	require.Len(t, projects, 1, "MySQL failed to grant permissions, returned an incorrect number of projects for user two")
	require.Equal(t, projectID, projects[0].ProjectID, "MySQL returned an incorrect project for user two")

	cb.ProjectRevokePermissions(projectID, userTwo.Username)

	projects, err = cb.UserGetProjects(userTwo.Username)
	require.Len(t, projects, 0, "MySQL failed to revoke permissions, returned incorrect number of projects for user two")
}

func TestMySQLRelationalStore_ProjectRename(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	err = cb.ProjectRename(projectID, projectOneName+"_renamed")
	require.Nil(t, err, "MySQL threw an error when renaming the project")

	projects, err := cb.UserGetProjects(userOne.Username)
	require.Nil(t, err, "MySQL threw an error when retrieving projects for user two")
	require.Len(t, projects, 1, "MySQL returned an incorrect number of projects")
	require.Equal(t, projectID, projects[0].ProjectID, "MySQL returned a project with incorrect projectID")
	require.Equal(t, projectOneName+"_renamed", projects[0].Name, "MySQL returned a project with incorrect Name")
	require.Len(t, projects[0].ProjectPermissions, 1, "MySQL returned a project with incorrect permissions map")
}

func TestMySQLRelationalStore_ProjectDelete(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	// test trying to delete a project that contains files
	_, err = cb.FileCreate(userOne.Username, projectID, fileOneName, ".")
	require.Nil(t, err, "MySQL threw an error when creating a new file")

	err = cb.ProjectDelete(projectID)
	require.Nil(t, err, "MySQL threw an error when deleting a project")

	err = cb.ProjectDelete(projectID)
	require.NotNil(t, err, "No such projectID found", "MySQL threw an incorrect error when deleting a nonexistent project")
}

func TestMySQLRelationalStore_FileCreate(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	fileID, err := cb.FileCreate(userOne.Username, projectID, fileOneName, ".")
	require.Nil(t, err, "MySQL threw an error when creating a new file")
	require.True(t, fileID >= 0, "MySQL returned an invalid FileID")

	files, err := cb.ProjectGetFiles(projectID)
	require.Nil(t, err, "MySQL threw an error when getting files for given project")
	require.Len(t, files, 1, "MySQL returned incorrect number of files for given project")
	assert.Equal(t, fileID, files[0].FileID, "MySQL returned file with incorrect FileID")
	assert.Equal(t, userOne.Username, files[0].Creator, "MySQL returned file with incorrect Creator")
	assert.Equal(t, ".", files[0].RelativePath, "MySQL returned file with incorrect RelativePath")
	assert.Equal(t, fileOneName, files[0].Filename, "MySQL returned file with incorrect Filename")
	assert.Equal(t, projectID, files[0].ProjectID, "MySQL returned file with incorrect ProjectID")

	fileID, err = cb.FileCreate(userOne.Username, projectID, fileOneName, ".")
	require.EqualError(t, err, "Project already contains file at the given location", "MySQL threw incorrect error for duplicate file")
}

func TestMySQLRelationalStore_FileGet(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	fileID, err := cb.FileCreate(userOne.Username, projectID, fileOneName, ".")
	require.Nil(t, err, "MySQL threw an error when creating a new file")
	require.True(t, fileID >= 0, "MySQL returned an invalid FileID")

	fileMeta, err := cb.FileGet(fileID)
	require.Nil(t, err, "MySQL threw an error when getting file metadata")
	assert.Equal(t, fileID, fileMeta.FileID, "MySQL returned file with incorrect FileID")
	assert.Equal(t, userOne.Username, fileMeta.Creator, "MySQL returned file with incorrect Creator")
	assert.Equal(t, ".", fileMeta.RelativePath, "MySQL returned file with incorrect RelativePath")
	assert.Equal(t, fileOneName, fileMeta.Filename, "MySQL returned file with incorrect Filename")
	assert.Equal(t, projectID, fileMeta.ProjectID, "MySQL returned file with incorrect ProjectID")

	err = cb.FileMove(fileID, "cc")
	require.Nil(t, err, "MySQL threw an error when renaming file")

	fileMeta, err = cb.FileGet(fileID)
	require.Nil(t, err, "MySQL threw an error when getting file metadata")
	assert.Equal(t, fileID, fileMeta.FileID, "MySQL returned file with incorrect FileID")
	assert.Equal(t, userOne.Username, fileMeta.Creator, "MySQL returned file with incorrect Creator")
	assert.Equal(t, "cc", fileMeta.RelativePath, "MySQL returned file with incorrect RelativePath")
	assert.Equal(t, fileOneName, fileMeta.Filename, "MySQL returned file with incorrect Filename")
	assert.Equal(t, projectID, fileMeta.ProjectID, "MySQL returned file with incorrect ProjectID")
}

func TestMySQLRelationalStore_FileMove(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	fileID, err := cb.FileCreate(userOne.Username, projectID, fileOneName, ".")
	require.Nil(t, err, "MySQL threw an error when creating a new file")
	require.True(t, fileID >= 0, "MySQL returned an invalid FileID")

	files, err := cb.ProjectGetFiles(projectID)
	require.Nil(t, err, "MySQL threw an error when getting files for given project")
	require.Len(t, files, 1, "MySQL returned incorrect number of files for given project")
	assert.Equal(t, fileID, files[0].FileID, "MySQL returned file with incorrect FileID")
	assert.Equal(t, ".", files[0].RelativePath, "MySQL returned file with incorrect RelativePath")
	assert.Equal(t, fileOneName, files[0].Filename, "MySQL returned file with incorrect Filename")

	err = cb.FileMove(fileID, "cc")
	require.Nil(t, err, "MySQL threw an error when renaming file")

	files, err = cb.ProjectGetFiles(projectID)
	require.Nil(t, err, "MySQL threw an error when getting files for given project")
	require.Len(t, files, 1, "MySQL returned incorrect number of files for given project")
	assert.Equal(t, fileID, files[0].FileID, "MySQL returned file with incorrect FileID")
	assert.Equal(t, "cc", files[0].RelativePath, "MySQL returned file with incorrect RelativePath")
	assert.Equal(t, fileOneName, files[0].Filename, "MySQL returned file with incorrect Filename")
}

func TestMySQLRelationalStore_FileRename(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	fileID, err := cb.FileCreate(userOne.Username, projectID, fileOneName, ".")
	require.Nil(t, err, "MySQL threw an error when creating a new file")
	require.True(t, fileID >= 0, "MySQL returned an invalid FileID")

	files, err := cb.ProjectGetFiles(projectID)
	require.Nil(t, err, "MySQL threw an error when getting files for given project")
	require.Len(t, files, 1, "MySQL returned incorrect number of files for given project")
	assert.Equal(t, fileID, files[0].FileID, "MySQL returned file with incorrect FileID")
	assert.Equal(t, ".", files[0].RelativePath, "MySQL returned file with incorrect RelativePath")
	assert.Equal(t, fileOneName, files[0].Filename, "MySQL returned file with incorrect Filename")

	err = cb.FileRename(fileID, fileOneName+"_renamed")
	require.Nil(t, err, "MySQL threw an error when renaming file")

	files, err = cb.ProjectGetFiles(projectID)
	require.Nil(t, err, "MySQL threw an error when getting files for given project")
	require.Len(t, files, 1, "MySQL returned incorrect number of files for given project")
	assert.Equal(t, fileID, files[0].FileID, "MySQL returned file with incorrect FileID")
	assert.Equal(t, ".", files[0].RelativePath, "MySQL returned file with incorrect RelativePath")
	assert.Equal(t, fileOneName+"_renamed", files[0].Filename, "MySQL returned file with incorrect Filename")
}

func TestMySQLRelationalStore_FileDelete(t *testing.T) {
	config.SetupTestingConfig(t, "../../../config")
	cb := NewMySQLRelationalStore(config.GetConfig().DataStoreConfig.RelationalStoreCfg)

	defer func() {
		require.Nil(t, recover(), "MySQL connect threw a fatal error")
	}()
	cb.Connect()

	err := cb.UserRegister(userOne)
	require.Nil(t, err, "MySQL threw an error when registering user")
	defer cb.UserDelete(userOne.Username)

	projectID, err := cb.ProjectCreate(userOne.Username, projectOneName)
	require.Nil(t, err, "MySQL threw an error when creating a new project")
	require.True(t, projectID >= 0, "MySQL returned an invalid ProjectID")

	fileID, err := cb.FileCreate(userOne.Username, projectID, fileOneName, ".")
	require.Nil(t, err, "MySQL threw an error when creating a new file")
	require.True(t, fileID >= 0, "MySQL returned an invalid FileID")

	files, err := cb.ProjectGetFiles(projectID)
	require.Nil(t, err, "MySQL threw an error when getting files for given project")
	require.Len(t, files, 1, "MySQL returned incorrect number of files for given project")
	assert.Equal(t, fileID, files[0].FileID, "MySQL returned file with incorrect FileID")

	err = cb.FileDelete(fileID)
	require.Nil(t, err, "MySQL threw an error when deleting the new file")

	files, err = cb.ProjectGetFiles(projectID)
	require.Nil(t, err, "MySQL threw an error when getting files for given project")
	require.Len(t, files, 0, "MySQL returned incorrect number of files for given project")

	err = cb.FileDelete(fileID)
	require.EqualError(t, err, "No such fileID found", "MySQL threw incorrect error for deletion of nonexistent fileID")
}
