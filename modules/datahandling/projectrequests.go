package datahandling

import (
	"fmt"
)

var projectRequestsSetup = false

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initProjectRequests() {
	if projectRequestsSetup {
		return
	}

	authenticatedRequestMap["Project.Create"] = func(req *abstractRequest) (request, error) {
		return commonJson(new(projectCreateRequest), req)
	}

	authenticatedRequestMap["Project.Rename"] = func(req *abstractRequest) (request, error) {
		return commonJson(new(projectRenameRequest), req)
	}

	authenticatedRequestMap["Project.GetPermissionsConstants"] = func(req *abstractRequest) (request, error) {
		return commonJson(new(projectGetPermissionConstantsRequest), req)
	}

	authenticatedRequestMap["Project.GrantPermissions"] = func(req *abstractRequest) (request, error) {
		return commonJson(new(projectGrantPermissionsRequest), req)
	}

	authenticatedRequestMap["Project.RevokePermissions"] = func(req *abstractRequest) (request, error) {
		return commonJson(new(projectRevokePermissionsRequest), req)
	}

	authenticatedRequestMap["Project.GetOnlineClients"] = func(req *abstractRequest) (request, error) {
		return commonJson(new(projectGetOnlineClientsRequest), req)
	}

	authenticatedRequestMap["Project.Lookup"] = func(req *abstractRequest) (request, error) {
		return commonJson(new(projectLookupRequest), req)
	}

	authenticatedRequestMap["Project.GetFiles"] = func(req *abstractRequest) (request, error) {
		return commonJson(new(projectGetFilesRequest), req)
	}

	authenticatedRequestMap["Project.Subscribe"] = func(req *abstractRequest) (request, error) {
		return commonJson(new(projectSubscribeRequest), req)
	}

	authenticatedRequestMap["Project.Delete"] = func(req *abstractRequest) (request, error) {
		return commonJson(new(projectDeleteRequest), req)
	}

	projectRequestsSetup = true
}

// Project.Create
type projectCreateRequest struct {
	Name string
	abstractRequest
}

func (p *projectCreateRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

func (p projectCreateRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project create request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.Rename
type projectRenameRequest struct {
	ProjectID int64
	NewName   string
	abstractRequest
}

func (p *projectRenameRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
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

func (p *projectGetPermissionConstantsRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

func (p projectGetPermissionConstantsRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project get permissions constants request from %s\n", p.SenderID)
	return nil, nil, nil
}

// Project.GrantPermissions
type projectGrantPermissionsRequest struct {
	ProjectID       int64
	GrantUsername   string
	PermissionLevel int
	abstractRequest
}

func (p projectGrantPermissionsRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project grant permissions request from %s\n", p.SenderID)
	return nil, nil, nil
}

func (p *projectGrantPermissionsRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.RevokePermissions
type projectRevokePermissionsRequest struct {
	ProjectID      int64
	RevokeUsername string
	abstractRequest
}

func (p projectRevokePermissionsRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project revoke permissions request from %s\n", p.SenderID)
	return nil, nil, nil
}

func (p *projectRevokePermissionsRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.GetOnlineClients
type projectGetOnlineClientsRequest struct {
	ProjectID int64
	abstractRequest
}

func (p projectGetOnlineClientsRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project get online clients request from %s\n", p.SenderID)
	return nil, nil, nil
}

func (p *projectGetOnlineClientsRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
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

func (p *projectLookupRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.GetFiles
type projectGetFilesRequest struct {
	ProjectID int64
	abstractRequest
}

func (p projectGetFilesRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved get project files request from %s\n", p.SenderID)
	return nil, nil, nil
}

func (p *projectGetFilesRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.Subscribe
type projectSubscribeRequest struct {
	ProjectID int64
	abstractRequest
}

func (p projectSubscribeRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project subscribe request from %s\n", p.SenderID)
	return nil, nil, nil
}

func (p *projectSubscribeRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.Delete
type projectDeleteRequest struct {
	ProjectID int64
	abstractRequest
}

func (p projectDeleteRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved project delete request from %s\n", p.SenderID)
	return nil, nil, nil
}

func (p *projectDeleteRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}




