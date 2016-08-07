package datahandling

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/CodeCollaborate/Server/modules/dbfs"
)

func TestFileCreateRequest_Process(t *testing.T) {
	configSetup()
	req := *new(fileCreateRequest)
	setBaseFields(&req)

	db := dbfs.NewDBMock()
	meta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	db.MySQLUserRegister(meta)
	projectid, err := db.MySQLProjectCreate("loganga", "hi")

	req.Resource = "File"
	req.Method = "Create"
	req.Name = "new file"
	req.ProjectID = projectid
	req.RelativePath = ""
	req.FileBytes = []byte{}

	db.FunctionCallCount = 0

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 2 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	not := closures[1].(toChannelClosure).msg
	// did the server return success status
	if resp.Status != success {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}
	// is the data actually correct
	FileID := reflect.ValueOf(resp.Data).FieldByName("FileID").Interface().(int64)

	route, err := strconv.ParseInt(not.RoutingKey, 10, 64)
	if route != projectid {
		t.Fatal("notification sent to wrong channel")
	}

	notFileID := reflect.ValueOf(not.ServerMessage.(notification).Data).FieldByName("FileID").Interface().(int64)
	if FileID != notFileID {
		t.Fatal("recieved different data from notification and responce")
	}

}

func TestFileRenameRequest_Process(t *testing.T) {
	configSetup()
	req := *new(fileRenameRequest)
	setBaseFields(&req)

	db := dbfs.NewDBMock()
	meta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	db.MySQLUserRegister(meta)
	projectid, err := db.MySQLProjectCreate("loganga", "hi")
	fileid, err := db.MySQLFileCreate("loganga", "new file", "", projectid)

	req.Resource = "File"
	req.Method = "Rename"
	req.FileID = fileid
	req.NewName = "new name"

	db.FunctionCallCount = 0

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 2 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	not := closures[1].(toChannelClosure).msg
	// did the server return success status
	if resp.Status != success {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	route, err := strconv.ParseInt(not.RoutingKey, 10, 64)
	if route != projectid {
		t.Fatal("notification sent to wrong channel")
	}

	notFileID := not.ServerMessage.(notification).ResourceID
	if fileid != notFileID {
		t.Fatal("wrong FileID recieved in notification")
	}

	filename := reflect.ValueOf(not.ServerMessage.(notification).Data).FieldByName("NewName").Interface().(string)
	if filename != req.NewName {
		t.Fatal("wrong new filename recieved in notification")
	}

	// TODO: check the file actually moved

}

func TestFileMoveRequest_Process(t *testing.T) {
	configSetup()
	req := *new(fileMoveRequest)
	setBaseFields(&req)

	db := dbfs.NewDBMock()
	meta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	db.MySQLUserRegister(meta)
	projectid, err := db.MySQLProjectCreate("loganga", "hi")
	fileid, err := db.MySQLFileCreate("loganga", "new file", "", projectid)

	req.Resource = "File"
	req.Method = "Move"
	req.FileID = fileid
	req.NewPath = "random/"

	db.FunctionCallCount = 0

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 2 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	not := closures[1].(toChannelClosure).msg
	// did the server return success status
	if resp.Status != success {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	route, err := strconv.ParseInt(not.RoutingKey, 10, 64)
	if route != projectid {
		t.Fatal("notification sent to wrong channel")
	}

	notFileID := not.ServerMessage.(notification).ResourceID
	if fileid != notFileID {
		t.Fatal("wrong FileID recieved in notification")
	}

	filepath := reflect.ValueOf(not.ServerMessage.(notification).Data).FieldByName("NewPath").Interface().(string)
	if filepath != req.NewPath {
		t.Fatal("wrong new filepath recieved in notification")
	}

	// TODO: check the file actually moved

}

func TestFileDeleteRequest_Process(t *testing.T) {
	configSetup()
	req := *new(fileDeleteRequest)
	setBaseFields(&req)

	db := dbfs.NewDBMock()
	meta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	db.MySQLUserRegister(meta)
	projectid, err := db.MySQLProjectCreate("loganga", "hi")
	fileid, err := db.MySQLFileCreate("loganga", "new file", "", projectid)

	req.Resource = "File"
	req.Method = "Delete"
	req.FileID = fileid

	db.FunctionCallCount = 0

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 4 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	not := closures[1].(toChannelClosure).msg
	// did the server return success status
	if resp.Status != success {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	route, err := strconv.ParseInt(not.RoutingKey, 10, 64)
	if route != projectid {
		t.Fatal("notification sent to wrong channel")
	}

	notFileID := not.ServerMessage.(notification).ResourceID
	if fileid != notFileID {
		t.Fatal("wrong FileID recieved in notification")
	}

	if _, ok := db.Files[fileid]; ok {
		t.Fatal("File still exists")
	}

}

func TestFileChangeRequest_Process(t *testing.T) {
	configSetup()
	req := *new(fileChangeRequest)
	setBaseFields(&req)

	db := dbfs.NewDBMock()
	meta := dbfs.UserMeta{
		FirstName: "Gene",
		LastName:  "Logan",
		Email:     "loganga@codecollaborate.com",
		Password:  "correct horse battery staple",
		Username:  "loganga",
	}
	db.MySQLUserRegister(meta)
	projectid, err := db.MySQLProjectCreate("loganga", "hi")
	fileid, err := db.MySQLFileCreate("loganga", "new file", "", projectid)
	db.CBInsertNewFile(fileid, newFileVersion, []string{})

	req.Resource = "File"
	req.Method = "Change"
	req.FileID = fileid
	req.Changes = []string{"change 1"}
	req.BaseFileVersion = newFileVersion

	db.FunctionCallCount = 0

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 2 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	not := closures[1].(toChannelClosure).msg
	// did the server return success status
	if resp.Status != success {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	route, err := strconv.ParseInt(not.RoutingKey, 10, 64)
	if route != projectid {
		t.Fatal("notification sent to wrong channel")
	}

	notFileID := not.ServerMessage.(notification).ResourceID
	if fileid != notFileID {
		t.Fatal("wrong FileID recieved in notification")
	}

	changes := reflect.ValueOf(not.ServerMessage.(notification).Data).FieldByName("Changes").Interface().([]string)
	if changes[0] != req.Changes[0] {
		t.Fatal("wrong changes recieved in notification")
	}

	if db.FileChanges[fileid][0] != changes[0] {
		t.Fatal("changes not inserted")
	}

	newVersion := reflect.ValueOf(not.ServerMessage.(notification).Data).FieldByName("FileVersion").Interface().(int64)
	if newVersion != req.BaseFileVersion+1 {
		t.Fatalf("wrong file version, expected: %d, got: %d", req.BaseFileVersion+1, newVersion)
	}

	// try the request again to prove that it rejects lower file versions

	db.FunctionCallCount = 0

	closures, err = req.process(db)
	if err != dbfs.ErrVersionOutOfDate {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 2 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(closures) != 1 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp = closures[0].(toSenderClosure).msg.ServerMessage.(response)
	// did the server return out of date status
	if resp.Status != versionOutOfDate {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

}
