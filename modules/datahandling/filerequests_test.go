package datahandling

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/CodeCollaborate/Server/modules/datahandling/messages"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/stretchr/testify/assert"
)

var geneMeta = dbfs.UserMeta{
	FirstName: "Gene",
	LastName:  "Logan",
	Email:     "loganga@codecollaborate.com",
	Password:  "correct horse battery staple",
	Username:  "loganga",
}

func TestFileCreateRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(fileCreateRequest)
	setBaseFields(&req)

	db := dbfs.NewDBMock()
	db.MySQLUserRegister(geneMeta)
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
	assert.Equal(t, 4, db.FunctionCallCount, "did not call correct number of db functions")

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toRabbitChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	closure := closures[1].(toRabbitChannelClosure)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}
	// is the data actually correct
	FileID := reflect.ValueOf(resp.Data).FieldByName("FileID").Interface().(int64)

	if closure.key != fmt.Sprintf("Project-%d", projectid) {
		t.Fatal("notification sent to wrong channel")
	}

	notFileID := reflect.ValueOf(closure.msg.ServerMessage.(messages.Notification).Data).FieldByName("File").FieldByName("FileID").Interface().(int64)
	if FileID != notFileID {
		t.Fatal("recieved different data from notification and response")
	}
}

func TestFileRenameRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(fileRenameRequest)
	setBaseFields(&req)

	db := dbfs.NewDBMock()
	db.MySQLUserRegister(geneMeta)
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
	assert.Equal(t, 4, db.FunctionCallCount, "did not call correct number of db functions")

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toRabbitChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	closure := closures[1].(toRabbitChannelClosure)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	if closure.key != fmt.Sprintf("Project-%d", projectid) {
		t.Fatal("notification sent to wrong channel")
	}

	notFileID := closure.msg.ServerMessage.(messages.Notification).ResourceID
	if fileid != notFileID {
		t.Fatal("wrong FileID recieved in notification")
	}

	filename := reflect.ValueOf(closure.msg.ServerMessage.(messages.Notification).Data).FieldByName("NewName").Interface().(string)
	if filename != req.NewName {
		t.Fatal("wrong new filename recieved in notification")
	}

	// TODO: check the file actually moved

}

func TestFileMoveRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(fileMoveRequest)
	setBaseFields(&req)

	db := dbfs.NewDBMock()
	db.MySQLUserRegister(geneMeta)
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
	assert.Equal(t, 4, db.FunctionCallCount, "did not call correct number of db functions")

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toRabbitChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	closure := closures[1].(toRabbitChannelClosure)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	if closure.key != fmt.Sprintf("Project-%d", projectid) {
		t.Fatal("notification sent to wrong channel")
	}

	notFileID := closure.msg.ServerMessage.(messages.Notification).ResourceID
	if fileid != notFileID {
		t.Fatal("wrong FileID recieved in notification")
	}

	filepath := reflect.ValueOf(closure.msg.ServerMessage.(messages.Notification).Data).FieldByName("NewPath").Interface().(string)
	if filepath != req.NewPath {
		t.Fatal("wrong new filepath recieved in notification")
	}

	// TODO: check the file actually moved

}

func TestFileDeleteRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(fileDeleteRequest)
	setBaseFields(&req)

	db := dbfs.NewDBMock()
	db.MySQLUserRegister(geneMeta)
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
	assert.Equal(t, 5, db.FunctionCallCount, "did not call correct number of db functions")

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toRabbitChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	closure := closures[1].(toRabbitChannelClosure)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	if closure.key != fmt.Sprintf("Project-%d", projectid) {
		t.Fatal("notification sent to wrong channel")
	}

	notFileID := closure.msg.ServerMessage.(messages.Notification).ResourceID
	if fileid != notFileID {
		t.Fatal("wrong FileID recieved in notification")
	}

	if _, ok := db.Files[fileid]; ok {
		t.Fatal("File still exists")
	}

}

func TestFileChangeRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(fileChangeRequest)
	setBaseFields(&req)

	db := dbfs.NewDBMock()
	db.MySQLUserRegister(geneMeta)
	projectid, err := db.MySQLProjectCreate("loganga", "hi")
	fileid, err := db.MySQLFileCreate("loganga", "new file", "", projectid)
	db.CBInsertNewFile(fileid, newFileVersion, []string{})

	req.Resource = "File"
	req.Method = "Change"
	req.FileID = fileid
	req.Changes = "v0:\n0:+1:a:\n10"
	baseFileVersion := int64(1)

	db.FunctionCallCount = 0

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	assert.Equal(t, 3, db.FunctionCallCount, "did not call correct number of db functions")

	// are we notifying the right people
	if len(closures) != 2 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" ||
		reflect.TypeOf(closures[1]).String() != "datahandling.toRabbitChannelClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	closure := closures[1].(toRabbitChannelClosure)
	// did the server return success status
	if resp.Status != messages.StatusSuccess {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	if closure.key != fmt.Sprintf("Project-%d", projectid) {
		t.Fatal("notification sent to wrong channel")
	}

	notFileID := closure.msg.ServerMessage.(messages.Notification).ResourceID
	if fileid != notFileID {
		t.Fatal("wrong FileID recieved in notification")
	}

	changes := reflect.ValueOf(closure.msg.ServerMessage.(messages.Notification).Data).FieldByName("Changes").Interface().(string)
	if changes != req.Changes {
		t.Fatal("wrong changes recieved in notification")
	}

	if db.FileChanges[fileid][0] != changes {
		t.Fatal("changes not inserted")
	}

	newVersion := reflect.ValueOf(closure.msg.ServerMessage.(messages.Notification).Data).FieldByName("FileVersion").Interface().(int64)
	if newVersion != baseFileVersion+1 {
		t.Fatalf("wrong file version, expected: %d, got: %d", baseFileVersion+1, newVersion)
	}

	// try the request again to prove that it rejects higher file versions

	req.Changes = "v9999:\n0:+1:a:\n10"
	db.FunctionCallCount = 0

	closures, err = req.process(db)
	if err != dbfs.ErrVersionOutOfDate {
		t.Fatal(err)
	}

	// didn't call extra db functions
	assert.Equal(t, 3, db.FunctionCallCount, "did not call correct number of db functions")

	// are we notifying the right people
	if len(closures) != 1 ||
		reflect.TypeOf(closures[0]).String() != "datahandling.toSenderClosure" {
		t.Fatalf("did not properly process, recieved %d closure(s)", len(closures))
	}

	resp = closures[0].(toSenderClosure).msg.ServerMessage.(messages.Response)
	// did the server return out of date status
	if resp.Status != messages.StatusVersionOutOfDate {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

}

func TestFilePullRequest_Process(t *testing.T) {
	configSetup(t)
	req := *new(filePullRequest)
	setBaseFields(&req)

	db := dbfs.NewDBMock()
	db.MySQLUserRegister(geneMeta)
	projectID, err := db.MySQLProjectCreate("loganga", "hi")
	fileid, err := db.MySQLFileCreate("loganga", "new file", "", projectID)
	db.FileWrite("./", "new file", projectID, []byte{})

	changes := "v0:\n0:+1:a:\n10"
	db.CBAppendFileChange(dbfs.FileMeta{FileID: fileid}, changes)

	req.Resource = "File"
	req.Method = "Pull"
	req.FileID = fileid

	db.FunctionCallCount = 0

	closures, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 3 {
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
	fileChanges := reflect.ValueOf(resp.Data).FieldByName("Changes").Interface().([]string)
	if changes != fileChanges[0] {
		t.Fatalf("wrong file changes, expected: %v, got: %v", changes, fileChanges)
	}
}
