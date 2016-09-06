package dbfs

import (
	"testing"
	"time"
)

func TestDatabaseImpl_OpenMySQLConn(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	my, err := di.getMySQLConn()
	defer di.CloseMySQL()

	if err != nil {
		t.Fatal(err)
	}

	err = my.db.Ping()

	if err != nil {
		t.Fatal(err)
	}

}

func TestDatabaseImpl_CloseMySQL(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	_, err := di.getMySQLConn()
	if err != nil {
		t.Fatal(err)
	}
	err = di.CloseMySQL()
	if err != nil {
		t.Fatal(err)
	}
	err = di.CloseMySQL()
	if err != ErrDbNotInitialized {
		t.Fatal("Wrong error recieved")
	}
}

func TestDatabaseImpl_MySQLUserRegister(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	di.MySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	err := di.MySQLUserRegister(user)
	if err != nil {
		t.Fatal(err)
	}
	err = di.MySQLUserDelete("jshap70", "secret")
	if err == ErrNoDbChange {
		t.Fatal("No user added")
	}
}

func TestDatabaseImpl_MySQLUserGetPass(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)
	di.MySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	err := di.MySQLUserRegister(user)
	if err != nil {
		t.Fatal(err)
	}

	pass, err := di.MySQLUserGetPass("jshap70")
	if err != nil {
		t.Fatal(err)
	}
	if pass != "secret" {
		t.Fatal("Wrong password returned")
	}

	err = di.MySQLUserDelete("jshap70", "secret")
}

func TestDatabaseImpl_MySQLUserLookup(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)
	di.MySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	err := di.MySQLUserRegister(user)
	if err != nil {
		t.Fatal(err)
	}

	userRet, err := di.MySQLUserLookup("jshap70")
	if err != nil {
		t.Fatal(err)
	}
	if userRet.FirstName != "Joel" || userRet.LastName != "Shapiro" || userRet.Email != "joel@codecollab.cc" {
		t.Fatalf("Wrong return, got: %v %v, email: %v", userRet.FirstName, userRet.LastName, userRet.Email)
	}

	userRet, err = di.MySQLUserLookup("notjshap70")
	if err == nil {
		t.Fatal("Expected lookup with incorrect username to fail, but it did not")
	}

	err = di.MySQLUserDelete("jshap70", "secret")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseImpl_MySQLUserProjects(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)
	di.MySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	erro := di.MySQLUserRegister(user)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")

	projects, err := di.MySQLUserProjects("jshap70")
	_ = di.MySQLProjectDelete(projectID, "jshap70")
	_ = di.MySQLUserDelete("jshap70", "secret")
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
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	erro := di.MySQLUserRegister(user)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, err := di.MySQLProjectCreate("jshap70", "codecollabcore")
	if err != nil {
		t.Fatal(err)
	}
	if projectID < 0 {
		t.Fatal("incorrect ProjectID")
	}

	_, err = di.MySQLProjectCreate("jshap70", "codecollabcore")
	if err == nil {
		t.Fatal("unexpected opperation allowed")
	}

	err = di.MySQLProjectDelete(projectID, "jshap70")
	if err != nil {
		t.Fatal(err, projectID)
	}
	err = di.MySQLUserDelete("jshap70", "secret")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseImpl_MySQLProjectDelete(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	erro := di.MySQLUserRegister(user)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, err := di.MySQLProjectCreate("jshap70", "codecollabcore")
	if err != nil {
		t.Fatal(err)
	}

	// test trying to delete a project that contains files
	_, err = di.MySQLFileCreate("jshap70", "file-y", ".", projectID)
	if err != nil {
		t.Fatal(err)
	}

	err = di.MySQLProjectDelete(projectID, "jshap70")
	if err != nil {
		t.Fatal(err)
	}
	err = di.MySQLProjectDelete(projectID, "jshap70")
	if err == nil {
		t.Fatal("project delete succeded 2x on the same projectID")
	}

	err = di.MySQLUserDelete("jshap70", "secret")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseImpl_MySQLProjectGetFiles(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	erro := di.MySQLUserRegister(user)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, err := di.MySQLProjectCreate("jshap70", "codecollabcore")
	di.MySQLFileCreate("jshap70", "file-y", ".", projectID)

	files, err := di.MySQLProjectGetFiles(projectID)

	_ = di.MySQLProjectDelete(projectID, "jshap70")
	_ = di.MySQLUserDelete("jshap70", "secret")

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatalf("Project %v returned not the correct length, expected: 1, actual: %v", projectID, len(files))
	}
	if files[0].FileID == -1 || files[0].Creator != "jshap70" || files[0].RelativePath != "." || files[0].Filename != "file-y" || files[0].ProjectID != projectID {
		t.Fatalf("Wrong return, got project: %v", files[0])
	}

	files, err = di.MySQLProjectGetFiles(projectID + 1000)
	if err == nil {
		t.Fatal("Expected lookup to fail when using an incorrect projectID")
	}
}

