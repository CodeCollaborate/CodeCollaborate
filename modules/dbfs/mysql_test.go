package dbfs

import (
	"testing"
	"time"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/stretchr/testify/assert"
)

var userOne = UserMeta{
	Username:  "_test_user1",
	Password:  "secret",
	Email:     "_test_email1@codecollab.cc",
	FirstName: "Joel",
	LastName:  "Shapiro"}

var userTwo = UserMeta{
	Username:  "_test_user2",
	Password:  "secret",
	Email:     "_test_email2@codecollab.cc",
	FirstName: "Austin",
	LastName:  "Fahsl"}

func TestDatabaseImpl_MySQLUserRegister(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	di.MySQLUserDelete(userOne.Username)

	err := di.MySQLUserRegister(userOne)
	if err != nil {
		t.Fatal(err)
	}
	_, err = di.MySQLUserDelete(userOne.Username)
	if err == ErrNoDbChange {
		t.Fatal("No user added")
	}
}

func TestDatabaseImpl_MySQLUserGetPass(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)
	di.MySQLUserDelete(userOne.Username)

	err := di.MySQLUserRegister(userOne)
	if err != nil {
		t.Fatal(err)
	}

	pass, err := di.MySQLUserGetPass(userOne.Username)
	if err != nil {
		t.Fatal(err)
	}
	if pass != userOne.Password {
		t.Fatal("Wrong password returned")
	}

	di.MySQLUserDelete(userOne.Username)
}

func TestDatabaseImpl_MySQLUserDelete(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)
	di.MySQLUserDelete(userOne.Username)
	di.MySQLUserDelete(userTwo.Username)

	//db.MySQLUserRegister(geneMeta)
	err := di.MySQLUserRegister(userOne)
	assert.NoError(t, err)

	//closures, err := req.process(db)
	projectIDs, err := di.MySQLUserDelete(userOne.Username)
	assert.NoError(t, err)

	assert.Empty(t, projectIDs, "expected 0 projects to be deleted")

	// check user actually deleted
	returnedUser, err := di.MySQLUserLookup(userOne.Username)
	assert.EqualError(t, err, "No such username found", "expected no user to be returned")
	assert.Equal(t, UserMeta{}, returnedUser, "expected no user to be returned, also no error was thrown on empty data")

	// test with projects for notifications
	err = di.MySQLUserRegister(userOne)
	assert.NoError(t, err)
	err = di.MySQLUserRegister(userTwo)
	assert.NoError(t, err)

	projectID1, err := di.MySQLProjectCreate(userOne.Username, "_test_project_1")
	projectID2, err := di.MySQLProjectCreate(userOne.Username, "_test_project_2")

	writePerm, err := config.PermissionByLabel("write")
	assert.NoError(t, err, "api permissions error")

	err = di.MySQLProjectGrantPermission(projectID1, userTwo.Username, writePerm.Level, userOne.Username)
	assert.NoError(t, err, "project grant permission error")
	err = di.MySQLProjectGrantPermission(projectID2, userTwo.Username, writePerm.Level, userOne.Username)
	assert.NoError(t, err, "project grant permission error")

	projectIDs, err = di.MySQLUserDelete(userOne.Username)
	assert.NoError(t, err)

	// check that it claims to have deleted both owned projects
	assert.Len(t, projectIDs, 2, "expected 2 projects to be deleted")
	assert.Contains(t, projectIDs, projectID1, "didn't delete _test_project_1")
	assert.Contains(t, projectIDs, projectID2, "didn't delete _test_project_2")

	// check user actually deleted
	returnedUser, err = di.MySQLUserLookup(userOne.Username)
	assert.EqualError(t, err, "No such username found", "expected no user to be returned")
	assert.Equal(t, UserMeta{}, returnedUser, "expected no user to be returned, also no error was thrown on empty data")

	// check projects actually deleted
	_, _, err = di.MySQLProjectLookup(projectID1, userTwo.Username)
	assert.EqualError(t, err, "No such projectID found", "expected project1 to not exist")

	_, _, err = di.MySQLProjectLookup(projectID2, userTwo.Username)
	assert.EqualError(t, err, "No such projectID found", "expected project2 to not exist")
}

