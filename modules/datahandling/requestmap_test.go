package datahandling

import (
	"testing"
	"encoding/json"
)

func TestGetRequestMap(t *testing.T) {
	req := *new(AbstractRequest)

	req.Resource = "Project"
	req.Method = "Lookup"

	req.Data = json.RawMessage("{\"ProjectIds\": [{\"ProjectId\" :12345}]}")

	newRequest, err := GetRequestMap(req)
	if (err != nil) {
		t.Fatal(err)
	}
	err2 := newRequest.Process()
	if (err2 != nil) {
		t.Fatal("Dicks")
	}
}