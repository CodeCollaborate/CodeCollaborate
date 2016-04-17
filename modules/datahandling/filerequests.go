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

	requestMap["FileCreate"] = func(req AbstractRequest) (Request, error) {
		p := new(fileCreateRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["FileRename"] = func(req AbstractRequest) (Request, error) {
		p := new(fileRenameRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["FileMove"] = func(req AbstractRequest) (Request, error) {
		p := new(fileMoveRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["FileDelete"] = func(req AbstractRequest) (Request, error) {
		p := new(fileDeleteRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["FileChange"] = func(req AbstractRequest) (Request, error) {
		p := new(fileChangeRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["FilePull"] = func(req AbstractRequest) (Request, error) {
		p := new(filePullRequest)
		p.AbstractRequest = req
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
	AbstractRequest
}

func (p fileCreateRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Printf("Recieved file create request from %s\n", p.SenderID)
	return nil, nil, nil
}

// File.Rename
type fileRenameRequest struct {
	FileID  string
	NewName string
	AbstractRequest
}

func (p fileRenameRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Printf("Recieved file rename request from %s\n", p.SenderID)
	return nil, nil, nil
}

// File.Move
type fileMoveRequest struct {
	FileID  string
	NewPath string
	AbstractRequest
}

func (p fileMoveRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Printf("Recieved file move request from %s\n", p.SenderID)
	return nil, nil, nil
}

// File.Delete
type fileDeleteRequest struct {
	FileID string
	AbstractRequest
}

func (p fileDeleteRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Printf("Recieved file delete request from %s\n", p.SenderID)
	return nil, nil, nil
}

// File.Change
type fileChangeRequest struct {
	FileID  string
	Changes []string
	AbstractRequest
}

func (p fileChangeRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Printf("Recieved file change request from %s\n", p.SenderID)
	return nil, nil, nil
}

// File.Pull
type filePullRequest struct {
	FileID string
	AbstractRequest
}

func (p filePullRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Printf("Recieved file pull request from %s\n", p.SenderID)
	return nil, nil, nil
}
