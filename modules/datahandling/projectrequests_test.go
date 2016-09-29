package datahandling

import (
	"reflect"
	"testing"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datahandling/messages"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
)

func setBaseFields(req request) {
	req.setAbstractRequest(&abstractRequest{
		SenderID:    "loganga",
		SenderToken: "supersecure",
	})
}

var datahanly = DataHandler{
	MessageChan: make(chan rabbitmq.AMQPMessage, 1),
	WebsocketID: 1,
}

func TestProjectCreateRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(projectCreateRequest)
	setBaseFields(&req)

	req.Resource = "Project"
	req.Method = "Create"
	req.Name = "new stuff"

	db := dbfs.NewDBMock()
	db.Users["loganga"] = geneMeta

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

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}
	// is the data actually correct
	projectID := reflect.ValueOf(resp.Data).FieldByName("ProjectID").Interface().(int64)
	if projectID != db.ProjectIDCounter-1 {
		t.Fatal("Incorrect projectID was returned")
	}

	if len(db.Projects["loganga"]) != 1 {
		t.Fatal("did not actually add project")
	}

	project := db.Projects["loganga"][0]
	if project.Name != "new stuff" || project.ProjectID != projectID {
		t.Fatal("wrong project added somehow")
	}

}

func TestProjectRenameRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(projectRenameRequest)
	setBaseFields(&req)

	req.Resource = "Project"
	req.Method = "Rename"
	req.ProjectID = 1
	req.NewName = "newer stuff"

	db := dbfs.NewDBMock()
	db.Users["loganga"] = geneMeta
	projectmeta := dbfs.ProjectMeta{
		ProjectID:       req.ProjectID,
		Name:            "new stuff",
		PermissionLevel: 10,
	}
	db.Projects["loganga"] = []dbfs.ProjectMeta{projectmeta}
	db.ProjectIDCounter = 2

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 1 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toRabbitChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	not := closures[1].(toRabbitChannelClosure).msg.ServerMessage.(messages.Notification)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	// is the data actually correct
	if not.ResourceID != db.ProjectIDCounter-1 {
		t.Fatalf("Incorrect projectID was returned, expected %d, recieved %d", db.ProjectIDCounter-1, not.ResourceID)
	}

}

// projectGetPermissionConstantsRequest.process is unimplemented

func TestProjectGrantPermissionsRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(projectGrantPermissionsRequest)
	setBaseFields(&req)

	permlvl, _ := config.GetPermissionLevel("Write")

	req.Resource = "Project"
	req.Method = "GrantPermissions"
	req.GrantUsername = "notloganga"
	req.PermissionLevel = permlvl

	db := dbfs.NewDBMock()
	notgenemeta := dbfs.UserMeta{
		FirstName: "Notgene",
		LastName:  "NotLogan",
		Email:     "notloganga@codecollaborate.com",
		Password:  "incorrect horse battery staple",
		Username:  "notloganga",
	}
	db.Users["loganga"] = geneMeta
	db.Users["notloganga"] = notgenemeta

	projectID, err := db.MySQLProjectCreate("loganga", "new stuff")

	db.FunctionCallCount = 0
	req.ProjectID = projectID

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 1 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(closures) != 3 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toRabbitChannelClosure" ||
		reflect.TypeOf(closures[2]).String() != "datahandling.toRabbitChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	not := closures[1].(toRabbitChannelClosure).msg.ServerMessage.(messages.Notification)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	// is the data actually correct
	username := reflect.ValueOf(not.Data).FieldByName("GrantUsername").Interface().(string)
	if username != req.GrantUsername {
		t.Fatalf("Incorrect username was returned, expected %v, recieved %v", req.GrantUsername, username)
	}

	// did the user actually get added
	if len(db.Projects[req.GrantUsername]) != 1 || db.Projects[req.GrantUsername][0].PermissionLevel != req.PermissionLevel {
		t.Fatal("Database was not properly modified")
	}

}

func TestProjectRevokePermissionsRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(projectRevokePermissionsRequest)
	setBaseFields(&req)

	req.Resource = "Project"
	req.Method = "RevokePermissions"
	req.RevokeUsername = "notloganga"

	db := dbfs.NewDBMock()
	notgenemeta := dbfs.UserMeta{
		FirstName: "Notgene",
		LastName:  "NotLogan",
		Email:     "notloganga@codecollaborate.com",
		Password:  "incorrect horse battery staple",
		Username:  "notloganga",
	}
	db.Users["loganga"] = geneMeta
	db.Users["notloganga"] = notgenemeta

	projectID, err := db.MySQLProjectCreate("loganga", "new stuff")
	db.MySQLProjectGrantPermission(projectID, notgenemeta.Username, 5, geneMeta.Username)

	db.FunctionCallCount = 0
	req.ProjectID = projectID

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 1 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toRabbitChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	not := closures[1].(toRabbitChannelClosure).msg.ServerMessage.(messages.Notification)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	// is the data actually correct
	username := reflect.ValueOf(not.Data).FieldByName("RevokeUsername").Interface().(string)
	if username != req.RevokeUsername {
		t.Fatalf("Incorrect username was returned, expected %v, recieved %v", req.RevokeUsername, username)
	}

	// did the user actually get removed
	if len(db.Projects[req.RevokeUsername]) != 0 {
		t.Fatal("Database was not properly modified")
	}

}

