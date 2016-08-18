package datahandling

import (
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
)

func configSetup() {
	config.SetConfigDir("../../config")
	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: this is backup for the tests, and will likely fail on a
	// non-local system unless the DB's have been set up to allow for this
	if val := config.GetConfig().ConnectionConfig["MySQL"].Schema; val == "" {
		tempCon := config.GetConfig().ConnectionConfig["MySQL"]
		tempCon.Schema = "testing"
		config.GetConfig().ConnectionConfig["MySQL"] = tempCon
	}
	if val := config.GetConfig().ConnectionConfig["Couchbase"].Schema; val == "" {
		tempCon := config.GetConfig().ConnectionConfig["Couchbase"]
		tempCon.Schema = "testing"
		config.GetConfig().ConnectionConfig["Couchbase"] = tempCon
	}
}

func setBaseFields(req request) {
	req.setAbstractRequest(&abstractRequest{
		SenderID:    "loganga",
		SenderToken: "supersecure",
	})
}

var datahanly = DataHandler{
	MessageChan:      make(chan rabbitmq.AMQPMessage, 1),
	SubscriptionChan: make(chan rabbitmq.Subscription, 1),
	WebsocketID:      1,
}

func TestProjectCreateRequest_Process(t *testing.T) {
	configSetup()
	req := *new(projectCreateRequest)
	setBaseFields(&req)

	req.Resource = "Project"
	req.Method = "Create"
	req.Name = "new stuff"

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

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	// did the server return success status
	if resp.Status != success {
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
	if project.ProjectName != "new stuff" || project.ProjectID != projectID {
		t.Fatal("wrong project added somehow")
	}

}

func TestProjectRenameRequest_Process(t *testing.T) {
	configSetup()
	req := *new(projectRenameRequest)
	setBaseFields(&req)

	req.Resource = "Project"
	req.Method = "Rename"
	req.ProjectID = 1
	req.NewName = "newer stuff"

	db := dbfs.NewDBMock()
	usermeta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	db.Users["loganga"] = usermeta
	projectmeta := dbfs.ProjectMeta{
		ProjectID:       req.ProjectID,
		ProjectName:     "new stuff",
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

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	not := closures[1].(toRabbitChannelClosure).msg.ServerMessage.(notification)
	// did the server return success status
	if resp.Status != success {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	// is the data actually correct
	if not.ResourceID != db.ProjectIDCounter-1 {
		t.Fatalf("Incorrect projectID was returned, expected %d, recieved %d", db.ProjectIDCounter-1, not.ResourceID)
	}

}

// projectGetPermissionConstantsRequest.process is unimplemented

func TestProjectGrantPermissionsRequest_Process(t *testing.T) {
	configSetup()
	req := *new(projectGrantPermissionsRequest)
	setBaseFields(&req)

	req.Resource = "Project"
	req.Method = "GrantPermissions"
	req.GrantUsername = "notloganga"
	req.PermissionLevel = 5

	db := dbfs.NewDBMock()
	genemeta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	notgenemeta := dbfs.UserMeta{
		FirstName: "Notgene",
		LastName:  "NotLogan",
		Email:     "notloganga@codecollaborate.com",
		Password:  "incorrect horse battery staple",
		Username:  "notloganga",
	}
	db.Users["loganga"] = genemeta
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
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toRabbitChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	not := closures[1].(toRabbitChannelClosure).msg.ServerMessage.(notification)
	// did the server return success status
	if resp.Status != success {
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
	configSetup()
	req := *new(projectRevokePermissionsRequest)
	setBaseFields(&req)

	req.Resource = "Project"
	req.Method = "RevokePermissions"
	req.RevokeUsername = "notloganga"

	db := dbfs.NewDBMock()
	genemeta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	notgenemeta := dbfs.UserMeta{
		FirstName: "Notgene",
		LastName:  "NotLogan",
		Email:     "notloganga@codecollaborate.com",
		Password:  "incorrect horse battery staple",
		Username:  "notloganga",
	}
	db.Users["loganga"] = genemeta
	db.Users["notloganga"] = notgenemeta

	projectID, err := db.MySQLProjectCreate("loganga", "new stuff")
	db.MySQLProjectGrantPermission(projectID, notgenemeta.Username, 5, genemeta.Username)

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

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	not := closures[1].(toRabbitChannelClosure).msg.ServerMessage.(notification)
	// did the server return success status
	if resp.Status != success {
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
	configSetup()
	req := *new(projectLookupRequest)
	setBaseFields(&req)
	db := dbfs.NewDBMock()

	req.Resource = "Project"
	req.Method = "Lookup"

	usermeta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	db.Users["loganga"] = usermeta

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

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	// did the server return success status
	if resp.Status != success {
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
	configSetup()
	req := *new(projectGetFilesRequest)
	setBaseFields(&req)
	db := dbfs.NewDBMock()

	req.Resource = "Project"
	req.Method = "GetFiles"

	usermeta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	db.Users["loganga"] = usermeta

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

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	// did the server return success status
	if resp.Status != success {
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
	configSetup()
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
		reflect.TypeOf(closures[0]).String() != "datahandling.rabbitChannelSubscribeClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	sub := closures[0].(rabbitChannelSubscribeClosure)
	// did the server return success status
	val, err := strconv.ParseInt(sub.key, 10, 64)
	if val != req.ProjectID {
		t.Fatalf("Subscribe function wanted to subscribe to the wrong channel\n expected: %d, got: %d", req.ProjectID, val)
	}
}

func TestProjectUnsubscribe_Process(t *testing.T) {
	configSetup()
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
		reflect.TypeOf(closures[0]).String() != "datahandling.rabbitChannelUnsubscribeClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	sub := closures[0].(rabbitChannelUnsubscribeClosure)
	// did the server return success status
	val, err := strconv.ParseInt(sub.key, 10, 64)
	if val != req.ProjectID {
		t.Fatalf("Subscribe function wanted to subscribe to the wrong channel\n expected: %d, got: %d", req.ProjectID, val)
	}
}

func TestProjectDeleteRequest_process(t *testing.T) {
	configSetup()
	req := *new(projectDeleteRequest)
	setBaseFields(&req)

	req.Resource = "Project"
	req.Method = "Delete"

	db := dbfs.NewDBMock()
	usermeta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	db.Users["loganga"] = usermeta
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

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	not := closures[1].(toRabbitChannelClosure).msg.ServerMessage.(notification)
	// did the server return success status
	if resp.Status != success {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	// is the data actually correct
	if not.ResourceID != projID {
		t.Fatalf("Incorrect projectID was returned, expected %d, recieved %d", projID, not.ResourceID)
	}
}
