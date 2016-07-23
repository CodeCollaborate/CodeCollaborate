package dbfs

import (
	"testing"
	"time"
)

func TestOpenMySQLConn(t *testing.T) {
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

func TestCloseMySQL(t *testing.T) {
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

func TestMySQLUserRegister(t *testing.T) {
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

func TestMySQLUserGetPass(t *testing.T) {
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

func TestMySQLUserLookup(t *testing.T) {
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

	err = di.MySQLUserDelete("jshap70", "secret")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMySQLUserProjects(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)
	di.MySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	di.MySQLUserRegister(user)
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
	if projects[0].ProjectID == -1 || projects[0].ProjectName != "codecollabcore" || projects[0].PermissionLevel != 10 {
		t.Fatalf("Wrong return, got project:%v %v, perm: %v", projects[0].ProjectName, projects[0].ProjectID, projects[0].PermissionLevel)
	}
}

func TestMySQLProjectCreate(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	di.MySQLUserRegister(user)

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

func TestMySQLProjectDelete(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	di.MySQLUserRegister(user)

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

func TestMySQLProjectGetFiles(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	di.MySQLUserRegister(user)

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
}

func TestMySQLProjectGrantPermission(t *testing.T) {
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

	di.MySQLUserRegister(userJoel)
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
	if projects[0].ProjectID != projectID || projects[0].ProjectName != "codecollabcore" || projects[0].PermissionLevel != 5 {
		t.Fatalf("Wrong return, got project:%v %v, perm: %v", projects[0].ProjectName, projects[0].ProjectID, projects[0].PermissionLevel)
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

func TestMySQLProjectLookup(t *testing.T) {
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

	di.MySQLUserRegister(userJoel)
	di.MySQLUserRegister(userAustin)

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")

	err := di.MySQLProjectGrantPermission(projectID, "fahslaj", 5, "jshap70")
	if err != nil {
		t.Fatal(err)
	}

	name, perms, err := di.MySQLProjectLookup(projectID, "fahslaj")
	di.MySQLProjectDelete(projectID, "jshap70")
	di.MySQLUserDelete("fahslaj", "secret")
	di.MySQLUserDelete("jshap70", "secret")

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
}

func TestMySQLProjectRevokePermission(t *testing.T) {
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

	di.MySQLUserRegister(userJoel)
	di.MySQLUserRegister(userAustin)

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")

	di.MySQLProjectGrantPermission(projectID, "fahslaj", 5, "jshap70")

	projects, _ := di.MySQLUserProjects("fahslaj")
	if len(projects) != 1 {
		t.Fatalf("Projects returned not the correct length, expected: 1, actual: %v", len(projects))
	}
	if projects[0].ProjectID != projectID || projects[0].PermissionLevel != 5 {
		t.Fatalf("Wrong return, got project:%v %v, perm: %v", projects[0].ProjectName, projects[0].ProjectID, projects[0].PermissionLevel)
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

func TestMySQLProjectRename(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	di.MySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	di.MySQLUserRegister(user)

	projectID, _ := di.MySQLProjectCreate("jshap70", "codecollabcore")

	err := di.MySQLProjectRename(projectID, "newName")
	if err != nil {
		t.Fatal(err)
	}

	projects, err := di.MySQLUserProjects("jshap70")
	_ = di.MySQLProjectDelete(projectID, "jshap70")
	_ = di.MySQLUserDelete("jshap70", "secret")

	if projects[0].ProjectID != projectID || projects[0].ProjectName != "newName" {
		t.Fatalf("Wrong return, got project:%v %v", projects[0].ProjectName, projects[0].ProjectID)
	}
}

func TestMySQLFileCreate(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	di.MySQLUserRegister(user)

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

func TestMySQLFileDelete(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	di.MySQLUserRegister(user)

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

func TestMySQLFileMove(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	di.MySQLUserRegister(user)

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

func TestMySQLRenameFile(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	di.MySQLUserRegister(user)

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

func TestMySQLFileGetInfo(t *testing.T) {
	configSetup()
	di := new(DatabaseImpl)

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	di.MySQLUserRegister(user)

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
