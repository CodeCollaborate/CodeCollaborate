package datahandling

import (
	"encoding/json"
	"fmt"
)

var fileRequestsSetup = false

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initFileRequests() {
	if fileRequestsSetup {
		return
	}

	authenticatedRequestMap["FileCreate"] = func(req abstractRequest) (request, error) {
		p := new(fileCreateRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["FileRename"] = func(req abstractRequest) (request, error) {
		p := new(fileRenameRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["FileMove"] = func(req abstractRequest) (request, error) {
		p := new(fileMoveRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["FileDelete"] = func(req abstractRequest) (request, error) {
		p := new(fileDeleteRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["FileChange"] = func(req abstractRequest) (request, error) {
		p := new(fileChangeRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["FilePull"] = func(req abstractRequest) (request, error) {
		p := new(filePullRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
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

func (p filePullRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved file pull request from %s\n", p.SenderID)
	return nil, nil, nil
}
