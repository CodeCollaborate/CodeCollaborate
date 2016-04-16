package datahandling

import (
	"testing"
)

var testJson []byte = []byte(
		"{\"Tag\":12345, " +
		"\"Resource\":\"Project\", " +
		"\"SenderId\":\"loganga\", " +
		"\"SenderToken\":\"test\", " +
		"\"Method\":\"Lookup\", " +
		"\"Time\":1460839273, " +
		"\"Data\":{\"ProjectIds\": [{\"ProjectId\" :12345}]}}")


func TestCreateAbstractRequest(t *testing.T) {
	req, err := CreateAbstractRequest(testJson);
	if (err != nil) {
		t.Fatal(err)
	}
	if (req.Tag != 12345) {
		t.Fatal(req.Tag)
	}
	if (req.Resource != "Project") {
		t.Fatal(req.Resource)
	}
	if (req.SenderId != "loganga") {
		t.Fatal(req.SenderId)
	}
	if (req.SenderToken != "test") {
		t.Fatal(req.SenderToken)
	}
	if (req.Method != "Lookup") {
		t.Fatal(req.Method)
	}
	if (req.Time != 1460839273) {
		t.Fatal(req.Time)
	}
	if (req.Data == nil) {
		t.Fail()
	}
}