// projectGetOnlineClientsRequest.process is unimplemented

func TestProjectLookupRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(projectLookupRequest)
	setBaseFields(&req)
	db := dbfs.NewDBMock()

	req.Resource = "Project"
	req.Method = "Lookup"

	db.Users["loganga"] = geneMeta

	projid1, err := db.MySQLProjectCreate("loganga", "new shit")
	projid2, err := db.MySQLProjectCreate("loganga", "newer shit")

	req.ProjectIDs = []int64{projid1, projid2}
	db.FunctionCallCount = 0

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != len(req.ProjectIDs) {
		t.Fatalf("did not call correct number of db functions, expected: %d, actual: %d", len(req.ProjectIDs), db.FunctionCallCount)
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
	if len(projects) != 2 {
		t.Fatalf("Incorrect project count, expected %d, recieved %d", 2, len(projects))
	}

	if projects[0].Name != "new shit" || projects[1].Name != "newer shit" {
		t.Fatal("incorrect project name(s)")
	}

}

func TestProjectGetFilesRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(projectGetFilesRequest)
	setBaseFields(&req)
	db := dbfs.NewDBMock()

	req.Resource = "Project"
	req.Method = "GetFiles"

	db.Users["loganga"] = geneMeta

	projid1, err := db.MySQLProjectCreate("loganga", "new shit")
	db.MySQLFileCreate("loganga", "file1", "", projid1)
	db.MySQLFileCreate("loganga", "file2", "", projid1)
	db.MySQLFileCreate("loganga", "file3", "", projid1)

	req.ProjectID = projid1
	db.FunctionCallCount = 0

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	// NOTE: 4 comes from 1 db.MySQLProjectGetFiles + 3 db.CBGetFileVersion
	if db.FunctionCallCount != 4 {
		t.Fatalf("did not call correct number of db functions, expected: %d, actual: %d", 4, db.FunctionCallCount)
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
	files := reflect.ValueOf(resp.Data).FieldByName("Files").Interface().([]fileLookupResult)
	if len(files) != 3 {
		t.Fatalf("Incorrect file count, expected %d, recieved %d", 3, len(files))
	}

	if files[0].Filename != "file1" || files[1].Filename != "file2" || files[2].Filename != "file3" {
		t.Fatal("incorrect filename(s)")
	}

}

func TestProjectSubscribe_Process(t *testing.T) {
	configSetup(t)
	req := *new(projectSubscribeRequest)
	setBaseFields(&req)
	db := dbfs.NewDBMock()

	req.Resource = "Project"
	req.Method = "Subscribe"
	req.ProjectID = 1

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// are we notifying the right people
	if len(closures) != 1 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.rabbitCommandClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	sub := closures[0].(rabbitCommandClosure)
	// did the server return success status
	channelKey := rabbitmq.RabbitProjectQueueName(req.ProjectID)
	if sub.Data.(rabbitmq.RabbitQueueData).Key != channelKey {
		t.Fatalf("Subscribe function wanted to subscribe to the wrong channel\n expected: %s, got: %s", channelKey, sub.Data.(rabbitmq.RabbitQueueData).Key)
	}
}

func TestProjectUnsubscribe_Process(t *testing.T) {
	configSetup(t)
	req := *new(projectUnsubscribeRequest)
	setBaseFields(&req)
	db := dbfs.NewDBMock()

	req.Resource = "Project"
	req.Method = "Unsubscribe"
	req.ProjectID = 1

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// are we notifying the right people
	if len(closures) != 1 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.rabbitCommandClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	sub := closures[0].(rabbitCommandClosure)
	// did the server return success status
	channelKey := rabbitmq.RabbitProjectQueueName(req.ProjectID)
	if sub.Data.(rabbitmq.RabbitQueueData).Key != channelKey {
		t.Fatalf("Subscribe function wanted to subscribe to the wrong channel\n expected: %s, got: %s", channelKey, sub.Data.(rabbitmq.RabbitQueueData).Key)
	}
}

func TestProjectDeleteRequest_process(t *testing.T) {
	configSetup(t)
	req := *new(projectDeleteRequest)
	setBaseFields(&req)

	req.Resource = "Project"
	req.Method = "Delete"

	db := dbfs.NewDBMock()
	db.Users["loganga"] = geneMeta
	projID, err := db.MySQLProjectCreate("loganga", "new project")

	db.FunctionCallCount = 0
	req.ProjectID = projID

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 1 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toRabbitChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	not := closures[1].(toRabbitChannelClosure).msg.ServerMessage.(messages.Notification)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	// is the data actually correct
	if not.ResourceID != projID {
		t.Fatalf("Incorrect projectID was returned, expected %d, recieved %d", projID, not.ResourceID)
	}
}
