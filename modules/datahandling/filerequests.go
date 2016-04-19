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

	authenticatedRequestMap["File.Create"] = func(req abstractRequest) (request, error) {
		return commonJson(new(fileCreateRequest), req)
	}

	authenticatedRequestMap["File.Rename"] = func(req abstractRequest) (request, error) {
		return commonJson(new(fileRenameRequest), req)
	}

	authenticatedRequestMap["File.Move"] = func(req abstractRequest) (request, error) {
		return commonJson(new(fileMoveRequest), req)
	}

	authenticatedRequestMap["File.Delete"] = func(req abstractRequest) (request, error) {
		return commonJson(new(fileDeleteRequest), req)
	}

	authenticatedRequestMap["File.Change"] = func(req abstractRequest) (request, error) {
		return commonJson(new(fileChangeRequest), req)
	}

	authenticatedRequestMap["File.Pull"] = func(req abstractRequest) (request, error) {
		return commonJson(new(filePullRequest), req)
	}

	fileRequestsSetup = true
}

// File.Create
type fileCreateRequest struct {
	Name         string
	RelativePath string
	ProjectID    string
	FileBytes    []byte
	abstractRequest
}

func (f *fileCreateRequest)setAbstractRequest(req abstractRequest) {
	f.abstractRequest = req
}

func (p fileCreateRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file create request from %s\n", p.SenderID)
	return nil, nil, nil
}

// File.Rename
type fileRenameRequest struct {
	FileID  string
	NewName string
	abstractRequest
}

func (f *fileRenameRequest)setAbstractRequest(req abstractRequest) {
	f.abstractRequest = req
}

func (p fileRenameRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file rename request from %s\n", p.SenderID)
	return nil, nil, nil
}

// File.Move
type fileMoveRequest struct {
	FileID  string
	NewPath string
	abstractRequest
}

func (f *fileMoveRequest)setAbstractRequest(req abstractRequest) {
	f.abstractRequest = req
}

func (p fileMoveRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file move request from %s\n", p.SenderID)
	return nil, nil, nil
}

// File.Delete
type fileDeleteRequest struct {
	FileID string
	abstractRequest
}

func (f *fileDeleteRequest)setAbstractRequest(req abstractRequest) {
	f.abstractRequest = req
}

func (p fileDeleteRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file delete request from %s\n", p.SenderID)
	return nil, nil, nil
}

// File.Change
type fileChangeRequest struct {
	FileID  string
	Changes []string
	abstractRequest
}

func (f *fileChangeRequest)setAbstractRequest(req abstractRequest) {
	f.abstractRequest = req
}

func (p fileChangeRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file change request from %s\n", p.SenderID)
	return nil, nil, nil
}

// File.Pull
type filePullRequest struct {
	FileID string
	abstractRequest
}

func (f *filePullRequest)setAbstractRequest(req abstractRequest) {
	f.abstractRequest = req
}

func (p filePullRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file pull request from %s\n", p.SenderID)
	return nil, nil, nil
}
