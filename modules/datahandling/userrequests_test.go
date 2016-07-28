package datahandling

import (
	"reflect"
	"testing"

	"github.com/CodeCollaborate/Server/modules/dbfs"
)

func TestUserRegisterRequest_Process(t *testing.T) {
	configSetup()
	req := *new(userRegisterRequest)

	req.SenderID = "loganga"
	req.Resource = "User"
	req.Method = "Register"
	req.SenderToken = "supersecure"
	req.Username = "loganga"
	req.FirstName = "Gene"
	req.LastName = "Logan"
	req.Email = "loganga@codecollaborate.com"
	req.Password = "correct horse battery staple"

	db := dbfs.NewDBMock()

	continuations, err := req.process(db)

	if err != nil {
		t.Fatal(err)
	}

	if len(continuations) != 1 ||
		reflect.TypeOf(continuations[0]).String() != "func(datahandling.DataHandler) error" {
		t.Fatal("did not properly process")
	}

	if db.FunctionCallCount != 1 {
		t.Fatal("did not call correct number of db functions")
	}

	if _, ok := db.Users["loganga"]; !ok {
		t.Fatal("did not correctly call db function")
	}
}

// this is commented out because userLoginRequest.process is unimplemented
//
//func TestUserLoginRequest_Process(t *testing.T) {
//	configSetup()
//	req := *new(userRegisterRequest)
//
//	req.SenderID = "loganga"
//	req.Resource = "User"
//	req.Method = "Login"
//	req.SenderToken = "supersecure"
//	req.Username  = "loganga"
//	req.Password  = "correct horse battery staple"
//
//	db := dbfs.NewDBMock()
//
//	meta := dbfs.UserMeta{
//		FirstName:"Gene",
//		LastName:"Logan",
//		Email:"loganga@codecollaborate.com",
//		Password:"correct horse battery staple",
//		Username:"loganga",
//	}
//	db.Users["loganga"] = meta
//
//	continuations, err := req.process(db)
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if len(continuations) != 1 {
//		t.Fatal("did not properly process")
//	}
//
//	if db.FunctionCallCount != 1 {
//		t.Fatal("did not call correct number of db functions")
//	}
//
//	if _, ok := db.Users["loganga"]; !ok {
//		t.Fatal("did not correctly call db function")
//	}
//}

func TestUserLookupRequest_Process(t *testing.T) {
	configSetup()
	req := *new(userRegisterRequest)

	req.SenderID = "loganga"
	req.Resource = "User"
	req.Method = "Login"
	req.SenderToken = "supersecure"
	req.Username = "loganga"
	req.Password = "correct horse battery staple"

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

	if len(continuations) != 1 ||
		reflect.TypeOf(continuations[0]).String() != "func(datahandling.DataHandler) error" {
		t.Fatal("did not properly process")
	}

	if db.FunctionCallCount != 1 {
		t.Fatal("did not call correct number of db functions")
	}

	if _, ok := db.Users["loganga"]; !ok {
		t.Fatal("did not correctly call db function")
	}
}
