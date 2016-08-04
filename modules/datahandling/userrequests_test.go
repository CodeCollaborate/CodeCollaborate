package datahandling

import (
	"reflect"
	"testing"

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
	MessageChan:      make(chan rabbitmq.AMQPMessage, 1),
	SubscriptionChan: make(chan rabbitmq.Subscription, 1),
	WebsocketID:      1,
}

func TestUserRegisterRequest_Process(t *testing.T) {
	configSetup()
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

	continuations, err := req.process(db)
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
	if len(continuations) != 1 ||
		reflect.TypeOf(continuations[0]).String() != "datahandling.toSenderClos" {
		t.Fatal("did not properly process")
	}
	// did the server return success status
	cont := continuations[0].(toSenderClos).msg.ServerMessage.(response).Status
	if cont != success {
		t.Fatalf("Process function responded with status: %d", cont)
	}

	continuations, err = req.process(db)
	if err == nil {
		t.Fatal("Should have failed to register user that already exists")
	}
}

// userLoginRequest.process is unimplemented

func TestUserLookupRequest_Process(t *testing.T) {
	configSetup()
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

	continuations, err := req.process(db)
	if err != nil {
		t.Fatal(err)
	}

	// didn't call extra db functions
	if db.FunctionCallCount != 1 {
		t.Fatal("did not call correct number of db functions")
	}

	// are we notifying the right people
	if len(continuations) != 1 ||
		reflect.TypeOf(continuations[0]).String() != "datahandling.toSenderClos" {
		t.Fatal("did not properly process")
	}
	response := continuations[0].(toSenderClos).msg.ServerMessage.(response)
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