func TestDatabaseImpl_MySQLProjectGrantPermission(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	userJoel := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	userAustin := UserMeta{
		Username:  "fahslaj",
		Password:  "secret",
		Email:     "austin@codecollab.cc",
		FirstName: "Austin",
		LastName:  "Fahsl"}

	erro := di.MySQLUserRegister(userJoel)
	if erro != nil {
		t.Fatal(erro)
	}

	di.MySQLUserRegister(userAustin)

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")

	err := di.MySQLProjectGrantPermission(projectID, "fahslaj", 5, "jshap70")
	if err != nil {
		t.Fatal(err)
	}

	projects, err := di.MySQLUserProjects("fahslaj")
	if err != nil {
		t.Fatal(err)
	}

	if len(projects) != 1 {
		t.Fatalf("Projects returned not the correct length, expected: 1, actual: %v", len(projects))
	}
	if projects[0].ProjectID != projectID || projects[0].Name != "codecollabcore" || projects[0].PermissionLevel != 5 {
		t.Fatalf("Wrong return, got project:%v %v, perm: %v", projects[0].Name, projects[0].ProjectID, projects[0].PermissionLevel)
	}

	err = di.MySQLProjectDelete(projectID, "jshap70")
	if err != nil {
		t.Fatal(err)
	}
	err = di.MySQLUserDelete("fahslaj", "secret")
	if err != nil {
		t.Fatal(err)
	}
	err = di.MySQLUserDelete("jshap70", "secret")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseImpl_MySQLProjectLookup(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	userJoel := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	userAustin := UserMeta{
		Username:  "fahslaj",
		Password:  "secret",
		Email:     "austin@codecollab.cc",
		FirstName: "Austin",
		LastName:  "Fahsl"}

	erro := di.MySQLUserRegister(userJoel)
	if erro != nil {
		t.Fatal(erro)
	}

	di.MySQLUserRegister(userAustin)

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")

	defer di.MySQLUserDelete("fahslaj", "secret")
	defer di.MySQLUserDelete("jshap70", "secret")
	defer di.MySQLProjectDelete(projectID, "jshap70")

	name, perms, err := di.MySQLProjectLookup(projectID, "fahslaj")
	if err == nil {
		t.Fatal("Expected failure when given a projectID you don't have access to")
	}

	err = di.MySQLProjectGrantPermission(projectID, "fahslaj", 5, "jshap70")
	if err != nil {
		t.Fatal(err)
	}

	name, perms, err = di.MySQLProjectLookup(projectID, "fahslaj")

	if err != nil {
		t.Fatal(err)
	}
	if name != "codecollabcore" {
		t.Fatalf("Incorrect name: %v", name)
	}
	if len(perms) != 2 {
		t.Fatalf("Projects returned not the correct length, expected: 1, actual: %v", len(perms))
	}

	if perms["jshap70"].PermissionLevel != 10 {
		t.Fatalf("jshap70 had permision level: %v", perms["jshap70"].PermissionLevel)
	}
	if perms["fahslaj"].PermissionLevel != 5 {
		t.Fatalf("fahslaj had permision level: %v", perms["fahslaj"].PermissionLevel)
	}
	if perms["fahslaj"].GrantedDate == time.Unix(0, 0) {
		t.Fatal("time did not correctly parse")
	}

	name, perms, err = di.MySQLProjectLookup(projectID+1000, "fahslaj")
	if err == nil {
		t.Fatal("Expected failure when given a non-existant projectID")
	}
}

func TestDatabaseImpl_MySQLProjectRevokePermission(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	userJoel := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	userAustin := UserMeta{
		Username:  "fahslaj",
		Password:  "secret",
		Email:     "austin@codecollab.cc",
		FirstName: "Austin",
		LastName:  "Fahsl"}

	erro := di.MySQLUserRegister(userJoel)
	if erro != nil {
		t.Fatal(erro)
	}

	di.MySQLUserRegister(userAustin)

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")

	di.MySQLProjectGrantPermission(projectID, "fahslaj", 5, "jshap70")

	projects, _ := di.MySQLUserProjects("fahslaj")
	if len(projects) != 1 {
		t.Fatalf("Projects returned not the correct length, expected: 1, actual: %v", len(projects))
	}
	if projects[0].ProjectID != projectID || projects[0].PermissionLevel != 5 {
		t.Fatalf("Wrong return, got project:%v %v, perm: %v", projects[0].Name, projects[0].ProjectID, projects[0].PermissionLevel)
	}

	di.MySQLProjectRevokePermission(projectID, "fahslaj", "jshap70")
	_ = di.MySQLProjectDelete(projectID, "jshap70")
	_ = di.MySQLUserDelete("jshap70", "secret")
	_ = di.MySQLUserDelete("fahslaj", "secret")

	projects, _ = di.MySQLUserProjects("fahslaj")
	if len(projects) > 0 {
		t.Fatalf("Projects returned not the correct length, expected: 0, actual: %v", len(projects))
	}
}

func TestDatabaseImpl_MySQLProjectRename(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	di.MySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	erro := di.MySQLUserRegister(user)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")

	err := di.MySQLProjectRename(projectID, "newName")
	if err != nil {
		t.Fatal(err)
	}

	projects, err := di.MySQLUserProjects("jshap70")
	_ = di.MySQLProjectDelete(projectID, "jshap70")
	_ = di.MySQLUserDelete("jshap70", "secret")

	if projects[0].ProjectID != projectID || projects[0].Name != "newName" {
		t.Fatalf("Wrong return, got project:%v %v", projects[0].Name, projects[0].ProjectID)
	}
}

func TestDatabaseImpl_MySQLFileCreate(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	erro := di.MySQLUserRegister(user)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")
	fileID, err := di.MySQLFileCreate("jshap70", "file-y", ".", projectID)

	files, _ := di.MySQLProjectGetFiles(projectID)

	_ = di.MySQLProjectDelete(projectID, "jshap70")
	_ = di.MySQLUserDelete("jshap70", "secret")

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatalf("Project %v returned not the correct length, expected: 1, actual: %v", projectID, len(files))
	}
	if files[0].FileID != fileID || files[0].Creator != "jshap70" || files[0].RelativePath != "." || files[0].Filename != "file-y" || files[0].ProjectID != projectID {
		t.Fatalf("Wrong return, got project: %v", files[0])
	}
}

