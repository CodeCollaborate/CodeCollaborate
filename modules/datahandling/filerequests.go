package datahandling

import (
	"fmt"
)

var fileRequestsSetup = false

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initFileRequests() {
	if fileRequestsSetup {
		return
	}

	authenticatedRequestMap["File.Create"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(fileCreateRequest), req)
	}

	authenticatedRequestMap["File.Rename"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(fileRenameRequest), req)
	}

	authenticatedRequestMap["File.Move"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(fileMoveRequest), req)
	}

	authenticatedRequestMap["File.Delete"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(fileDeleteRequest), req)
	}

	authenticatedRequestMap["File.Change"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(fileChangeRequest), req)
	}

	authenticatedRequestMap["File.Pull"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(filePullRequest), req)
	}

	fileRequestsSetup = true
}

// File.Create
type fileCreateRequest struct {
	Name         string
	RelativePath string
	ProjectID    int64
	FileBytes    []byte
	abstractRequest
}

func (f *fileCreateRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileCreateRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file create request from %s\n", f.SenderID)
	return nil, nil, nil
}

// File.Rename
type fileRenameRequest struct {
	FileID  int64
	NewName string
	abstractRequest
}

func (f *fileRenameRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileRenameRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file rename request from %s\n", f.SenderID)
	return nil, nil, nil
}

// File.Move
type fileMoveRequest struct {
	FileID  int64
	NewPath string
	abstractRequest
}

func (f *fileMoveRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileMoveRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file move request from %s\n", f.SenderID)
	return nil, nil, nil
}

// File.Delete
type fileDeleteRequest struct {
	FileID int64
	abstractRequest
}

func (f *fileDeleteRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileDeleteRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file delete request from %s\n", f.SenderID)
	return nil, nil, nil
}

// File.Change
type fileChangeRequest struct {
	FileID      int64
	Changes     []string
	FileVersion int64
	abstractRequest
}

func (f *fileChangeRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileChangeRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file change request from %s\n", f.SenderID)
	return nil, nil, nil
}

// File.Pull
type filePullRequest struct {
	FileID int64
	abstractRequest
}

func (f *filePullRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f filePullRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file pull request from %s\n", f.SenderID)
	return nil, nil, nil
}
