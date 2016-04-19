package datahandling

import (
	"encoding/json"
	"reflect"
	"testing"
)

// authenticated

// Project functions
func TestProjectCreateRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "Create"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
		"\"Name\": \"Namey\"" +
		"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *projectCreateRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *projectCreateRequest:

	}
}

func TestProjectRenameRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "Rename"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"NewName\": \"Namey\", " +
	"\"ProjectID\": 12345" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *projectRenameRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *projectRenameRequest:

	}
}

func TestProjectGetPermissionConstantsRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "GetPermissionsConstants"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *projectGetPermissionConstantsRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *projectGetPermissionConstantsRequest:

	}
}

func TestProjectGrantPermissionsRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "GrantPermissions"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"ProjectID\": 12345, " +
	"\"GrantUsername\": \"loganga\", " +
	"\"PermissionLevel\": 1" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *projectGrantPermissionsRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *projectGrantPermissionsRequest:

	}
}

func TestProjectRevokePermissionsRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "RevokePermissions"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"ProjectID\": 12345, " +
	"\"RevokeUsername\": \"loganga\"" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *projectRevokePermissionsRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *projectRevokePermissionsRequest:

	}
}

func TestProjectGetOnlineClientsRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "GetOnlineClients"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"ProjectID\": 12345" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *projectGetOnlineClientsRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *projectGetOnlineClientsRequest:

	}
}

func TestProjectLookupRequest(t *testing.T) {
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

func TestProjectGetFilesRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "GetFiles"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"ProjectID\": 12345" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *projectGetFilesRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *projectGetFilesRequest:

	}
}

func TestProjectSubscribeRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "Subscribe"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"ProjectID\": 12345" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *projectSubscribeRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *projectSubscribeRequest:

	}
}

func TestProjectDeleteRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "Delete"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"ProjectID\": 12345" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *projectDeleteRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *projectDeleteRequest:

	}
}

// File functions

func TestFileCreateRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Create"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"Name\": \"Namey\", " +
	"\"RelativePath\": \"src/\", " +
	"\"ProjectID\": 12345, " +
	"\"FileBytes\": [2]" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *fileCreateRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *fileCreateRequest:

	}
}

func TestFileRenameRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Rename"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"NewName\": \"Namey\", " +
	"\"FileID\": 12345" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *fileRenameRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *fileRenameRequest:

	}
}

func TestFileMoveRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Move"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"NewPath\": \"golang/\", " +
	"\"FileID\": 12345" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *fileMoveRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *fileMoveRequest:

	}
}

func TestFileDeleteRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Delete"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"FileID\": 12345" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *fileDeleteRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *fileDeleteRequest:

	}
}

func TestFileChangeRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Change"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"FileID\": 12345," +
	"\"Changes\": [\"ok\", \"k\"]" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *fileChangeRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *fileChangeRequest:

	}
}

func TestFilePullRequest(t *testing.T)	{
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Pull"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{" +
	"\"FileID\": 12345" +
	"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *filePullRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *filePullRequest:

	}
}

// User functions

func TestUserLookupRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "User"
	req.Method = "Lookup"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage(
		"{\"Usernames\": [\"jshap70\"]" +
		"}")
	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *userRegiserRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *userLookupRequest:

	}
}

func TestUserProjectsRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "User"
	req.Method = "Projects"
	req.SenderToken = "supersecure"
	req.Data = json.RawMessage("{}")
	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *userRegiserRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *userProjectsRequest:

	}
}

// unauthenticated
func TestUserRegisterRequest(t *testing.T) {
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
		t.Fatalf("newRequest is the wrong type, expected: *userRegisterRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *userRegisterRequest:

	}
}

func TestUserLoginRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "User"
	req.Method = "Login"
	req.Data = json.RawMessage(
		"{\"Username\": \"loganga\", " +
		"\"Password\":\"correct horse battery staple\"" +
		"}")

	newRequest, err := getFullRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	switch newRequest.(type) {
	default:
		t.Fatalf("newRequest is the wrong type, expected: *userLoginRequest, actual: %s\n", reflect.TypeOf(newRequest))
	case *userLoginRequest:

	}
}
