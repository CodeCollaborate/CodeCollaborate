package datahandling

import (
	"reflect"
	"testing"

	"github.com/CodeCollaborate/Server/modules/dbfs"
)

func init() {
	disableAuth = true
}

func TestUserRegisterRequest_Process(t *testing.T) {
	configSetupUnauthenticated(t)
	req := *new(userRegisterRequest)
	setBaseFields(&req)

	req.Resource = "User"
	req.Method = "Register"

	req.Username = "loganga"
	req.FirstName = "Gene"
	req.LastName = "Logan"
	req.Email = "loganga@codecollaborate.com"
	req.Password = "correct horse battery staple"

	db := dbfs.NewDBMock()
	datahanly.Db = db

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 1 {
		t.Fatal("did not call correct number of db functions")
	}
	// did gene it actually added
	if _, ok := db.Users["loganga"]; !ok {
		t.Fatal("did not correctly call db function")
	}

	// are we notifying the right people
	if len(closures) != 1 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}
	// did the server return success status
	cont := closures[0].(toSenderClosure).msg.ServerMessage.(response).Status
	if cont != success {
		t.Fatalf("Process function responded with status: %d", cont)
	}

	closures, err = req.process(db)
	if err == nil {
		t.Fatal("Should have failed to register user that already exists")
	}
}

// userLoginRequest.process is unimplemented

func TestUserLookupRequest_Process(t *testing.T) {
	configSetupUnauthenticated(t)
	req := *new(userLookupRequest)
	setBaseFields(&req)

	req.Resource = "User"
	req.Method = "Lookup"

	req.Usernames = []string{"loganga"}

	db := dbfs.NewDBMock()
	meta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	db.Users["loganga"] = meta

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 1 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(closures) != 1 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}
	response := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	// did the server return success status
	if response.Status != success {
		t.Fatalf("Process function responded with status: %d", response.Status)
	}
	// is the data actually correct
	users := reflect.ValueOf(response.Data).FieldByName("Users").Interface().([]dbfs.UserMeta)
	if len(users) != 1 && users[0] != meta {
		t.Fatal("Incorrect user was returned")
	}
}

func TestUserProjectsRequest_Process(t *testing.T) {
	configSetupUnauthenticated(t)
	req := *new(userProjectsRequest)
	setBaseFields(&req)

	req.Resource = "User"
	req.Method = "Projects"

	db := dbfs.NewDBMock()
	gene := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	db.MySQLUserRegister(gene)

	notgene := dbfs.UserMeta{
		FirstName: "Not",
		LastName:  "Gene",
		Email:     "notloganga@codecollaborate.com",
		Password:  "incorrect horse battery staple",
		Username:  "notloganga",
	}
	db.MySQLUserRegister(notgene)

	db.MySQLProjectCreate("loganga", "my project")
	genesproject := db.Projects["loganga"][0]

	db.MySQLProjectCreate("notloganga", "not his project")
	notgenesproject := db.Projects["notloganga"][0]

	db.FunctionCallCount = 0

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 1 {
		t.Fatalf("did not call correct number of db functions, called %d # of arguments", db.FunctionCallCount)
	}

	// are we notifying the right people
	if len(closures) != 1 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	// did the server return success status
	if resp.Status != success {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}
	// is the data actually correct
	projects := reflect.ValueOf(resp.Data).FieldByName("Projects").Interface().([]dbfs.ProjectMeta)
	if len(projects) != 1 && projects[0] != genesproject {
		t.Fatal("Incorrect user was returned")
	}

	// add gene to a new project and see if the process function updates as expected
	db.MySQLProjectGrantPermission(notgenesproject.ProjectID, "loganga", 5, "notloganga")
	db.FunctionCallCount = 0

	closures, err = req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// are we notifying the right people
	if len(closures) != 1 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp = closures[0].(toSenderClosure).msg.ServerMessage.(response)
	// did the server return success status
	if resp.Status != success {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}
	// is the data actually correct
	projects = reflect.ValueOf(resp.Data).FieldByName("Projects").Interface().([]dbfs.ProjectMeta)
	if len(projects) != 2 && projects[0] != genesproject && projects[1] != notgenesproject {
		t.Fatal("Incorrect user was returned")
	}
}
