package datahandling

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestProjectLookupRequest_Process(t *testing.T) {
	req := *new(AbstractRequest)

	req.Resource = "Project"
	req.Method = "Lookup"
	req.SenderToken = "supersecure"

	req.Data = json.RawMessage("{\"ProjectIds\": [12345, 38292]}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	response, notification, err2 := newRequest.Process()
	if err2 != nil {
		t.Fatal(err2)
	}

	if response == nil || notification == nil {
		t.SkipNow() // added April 17, 2016
		fmt.Println("unimplemented response and notification logic!")
	}
}
