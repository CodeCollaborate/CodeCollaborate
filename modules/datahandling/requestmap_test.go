package datahandling

import (
	"encoding/json"
	"reflect"
	"testing"
)

/*
 *
 * authenticated
 *
 */

const TestSenderID = "testUser"

// Project functions
func TestProjectCreateRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "Create"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"Name\": \"Namey\"" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.projectCreateRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestProjectRenameRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "Rename"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"NewName\": \"Namey\", " +
		"\"ProjectID\": 12345" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.projectRenameRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestProjectGetPermissionConstantsRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "GetPermissionsConstants"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.projectGetPermissionConstantsRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestProjectGrantPermissionsRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "GrantPermissions"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"ProjectID\": 12345, " +
		"\"GrantUsername\": \"loganga\", " +
		"\"PermissionLevel\": 1" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.projectGrantPermissionsRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestProjectRevokePermissionsRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "RevokePermissions"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"ProjectID\": 12345, " +
		"\"RevokeUsername\": \"loganga\"" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.projectRevokePermissionsRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestProjectGetOnlineClientsRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "GetOnlineClients"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"ProjectID\": 12345" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.projectGetOnlineClientsRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestProjectLookupRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "Lookup"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{\"ProjectIds\": [12345, 38292]}")
	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.projectLookupRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestProjectGetFilesRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "GetFiles"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"ProjectID\": 12345" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.projectGetFilesRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestProjectSubscribeRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "Subscribe"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"ProjectID\": 12345" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.projectSubscribeRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestProjectUnsubscribeRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "Unsubscribe"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"ProjectID\": 12345" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.projectUnsubscribeRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestProjectDeleteRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "Project"
	req.Method = "Delete"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"ProjectID\": 12345" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.projectDeleteRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

// File functions

func TestFileCreateRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Create"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"Name\": \"Namey\", " +
		"\"RelativePath\": \"src/\", " +
		"\"ProjectID\": 12345, " +
		"\"FileBytes\": [2]" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.fileCreateRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestFileRenameRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Rename"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"NewName\": \"Namey\", " +
		"\"FileID\": 12345" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.fileRenameRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestFileMoveRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Move"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"NewPath\": \"golang/\", " +
		"\"FileID\": 12345" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.fileMoveRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestFileDeleteRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Delete"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"FileID\": 12345" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.fileDeleteRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestFileChangeRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Change"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"FileID\": 12345," +
		"\"FileVersion\": 25," +
		"\"Changes\": [\"ok\", \"k\"]" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.fileChangeRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestFilePullRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "File"
	req.Method = "Pull"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{" +
		"\"FileID\": 12345" +
		"}")

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.filePullRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

// User functions

func TestUserLookupRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "User"
	req.Method = "Lookup"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage(
		"{\"Usernames\": [\"jshap70\"]" +
			"}")
	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.userLookupRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

func TestUserProjectsRequest(t *testing.T) {
	req := *new(abstractRequest)
	req.Resource = "User"
	req.Method = "Projects"
	req.SenderID = TestSenderID
	req.SenderToken = testToken(t, TestSenderID)
	req.Data = json.RawMessage("{}")
	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.userProjectsRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}

/*
 *
 * unauthenticated
 *
 */

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

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.userRegisterRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
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

	newRequest, err := getFullRequest(&req)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(newRequest).String() != "*datahandling.userLoginRequest" {
		t.Fatalf("wrong request type, got: %s", reflect.TypeOf(newRequest))
	}
}
