package dbfs

import (
	"testing"
	"time"
)

func TestOpenMySQLConn(t *testing.T) {
	configSetup()

	my, err := getMySQLConn()
	defer CloseMySQL()

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
	_, err := getMySQLConn()
	if err != nil {
		t.Fatal(err)
	}
	err = CloseMySQL()
	if err != nil {
		t.Fatal(err)
	}
	err = CloseMySQL()
	if err != ErrDbNotInitialized {
		t.Fatal("Wrong error recieved")
	}
}

func TestMySQLUserRegister(t *testing.T) {
	configSetup()
	mySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	err := MySQLUserRegister(user)
	if err != nil {
		t.Fatal(err)
	}
	err = mySQLUserDelete("jshap70", "secret")
	if err == ErrNoDbChange {
		t.Fatal("No user added")
	}
}

func TestMySQLUserGetPass(t *testing.T) {
	configSetup()
	mySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	err := MySQLUserRegister(user)
	if err != nil {
		t.Fatal(err)
	}

	pass, err := MySQLUserGetPass("jshap70")
	if err != nil {
		t.Fatal(err)
	}
	if pass != "secret" {
		t.Fatal("Wrong password returned")
	}

	err = mySQLUserDelete("jshap70", "secret")
}

func TestMySQLUserLookup(t *testing.T) {
	configSetup()
	mySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	err := MySQLUserRegister(user)
	if err != nil {
		t.Fatal(err)
	}

	userRet, err := MySQLUserLookup("jshap70")
	if err != nil {
		t.Fatal(err)
	}
	if userRet.FirstName != "Joel" || userRet.LastName != "Shapiro" || userRet.Email != "joel@codecollab.cc" {
		t.Fatalf("Wrong return, got: %v %v, email: %v", userRet.FirstName, userRet.LastName, userRet.Email)
	}

	err = mySQLUserDelete("jshap70", "secret")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMySQLUserProjects(t *testing.T) {
	configSetup()
	mySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	MySQLUserRegister(user)
	projectID, _ := MySQLProjectCreate("jshap70", "codecollabcore")

	projects, err := MySQLUserProjects("jshap70")
	_ = MySQLProjectDelete(projectID)
	_ = mySQLUserDelete("jshap70", "secret")
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

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	MySQLUserRegister(user)

	projectID, err := MySQLProjectCreate("jshap70", "codecollabcore")
	if err != nil {
		t.Fatal(err)
	}
	if projectID < 0 {
		t.Fatal("incorrect ProjectID")
	}

	_, err = MySQLProjectCreate("jshap70", "codecollabcore")
	if err == nil {
		t.Fatal("unexpected opperation allowed")
	}

	err = MySQLProjectDelete(projectID)
	if err != nil {
		t.Fatal(err, projectID)
	}
	err = mySQLUserDelete("jshap70", "secret")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMySQLProjectDelete(t *testing.T) {
	configSetup()

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	MySQLUserRegister(user)

	projectID, err := MySQLProjectCreate("jshap70", "codecollabcore")
	if err != nil {
		t.Fatal(err)
	}

	// test trying to delete a project that contains files
	_, err = MySQLFileCreate("jshap70", "file-y", ".", projectID)
	if err != nil {
		t.Fatal(err)
	}

	err = MySQLProjectDelete(projectID)
	if err != nil {
		t.Fatal(err)
	}
	err = MySQLProjectDelete(projectID)
	if err == nil {
		t.Fatal("project delete succeded 2x on the same projectID")
	}

	err = mySQLUserDelete("jshap70", "secret")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMySQLProjectGetFiles(t *testing.T) {
	configSetup()

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	MySQLUserRegister(user)

	projectID, err := MySQLProjectCreate("jshap70", "codecollabcore")
	MySQLFileCreate("jshap70", "file-y", ".", projectID)

	files, err := MySQLProjectGetFiles(projectID)

	_ = MySQLProjectDelete(projectID)
	_ = mySQLUserDelete("jshap70", "secret")

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

	MySQLUserRegister(userJoel)
	MySQLUserRegister(userAustin)

	projectID, _ := MySQLProjectCreate("jshap70", "codecollabcore")

	err := MySQLProjectGrantPermission(projectID, "fahslaj", 5, "jshap70")
	if err != nil {
		t.Fatal(err)
	}

	projects, err := MySQLUserProjects("fahslaj")
	if err != nil {
		t.Fatal(err)
	}

	if len(projects) != 1 {
		t.Fatalf("Projects returned not the correct length, expected: 1, actual: %v", len(projects))
	}
	if projects[0].ProjectID != projectID || projects[0].ProjectName != "codecollabcore" || projects[0].PermissionLevel != 5 {
		t.Fatalf("Wrong return, got project:%v %v, perm: %v", projects[0].ProjectName, projects[0].ProjectID, projects[0].PermissionLevel)
	}

	err = MySQLProjectDelete(projectID)
	if err != nil {
		t.Fatal(err)
	}
	err = mySQLUserDelete("fahslaj", "secret")
	if err != nil {
		t.Fatal(err)
	}
	err = mySQLUserDelete("jshap70", "secret")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMySQLProjectLookup(t *testing.T) {
	configSetup()

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

	MySQLUserRegister(userJoel)
	MySQLUserRegister(userAustin)

	projectID, _ := MySQLProjectCreate("jshap70", "codecollabcore")

	err := MySQLProjectGrantPermission(projectID, "fahslaj", 5, "jshap70")
	if err != nil {
		t.Fatal(err)
	}

	name, perms, err := MySQLProjectLookup(projectID)
	MySQLProjectDelete(projectID)
	mySQLUserDelete("fahslaj", "secret")
	mySQLUserDelete("jshap70", "secret")

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

	MySQLUserRegister(userJoel)
	MySQLUserRegister(userAustin)

	projectID, _ := MySQLProjectCreate("jshap70", "codecollabcore")

	MySQLProjectGrantPermission(projectID, "fahslaj", 5, "jshap70")

	projects, _ := MySQLUserProjects("fahslaj")
	if len(projects) != 1 {
		t.Fatalf("Projects returned not the correct length, expected: 1, actual: %v", len(projects))
	}
	if projects[0].ProjectID != projectID || projects[0].PermissionLevel != 5 {
		t.Fatalf("Wrong return, got project:%v %v, perm: %v", projects[0].ProjectName, projects[0].ProjectID, projects[0].PermissionLevel)
	}

	MySQLProjectRevokePermission(projectID, "fahslaj")
	_ = MySQLProjectDelete(projectID)
	_ = mySQLUserDelete("jshap70", "secret")
	_ = mySQLUserDelete("fahslaj", "secret")

	projects, _ = MySQLUserProjects("fahslaj")
	if len(projects) > 0 {
		t.Fatalf("Projects returned not the correct length, expected: 0, actual: %v", len(projects))
	}
}

func TestMySQLProjectRename(t *testing.T) {
	configSetup()
	mySQLUserDelete("jshap70", "secret")

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	MySQLUserRegister(user)

	projectID, _ := MySQLProjectCreate("jshap70", "codecollabcore")

	err := MySQLProjectRename(projectID, "newName")
	if err != nil {
		t.Fatal(err)
	}

	projects, err := MySQLUserProjects("jshap70")
	_ = MySQLProjectDelete(projectID)
	_ = mySQLUserDelete("jshap70", "secret")

	if projects[0].ProjectID != projectID || projects[0].ProjectName != "newName" {
		t.Fatalf("Wrong return, got project:%v %v", projects[0].ProjectName, projects[0].ProjectID)
	}
}

func TestMySQLFileCreate(t *testing.T) {
	configSetup()

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	MySQLUserRegister(user)

	projectID, _ := MySQLProjectCreate("jshap70", "codecollabcore")
	fileID, err := MySQLFileCreate("jshap70", "file-y", ".", projectID)

	files, _ := MySQLProjectGetFiles(projectID)

	_ = MySQLProjectDelete(projectID)
	_ = mySQLUserDelete("jshap70", "secret")

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

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	MySQLUserRegister(user)

	projectID, _ := MySQLProjectCreate("jshap70", "codecollabcore")
	fileID, _ := MySQLFileCreate("jshap70", "file-y", ".", projectID)
	err := MySQLFileDelete(fileID)

	files, _ := MySQLProjectGetFiles(projectID)
	_ = MySQLProjectDelete(projectID)
	_ = mySQLUserDelete("jshap70", "secret")

	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Fatalf("Project %v returned not the correct length, expected: 0, actual: %v", projectID, len(files))
	}
}

func TestMySQLFileMove(t *testing.T) {
	configSetup()

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	MySQLUserRegister(user)

	projectID, _ := MySQLProjectCreate("jshap70", "codecollabcore")
	fileID, _ := MySQLFileCreate("jshap70", "file-y", ".", projectID)

	err := MySQLFileMove(fileID, "cc")

	files, _ := MySQLProjectGetFiles(projectID)
	_ = MySQLProjectDelete(projectID)
	_ = mySQLUserDelete("jshap70", "secret")

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

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	MySQLUserRegister(user)

	projectID, _ := MySQLProjectCreate("jshap70", "codecollabcore")
	fileID, _ := MySQLFileCreate("jshap70", "file-y", ".", projectID)

	err := MySQLFileRename(fileID, "file-z")

	files, _ := MySQLProjectGetFiles(projectID)
	_ = MySQLProjectDelete(projectID)
	_ = mySQLUserDelete("jshap70", "secret")

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

	user := UserMeta{
		Username:  "jshap70",
		Password:  "secret",
		Email:     "joel@codecollab.cc",
		FirstName: "Joel",
		LastName:  "Shapiro"}

	MySQLUserRegister(user)

	projectID, _ := MySQLProjectCreate("jshap70", "codecollabcore")
	fileID, _ := MySQLFileCreate("jshap70", "file-y", ".", projectID)

	filebefore, err := MySQLFileGetInfo(fileID)
	_ = MySQLFileMove(fileID, "cc")
	fileafter, err := MySQLFileGetInfo(fileID)

	_ = MySQLProjectDelete(projectID)
	_ = mySQLUserDelete("jshap70", "secret")

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
