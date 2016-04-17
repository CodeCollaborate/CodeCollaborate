package datahandling

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestGetRequestMap(t *testing.T) {
	req := *new(AbstractRequest)

	req.Resource = "Project"
	req.Method = "Lookup"

	req.Data = json.RawMessage("{\"ProjectIds\": [12345, 38292]}")

	newRequest, err := GetRequestMap(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *projectLookupRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *projectLookupRequest:

	}

}
