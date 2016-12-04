package datahandling

import (
	"reflect"
	"testing"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datahandling/messages"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/stretchr/testify/assert"
)

func TestUserRegisterRequest_Process(t *testing.T) {
	configSetup(t)
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
	cont := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response).Status
	if cont != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", cont)
	}

	closures, err = req.process(db)
	if err == nil {
		t.Fatal("Should have failed to register user that already exists")
	}
}

// userLoginRequest.process is unimplemented

func TestUserDeleteRequest_Process(t *testing.T) {
	configSetup(t)

	req := *new(userDeleteRequest)
	setBaseFields(&req)

	req.Resource = "User"
	req.Method = "Delete"

	db := dbfs.NewDBMock()
	db.MySQLUserRegister(geneMeta)
	db.FunctionCallCount = 0

	closures, err := req.process(db)
	assert.Nil(t, err)
	assert.Equal(t, 2, db.FunctionCallCount, "unexpected db calls for user delete")

	assert.Equal(t, 1, len(closures), "unexpected number of returned closures")
	assert.IsType(t, toSenderClosure{}, closures[0], "incorrect closure type")

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)

	assert.Equal(t, messages.StatusSuccess, resp.Status, "unexpected response status")

	// test with projects
	req = *new(userDeleteRequest)
	setBaseFields(&req)

	req.Resource = "User"
	req.Method = "Delete"

	db = dbfs.NewDBMock()
	db.MySQLUserRegister(geneMeta)
	projectID1, _ := db.MySQLProjectCreate(geneMeta.Username, "_test_project1")
	projectID2, _ := db.MySQLProjectCreate(geneMeta.Username, "_test_project2")

	db.FunctionCallCount = 0

	closures, err = req.process(db)
	assert.Nil(t, err)
	assert.Equal(t, 2, db.FunctionCallCount, "unexpected db calls for user delete")

	assert.Equal(t, 3, len(closures), "unexpected number of returned closures")
	assert.IsType(t, toSenderClosure{}, closures[0], "incorrect closure type")
	assert.IsType(t, toRabbitChannelClosure{}, closures[1], "incorrect closure type")
	assert.IsType(t, toRabbitChannelClosure{}, closures[2], "incorrect closure type")

	resp = closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	assert.Equal(t, messages.StatusSuccess, resp.Status, "unexpected response status")

	not1 := closures[1].(toRabbitChannelClosure).msg.ServerMessage.(messages.Notification)
	assert.Equal(t, "Project", not1.Resource, "unexpected notification resource")
	assert.Equal(t, "Delete", not1.Method, "unexpected notification method")
	assert.Equal(t, projectID1, not1.ResourceID, "unexpected projectID deleted")

	not2 := closures[2].(toRabbitChannelClosure).msg.ServerMessage.(messages.Notification)
	assert.Equal(t, "Project", not2.Resource, "unexpected notification resource")
	assert.Equal(t, "Delete", not2.Method, "unexpected notification method")
	assert.Equal(t, projectID2, not2.ResourceID, "unexpected projectID deleted")
}

func TestUserLookupRequest_Process(t *testing.T) {
	configSetup(t)
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
	response := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	// did the server return success status
	if response.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", response.Status)
	}
	// is the data actually correct
	users := reflect.ValueOf(response.Data).FieldByName("Users").Interface().([]dbfs.UserMeta)
	if len(users) != 1 && users[0] != meta {
		t.Fatal("Incorrect user was returned")
	}
}

func TestUserProjectsRequest_Process(t *testing.T) {
	configSetup(t)
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
	if db.FunctionCallCount != 2 {
		t.Fatalf("did not call correct number of db functions, called %d # of arguments", db.FunctionCallCount)
	}

	// are we notifying the right people
	if len(closures) != 1 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}
	// is the data actually correct
	projects := reflect.ValueOf(resp.Data).FieldByName("Projects").Interface().([]projectLookupResult)
	if len(projects) != 1 && projects[0].ProjectID != genesproject.ProjectID {
		t.Fatal("Incorrect user was returned")
	}

	writePerm, err := config.PermissionByLabel("write")
	assert.Nil(t, err, "did not get permission")

	// add gene to a new project and see if the process function updates as expected
	err = db.MySQLProjectGrantPermission(notgenesproject.ProjectID, "loganga", writePerm.Level, "notloganga")
	assert.Nil(t, err, "couldn't grant project permission")

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

	resp = closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}
	// is the data actually correct
	projects = reflect.ValueOf(resp.Data).FieldByName("Projects").Interface().([]projectLookupResult)
	if len(projects) != 2 && projects[0].ProjectID != genesproject.ProjectID && projects[1].ProjectID != notgenesproject.ProjectID {
		t.Fatal("Incorrect user was returned")
	}

	assert.Len(t, projects[1].Permissions, 2, "incorrect number of permissions returned")

	// check to see if permission map is correct
	if projects[1].Permissions[gene.Username].PermissionLevel != writePerm.Level {
		t.Fatal("Permission map was not returned correctly in lookup")
	}

	ownerPerm, err := config.PermissionByLabel("owner")
	assert.Nil(t, err, "did not get permission")
	assert.Equal(t, ownerPerm.Level, projects[1].Permissions[notgene.Username].PermissionLevel, "not all permissions returned for project")
}