func TestDatabaseImpl_MySQLUserLookup(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)
	di.MySQLUserDelete(userOne.Username)

	err := di.MySQLUserRegister(userOne)
	if err != nil {
		t.Fatal(err)
	}

	userRet, err := di.MySQLUserLookup(userOne.Username)
	if err != nil {
		t.Fatal(err)
	}
	if userRet.FirstName != userOne.FirstName || userRet.LastName != userOne.LastName || userRet.Email != userOne.Email {
		t.Fatalf("Wrong return, got: %v %v, email: %v", userRet.FirstName, userRet.LastName, userRet.Email)
	}

	userRet, err = di.MySQLUserLookup("notjshap70")
	if err == nil {
		t.Fatal("Expected lookup with incorrect username to fail, but it did not")
	}

	_, err = di.MySQLUserDelete(userOne.Username)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseImpl_MySQLUserProjects(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)
	di.MySQLUserDelete(userOne.Username)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate(userOne.Username, "codecollabcore")

	projects, err := di.MySQLUserProjects(userOne.Username)
	_ = di.MySQLProjectDelete(projectID, userOne.Username)
	di.MySQLUserDelete(userOne.Username)
	if err != nil {
		t.Fatal(err)
	}

	if len(projects) != 1 {
		t.Fatalf("Projects returned not the correct length, expected: 1, actual: %v", len(projects))
	}
	if projects[0].ProjectID == -1 || projects[0].Name != "codecollabcore" || projects[0].PermissionLevel != 10 {
		t.Fatalf("Wrong return, got project:%v %v, perm: %v", projects[0].Name, projects[0].ProjectID, projects[0].PermissionLevel)
	}
}

func TestDatabaseImpl_MySQLProjectCreate(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, err := di.MySQLProjectCreate(userOne.Username, "codecollabcore")
	if err != nil {
		t.Fatal(err)
	}
	if projectID < 0 {
		t.Fatal("incorrect ProjectID")
	}

	_, err = di.MySQLProjectCreate(userOne.Username, "codecollabcore")
	if err == nil {
		t.Fatal("unexpected opperation allowed")
	}

	err = di.MySQLProjectDelete(projectID, userOne.Username)
	if err != nil {
		t.Fatal(err, projectID)
	}
	_, err = di.MySQLUserDelete(userOne.Username)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseImpl_MySQLProjectDelete(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, err := di.MySQLProjectCreate(userOne.Username, "codecollabcore")
	if err != nil {
		t.Fatal(err)
	}

	// test trying to delete a project that contains files
	_, err = di.MySQLFileCreate(userOne.Username, "file-y", ".", projectID)
	if err != nil {
		t.Fatal(err)
	}

	err = di.MySQLProjectDelete(projectID, userOne.Username)
	if err != nil {
		t.Fatal(err)
	}
	err = di.MySQLProjectDelete(projectID, userOne.Username)
	if err == nil {
		t.Fatal("project delete succeded 2x on the same projectID")
	}

	_, err = di.MySQLUserDelete(userOne.Username)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseImpl_MySQLProjectGetFiles(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, err := di.MySQLProjectCreate(userOne.Username, "codecollabcore")
	di.MySQLFileCreate(userOne.Username, "file-y", ".", projectID)

	files, err := di.MySQLProjectGetFiles(projectID)

	_ = di.MySQLProjectDelete(projectID, userOne.Username)
	_, _ = di.MySQLUserDelete(userOne.Username)

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatalf("Project %v returned not the correct length, expected: 1, actual: %v", projectID, len(files))
	}
	if files[0].FileID == -1 || files[0].Creator != userOne.Username || files[0].RelativePath != "." || files[0].Filename != "file-y" || files[0].ProjectID != projectID {
		t.Fatalf("Wrong return, got project: %v", files[0])
	}

	//files, err = di.MySQLProjectGetFiles(projectID + 1000)
	//if err == nil {
	//	t.Fatal("Expected lookup to fail when using an incorrect projectID")
	//}
}

func TestDatabaseImpl_MySQLProjectGrantPermission(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	err := di.MySQLUserRegister(userOne)
	if err != nil {
		di.MySQLUserDelete(userOne.Username)
		di.MySQLUserDelete(userTwo.Username)
		err = di.MySQLUserRegister(userOne)
		assert.NoError(t, err)
	}

	di.MySQLUserRegister(userTwo)

	projectID, _ := di.MySQLProjectCreate(userOne.Username, "codecollabcore")

	err = di.MySQLProjectGrantPermission(projectID, userTwo.Username, 5, userOne.Username)
	if err != nil {
		t.Fatal(err)
	}

	projects, err := di.MySQLUserProjects(userTwo.Username)
	if err != nil {
		t.Fatal(err)
	}

	if len(projects) != 1 {
		t.Fatalf("Projects returned not the correct length, expected: 1, actual: %v", len(projects))
	}
	if projects[0].ProjectID != projectID || projects[0].Name != "codecollabcore" || projects[0].PermissionLevel != 5 {
		t.Fatalf("Wrong return, got project:%v %v, perm: %v", projects[0].Name, projects[0].ProjectID, projects[0].PermissionLevel)
	}

	err = di.MySQLProjectDelete(projectID, userOne.Username)
	if err != nil {
		t.Fatal(err)
	}
	_, err = di.MySQLUserDelete(userTwo.Username)
	if err != nil {
		t.Fatal(err)
	}
	_, err = di.MySQLUserDelete(userOne.Username)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseImpl_MySQLProjectLookup(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}

	di.MySQLUserRegister(userTwo)

	projectID, _ := di.MySQLProjectCreate(userOne.Username, "codecollabcore")

	defer di.MySQLUserDelete(userTwo.Username)
	defer di.MySQLUserDelete(userOne.Username)
	defer di.MySQLProjectDelete(projectID, userOne.Username)

	err := di.MySQLProjectGrantPermission(projectID, userTwo.Username, 5, userOne.Username)
	if err != nil {
		t.Fatal(err)
	}

	name, perms, err := di.MySQLProjectLookup(projectID, userTwo.Username)

	if err != nil {
		t.Fatal(err)
	}
	if name != "codecollabcore" {
		t.Fatalf("Incorrect name: %v", name)
	}
	if len(perms) != 2 {
		t.Fatalf("Projects returned not the correct length, expected: 1, actual: %v", len(perms))
	}

	if perms[userOne.Username].PermissionLevel != 10 {
		t.Fatalf("jshap70 had permision level: %v", perms[userOne.Username].PermissionLevel)
	}
	if perms[userTwo.Username].PermissionLevel != 5 {
		t.Fatalf("fahslaj had permision level: %v", perms[userTwo.Username].PermissionLevel)
	}
	if perms[userTwo.Username].GrantedDate == time.Unix(0, 0) {
		t.Fatal("time did not correctly parse")
	}

	name, perms, err = di.MySQLProjectLookup(projectID+1000, userTwo.Username)
	if err == nil {
		t.Fatal("Expected failure when given a non-existant projectID")
	}
}

func TestDatabaseImpl_MySQLProjectRevokePermission(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}

	di.MySQLUserRegister(userTwo)

	projectID, _ := di.MySQLProjectCreate(userOne.Username, "codecollabcore")

	di.MySQLProjectGrantPermission(projectID, userTwo.Username, 5, userOne.Username)

	projects, _ := di.MySQLUserProjects(userTwo.Username)
	if len(projects) != 1 {
		t.Fatalf("Projects returned not the correct length, expected: 1, actual: %v", len(projects))
	}
	if projects[0].ProjectID != projectID || projects[0].PermissionLevel != 5 {
		t.Fatalf("Wrong return, got project:%v %v, perm: %v", projects[0].Name, projects[0].ProjectID, projects[0].PermissionLevel)
	}

	di.MySQLProjectRevokePermission(projectID, userTwo.Username, userOne.Username)
	_ = di.MySQLProjectDelete(projectID, userOne.Username)
	di.MySQLUserDelete(userOne.Username)
	di.MySQLUserDelete(userTwo.Username)

	projects, _ = di.MySQLUserProjects(userTwo.Username)
	if len(projects) > 0 {
		t.Fatalf("Projects returned not the correct length, expected: 0, actual: %v", len(projects))
	}
}