func TestDatabaseImpl_MySQLFileDelete(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	erro := di.MySQLUserRegister(user)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")
	fileID, _ := di.MySQLFileCreate("jshap70", "file-y", ".", projectID)
	err := di.MySQLFileDelete(fileID)

	files, _ := di.MySQLProjectGetFiles(projectID)
	_ = di.MySQLProjectDelete(projectID, "jshap70")
	_ = di.MySQLUserDelete("jshap70", "secret")

	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Fatalf("Project %v returned not the correct length, expected: 0, actual: %v", projectID, len(files))
	}
}

func TestDatabaseImpl_MySQLFileMove(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	erro := di.MySQLUserRegister(user)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")
	fileID, _ := di.MySQLFileCreate("jshap70", "file-y", ".", projectID)

	err := di.MySQLFileMove(fileID, "cc")

	files, _ := di.MySQLProjectGetFiles(projectID)
	_ = di.MySQLProjectDelete(projectID, "jshap70")
	_ = di.MySQLUserDelete("jshap70", "secret")

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
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	erro := di.MySQLUserRegister(user)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")
	fileID, _ := di.MySQLFileCreate("jshap70", "file-y", ".", projectID)

	err := di.MySQLFileRename(fileID, "file-z")

	files, _ := di.MySQLProjectGetFiles(projectID)
	_ = di.MySQLProjectDelete(projectID, "jshap70")
	_ = di.MySQLUserDelete("jshap70", "secret")

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
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	erro := di.MySQLUserRegister(user)
	if erro != nil {
		t.Fatal(erro)
	}

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")
	fileID, _ := di.MySQLFileCreate("jshap70", "file-y", ".", projectID)

	filebefore, err := di.MySQLFileGetInfo(fileID)
	_ = di.MySQLFileMove(fileID, "cc")
	fileafter, err := di.MySQLFileGetInfo(fileID)

	_ = di.MySQLProjectDelete(projectID, "jshap70")
	_ = di.MySQLUserDelete("jshap70", "secret")

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
