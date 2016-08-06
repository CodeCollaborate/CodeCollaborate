package datahandling

import (
	"log"
	"testing"

	"reflect"

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

// TODO (testing/required): testing... lots of testing

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
		t.Fatal("did not properly process")
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
		reflect.TypeOf(closures[1]).String() != "datahandling.toChannelClosure" {
		t.Fatal("did not properly process")
	}

	resp := closures[0].(toSenderClosure).msg.ServerMessage.(response)
	not := closures[1].(toChannelClosure).msg.ServerMessage.(notification)
	// did the server return success status
	if resp.Status != success {
		t.Fatalf("Process function responded with status: %d", resp.Status)
	}

	// is the data actually correct
	if not.ResourceID != db.ProjectIDCounter-1 {
		t.Fatalf("Incorrect projectID was returned, expected %d, recieved %d", db.ProjectIDCounter-1, not.ResourceID)
	}

}

func TestProjectLookupRequest_Process(t *testing.T) {
	configSetup()
	req := *new(projectLookupRequest)
	setBaseFields(&req)

	req.Resource = "Project"
	req.Method = "Lookup"
	req.ProjectIDs = []int64{12345, 38292}

	db := dbfs.NewDBMock()

	continuations, err := req.process(db)

	if err != nil {
		t.Fatal(err)
	}

	if len(continuations) != 1 {
		t.Fatal("did not properly process")
	}
}