func TestDatabaseImpl_MySqlUserProjectPermissionLookup(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	di.MySQLUserDelete(userOne.Username)
	di.MySQLUserDelete(userTwo.Username)

	defer func() {
		di.MySQLUserDelete(userOne.Username)
		di.MySQLUserDelete(userTwo.Username)
	}()

	err := di.MySQLUserRegister(userOne)
	assert.Nil(t, err)

	projectID, _ := di.MySQLProjectCreate(userOne.Username, "codecollabcore")
	defer di.MySQLProjectDelete(projectID, userOne.Username)

	permLevel, err := di.MySQLUserProjectPermissionLookup(projectID, userOne.Username)
	assert.Nil(t, err, "unexpected error from mysql permission lookup")
	ownerPerm, _ := config.PermissionByLabel("owner")
	assert.Equal(t, ownerPerm.Level, permLevel, "expected user to be owner")

	err = di.MySQLUserRegister(userTwo)
	assert.Nil(t, err)

	permLevel, err = di.MySQLUserProjectPermissionLookup(projectID, userTwo.Username)
	assert.Nil(t, err, "expected error from mysql permission lookup")
	assert.Equal(t, 0, permLevel, "expected user not have permission")

	readPerm, _ := config.PermissionByLabel("read")
	err = di.MySQLProjectGrantPermission(projectID, userTwo.Username, readPerm.Level, userOne.Username)
	assert.Nil(t, err)

	permLevel, err = di.MySQLUserProjectPermissionLookup(projectID, userTwo.Username)
	assert.Nil(t, err, "unexpected error from mysql permission lookup")
	assert.Equal(t, readPerm.Level, permLevel, "expected user have read permission")
}

