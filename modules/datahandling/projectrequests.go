package datahandling

import (
	"encoding/json"
	"fmt"
)

var projectRequestsSetup = false

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initProjectRequests() {
	if projectRequestsSetup {
		return
	}

	// This isn't an ideal way of populating the map since there is so much duplicate code, but extracting this
	// requires some kind of restructuring of the request hierarchy, which can come later because we can't come up
	// with an obvious way to do it.
	requestMap["ProjectCreateRequest"] = func(req AbstractRequest) (Request, error) {
		p := new(projectCreateRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["ProjectRename"] = func(req AbstractRequest) (Request, error) {
		p := new(projectRenameRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["ProjectGetPermissionsConstants"] = func(req AbstractRequest) (Request, error) {
		p := new(projectGetPermissionConstantsRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["ProjectGrantPermissions"] = func(req AbstractRequest) (Request, error) {
		p := new(projectGrantPermissionsRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["ProjectRevokePermissions"] = func(req AbstractRequest) (Request, error) {
		p := new(projectRevokePermissionsRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["ProjectGetOnlineClients"] = func(req AbstractRequest) (Request, error) {
		p := new(projectGetOnlineClientsRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["ProjectLookup"] = func(req AbstractRequest) (Request, error) {
		p := new(projectLookupRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["ProjectGetFiles"] = func(req AbstractRequest) (Request, error) {
		p := new(projectGetFilesRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["ProjectSubscribe"] = func(req AbstractRequest) (Request, error) {
		p := new(projectSubscribeRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["ProjectDelete"] = func(req AbstractRequest) (Request, error) {
		p := new(projectDeleteRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	projectRequestsSetup = true
}

// Project.Create
type projectCreateRequest struct {
	Name string
	AbstractRequest
}

func (p projectCreateRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project create request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.Rename
type projectRenameRequest struct {
	ProjectID string
	NewName   string
	AbstractRequest
}

func (p projectRenameRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project rename request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.GetPermissionConstants
type projectGetPermissionConstantsRequest struct {
	AbstractRequest
}

func (p projectGetPermissionConstantsRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project get permissions constants request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.GrantPermissions
type projectGrantPermissionsRequest struct {
	ProjectID       string
	GrantUsername   string
	PermissionLevel int
	AbstractRequest
}

func (p projectGrantPermissionsRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project grant permissions request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.RevokePermissions
type projectRevokePermissionsRequest struct {
	ProjectID      string
	RevokeUsername string
	AbstractRequest
}

func (p projectRevokePermissionsRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project revoke permissions request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.GetOnlineClients
type projectGetOnlineClientsRequest struct {
	ProjectID string
	AbstractRequest
}

func (p projectGetOnlineClientsRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project get online clients request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.Lookup
type projectLookupRequest struct {
	ProjectIDs []int64
	AbstractRequest
}

func (p projectLookupRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project lookup request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.GetFiles
type projectGetFilesRequest struct {
	ProjectID string
	AbstractRequest
}

func (p projectGetFilesRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved get project files request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.Subscribe
type projectSubscribeRequest struct {
	ProjectID string
	AbstractRequest
}

func (p projectSubscribeRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project subscribe request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.Delete
type projectDeleteRequest struct {
	ProjectID string
	AbstractRequest
}

func (p projectDeleteRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project delete request from %s\n", p.SenderID)
	return nil, nil, nil
}
