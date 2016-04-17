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
	authenticatedRequestMap["ProjectCreateRequest"] = func(req abstractRequest) (request, error) {
		p := new(projectCreateRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["ProjectRename"] = func(req abstractRequest) (request, error) {
		p := new(projectRenameRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["ProjectGetPermissionsConstants"] = func(req abstractRequest) (request, error) {
		p := new(projectGetPermissionConstantsRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["ProjectGrantPermissions"] = func(req abstractRequest) (request, error) {
		p := new(projectGrantPermissionsRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["ProjectRevokePermissions"] = func(req abstractRequest) (request, error) {
		p := new(projectRevokePermissionsRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["ProjectGetOnlineClients"] = func(req abstractRequest) (request, error) {
		p := new(projectGetOnlineClientsRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["ProjectLookup"] = func(req abstractRequest) (request, error) {
		p := new(projectLookupRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["ProjectGetFiles"] = func(req abstractRequest) (request, error) {
		p := new(projectGetFilesRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["ProjectSubscribe"] = func(req abstractRequest) (request, error) {
		p := new(projectSubscribeRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	authenticatedRequestMap["ProjectDelete"] = func(req abstractRequest) (request, error) {
		p := new(projectDeleteRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	projectRequestsSetup = true
}

// Project.Create
type projectCreateRequest struct {
	Name string
	abstractRequest
}

func (p projectCreateRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project create request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.Rename
type projectRenameRequest struct {
	ProjectID string
	NewName   string
	abstractRequest
}

func (p projectRenameRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project rename request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.GetPermissionConstants
type projectGetPermissionConstantsRequest struct {
	abstractRequest
}

func (p projectGetPermissionConstantsRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project get permissions constants request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.GrantPermissions
type projectGrantPermissionsRequest struct {
	ProjectID       string
	GrantUsername   string
	PermissionLevel int
	abstractRequest
}

func (p projectGrantPermissionsRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project grant permissions request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.RevokePermissions
type projectRevokePermissionsRequest struct {
	ProjectID      string
	RevokeUsername string
	abstractRequest
}

func (p projectRevokePermissionsRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project revoke permissions request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.GetOnlineClients
type projectGetOnlineClientsRequest struct {
	ProjectID string
	abstractRequest
}

func (p projectGetOnlineClientsRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project get online clients request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.Lookup
type projectLookupRequest struct {
	ProjectIDs []int64
	abstractRequest
}

func (p projectLookupRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project lookup request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.GetFiles
type projectGetFilesRequest struct {
	ProjectID string
	abstractRequest
}

func (p projectGetFilesRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved get project files request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.Subscribe
type projectSubscribeRequest struct {
	ProjectID string
	abstractRequest
}

func (p projectSubscribeRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project subscribe request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.Delete
type projectDeleteRequest struct {
	ProjectID string
	abstractRequest
}

func (p projectDeleteRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project delete request from %s\n", p.SenderID)
	return nil, nil, nil
}
