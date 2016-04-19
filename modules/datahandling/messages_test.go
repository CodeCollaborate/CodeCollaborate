package datahandling

import (
	"testing"
)

var testJSON = []byte(
	"{\"Tag\":12345, " +
		"\"Resource\":\"Project\", " +
		"\"Method\":\"Lookup\", " +
		"\"SenderID\":\"loganga\", " +
		"\"SenderToken\":\"test\", " +
		"\"Timestamp\":1460839273, " +
		"\"Data\":{\"ProjectIds\": [{\"ProjectId\" :12345}]}}")

func TestCreateAbstractRequest(t *testing.T) {
	req, err := createAbstractRequest(testJSON)
	if err != nil {
		t.Fatal(err)
	}
	if req.Tag != 12345 {
		t.Fatal(req.Tag)
	}
	if req.Resource != "Project" {
		t.Fatal(req.Resource)
	}
	if req.Method != "Lookup" {
		t.Fatal(req.Method)
	}
	if req.SenderID != "loganga" {
		t.Fatal(req.SenderID)
	}
	if req.SenderToken != "test" {
		t.Fatal(req.SenderToken)
	}
	if req.Timestamp != 1460839273 {
		t.Fatal(req.Timestamp)
	}
	if req.Data == nil {
		t.Fail()
	}
}
