package datahandling

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestProjectLookupRequest_Process(t *testing.T) {
	req := *new(abstractRequest)

	req.SenderID = "loganga"
	req.Resource = "Project"
	req.Method = "Lookup"
	req.SenderToken = "supersecure"

	req.Data = json.RawMessage("{\"ProjectIds\": [12345, 38292]}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	response, notification, err2 := newRequest.process()
	if err2 != nil {
		t.Fatal(err2)
	}

	if response == nil || notification == nil {
		// TODO: check if process worked
		fmt.Println("unimplemented response and notification logic!")
	}
}