func TestDatabaseImpl_MySQLProjectRename(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	di.MySQLUserDelete(userOne.Username)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate(userOne.Username, "codecollabcore")

	err := di.MySQLProjectRename(projectID, "newName")
	if err != nil {
		t.Fatal(err)
	}

	projects, err := di.MySQLUserProjects(userOne.Username)
	di.MySQLProjectDelete(projectID, userOne.Username)
	di.MySQLUserDelete(userOne.Username)

	if projects[0].ProjectID != projectID || projects[0].Name != "newName" {
		t.Fatalf("Wrong return, got project:%v %v", projects[0].Name, projects[0].ProjectID)
	}
}

func TestDatabaseImpl_MySQLFileCreate(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}
	filename := "file-y"

	projectID, _ := di.MySQLProjectCreate(userOne.Username, "codecollabcore")
	fileID, err := di.MySQLFileCreate(userOne.Username, filename, ".", projectID)

	files, _ := di.MySQLProjectGetFiles(projectID)

	defer di.MySQLUserDelete(userOne.Username)
	defer di.MySQLProjectDelete(projectID, userOne.Username)

	assert.NoError(t, err, "mysql error")
	assert.Equal(t, 1, len(files), "Project incorrect file count")
	assert.Equal(t, fileID, files[0].FileID, "incorrect fileID")
	assert.Equal(t, userOne.Username, files[0].Creator, "incorrect creator")
	assert.Equal(t, ".", files[0].RelativePath, "incorrect relative path")
	assert.Equal(t, filename, files[0].Filename, "incorrect filename")
	assert.Equal(t, projectID, files[0].ProjectID, "incorrect projectID")

	// should fail b/c location is already in use
	fileIDNew, err := di.MySQLFileCreate(userOne.Username, filename, ".", projectID)
	assert.EqualValues(t, -1, fileIDNew, "Expected invalid FileID to be returned")
	assert.Error(t, err, "expected duplicate insertion to fail")
}

func TestDatabaseImpl_MySQLFileDelete(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate(userOne.Username, "codecollabcore")
	fileID, _ := di.MySQLFileCreate(userOne.Username, "file-y", ".", projectID)
	err := di.MySQLFileDelete(fileID)

	files, _ := di.MySQLProjectGetFiles(projectID)
	_ = di.MySQLProjectDelete(projectID, userOne.Username)
	di.MySQLUserDelete(userOne.Username)

	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Fatalf("Project %v returned not the correct length, expected: 0, actual: %v", projectID, len(files))
	}
}

func TestDatabaseImpl_MySQLFileMove(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate(userOne.Username, "codecollabcore")
	fileID, _ := di.MySQLFileCreate(userOne.Username, "file-y", ".", projectID)

	err := di.MySQLFileMove(fileID, "cc")

	files, _ := di.MySQLProjectGetFiles(projectID)
	_ = di.MySQLProjectDelete(projectID, userOne.Username)
	di.MySQLUserDelete(userOne.Username)

	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("Project %v returned not the correct length, expected: 1, actual: %v", projectID, len(files))
	}
	if files[0].FileID != fileID || files[0].RelativePath != "cc" || files[0].ProjectID != projectID {
		t.Fatalf("Wrong return, got project: %v", files[0])
	}
}

func TestDatabaseImpl_MySQLRenameFile(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate(userOne.Username, "codecollabcore")
	fileID, _ := di.MySQLFileCreate(userOne.Username, "file-y", ".", projectID)

	err := di.MySQLFileRename(fileID, "file-z")

	files, _ := di.MySQLProjectGetFiles(projectID)
	_ = di.MySQLProjectDelete(projectID, userOne.Username)
	di.MySQLUserDelete(userOne.Username)

	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("Project %v returned not the correct length, expected: 1, actual: %v", projectID, len(files))
	}
	if files[0].FileID != fileID || files[0].Filename != "file-z" || files[0].ProjectID != projectID {
		t.Fatalf("Wrong return, got project: %v", files[0])
	}
}

func TestDatabaseImpl_MySQLFileGetInfo(t *testing.T) {
	testConfigSetup(t)
	di := new(DatabaseImpl)

	erro := di.MySQLUserRegister(userOne)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate(userOne.Username, "codecollabcore")
	fileID, _ := di.MySQLFileCreate(userOne.Username, "file-y", ".", projectID)

	filebefore, err := di.MySQLFileGetInfo(fileID)
	_ = di.MySQLFileMove(fileID, "cc")
	fileafter, err := di.MySQLFileGetInfo(fileID)

	_ = di.MySQLProjectDelete(projectID, userOne.Username)
	di.MySQLUserDelete(userOne.Username)

	if err != nil {
		t.Fatal(err)
	}
	if filebefore.FileID != fileID || filebefore.RelativePath != "." || filebefore.ProjectID != projectID {
		t.Fatalf("Wrong return, got project: %v", filebefore)
	}
	if fileafter.FileID != fileID || fileafter.RelativePath != "cc" || fileafter.ProjectID != projectID {
		t.Fatalf("Wrong return, got project: %v", filebefore)
	}
}
