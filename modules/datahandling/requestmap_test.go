package datahandling

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestAuthenticatedGetRequestMap(t *testing.T) {
	req := *new(abstractRequest)

	req.Resource = "Project"
	req.Method = "Lookup"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{\"ProjectIds\": [12345, 38292]}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *projectLookupRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *projectLookupRequest:

	}

}

func TestUnauthenticatedRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "User"
	req.Method = "Register"

	req.Data = json.RawMessage(
		"{\"Username\": \"loganga\", " +
			"\"FirstName\":\"Gene\", " +
			"\"LastName\":\"Logan\", " +
			"\"Email\":\"coolkid69@jithub.com\", " +
			"\"Password\":\"correct horse battery staple\"" +
			"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *userRegiserRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *userRegisterRequest:

	}
}
