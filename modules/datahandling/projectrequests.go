package datahandling

import (
	"fmt"
	"time"

	"github.com/CodeCollaborate/Server/modules/dbfs"
)

var projectRequestsSetup = false

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initProjectRequests() {
	if projectRequestsSetup {
		return
	}

	authenticatedRequestMap["Project.Create"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(projectCreateRequest), req)
	}

	authenticatedRequestMap["Project.Rename"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(projectRenameRequest), req)
	}

	authenticatedRequestMap["Project.GetPermissionsConstants"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(projectGetPermissionConstantsRequest), req)
	}

	authenticatedRequestMap["Project.GrantPermissions"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(projectGrantPermissionsRequest), req)
	}

	authenticatedRequestMap["Project.RevokePermissions"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(projectRevokePermissionsRequest), req)
	}

	authenticatedRequestMap["Project.GetOnlineClients"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(projectGetOnlineClientsRequest), req)
	}

	authenticatedRequestMap["Project.Lookup"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(projectLookupRequest), req)
	}

	authenticatedRequestMap["Project.GetFiles"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(projectGetFilesRequest), req)
	}

	authenticatedRequestMap["Project.Subscribe"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(projectSubscribeRequest), req)
	}

	authenticatedRequestMap["Project.Delete"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(projectDeleteRequest), req)
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

func (p projectCreateRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	projectID, err := dbfs.MySQLProjectCreate(p.SenderID, p.Name)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now()
	res.Type = "Responce"

	if err != nil {
		//if err == project already exists {
		// TODO: implement a specific error for this on the mysql.go side
		//}

		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   {"ProjectID": -1}}
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data:   {"ProjectID": projectID}}
	}

	return res, nil, nil
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

func (p projectRenameRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {

	// TODO: check if permission high enough on project

	err := dbfs.MySQLProjectRename(p.ProjectID, p.NewName)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	if err != nil {
		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   {}}
		not = nil // don't send anything
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data:   {}}
		not.ServerMessage = notification{
			Resource: p.Resource,
			Method:   p.Method,
			Data:     {"NewName": p.NewName}}
	}

	return res, not, nil
}

// Project.GetPermissionConstants
type projectGetPermissionConstantsRequest struct {
	abstractRequest
}

func (p *projectGetPermissionConstantsRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

func (p projectGetPermissionConstantsRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	// TODO: figure out how we want to do this on the db
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

func (p projectGrantPermissionsRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	// TODO: check if permission high enough on project

	err := dbfs.MySQLProjectGrantPermission(p.ProjectID, p.GrantUsername, p.PermissionLevel, p.SenderID)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	if err != nil {
		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   {}}
		not = nil
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data:   {}}
		not.ServerMessage = notification{
			Resource: p.Resource,
			Method:   p.Method,
			Data: {
				"GrantUsername":   p.GrantUsername,
				"PermissionLevel": p.PermissionLevel}}
	}

	return res, not, nil
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

func (p projectRevokePermissionsRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	// TODO: check if permission high enough on project
	err := dbfs.MySQLProjectRevokePermission(p.ProjectID, p.RevokeUsername, p.SenderID)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	if err != nil {
		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   {}}
		not = nil
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data:   {}}
		not.ServerMessage = notification{
			Resource: p.Resource,
			Method:   p.Method,
			Data:     {"RevokeUsername": p.RevokeUsername}}
	}

	return res, not, nil
}

func (p *projectRevokePermissionsRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.GetOnlineClients
type projectGetOnlineClientsRequest struct {
	ProjectID int64
	abstractRequest
}

func (p projectGetOnlineClientsRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	// TODO: add on redis
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

func (p projectLookupRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
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

func (p projectGetFilesRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
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

func (p projectSubscribeRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
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

func (p projectDeleteRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	// TODO
	fmt.Printf("Recieved project delete request from %s\n", p.SenderID)
	return nil, nil, nil
}

func (p *projectDeleteRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}
