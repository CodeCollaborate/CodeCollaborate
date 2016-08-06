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
		t.Fatal("did not properly process")
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
