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
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	if err != nil {
		//if err == project already exists {
		// TODO: implement a specific error for this on the mysql.go side
		//}

		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   struct{ ProjectID int64 }{-1},
		}
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data: struct {
				ProjectID int64
			}{
				projectID,
			}}
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
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	if err != nil {
		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   struct{}{}}
		not = nil // don't send anything
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data:   struct{}{}}
		not.ServerMessage = notification{
			Resource: p.Resource,
			Method:   p.Method,
			Data: struct {
				NewName string
			}{
				p.NewName,
			}}
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
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"
	res.ServerMessage = response{
		Status: unimplemented,
		Tag:    p.Tag,
		Data:   struct{}{}}
	return res, nil, nil
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
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	if err != nil {
		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   struct{}{}}
		not = nil
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data:   struct{}{}}
		not.ServerMessage = notification{
			Resource: p.Resource,
			Method:   p.Method,
			Data: struct {
				GrantUsername   string
				PermissionLevel int
			}{
				p.GrantUsername,
				p.PermissionLevel,
			}}
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
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	if err != nil {
		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   struct{}{}}
		not = nil
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data:   struct{}{}}
		not.ServerMessage = notification{
			Resource: p.Resource,
			Method:   p.Method,
			Data: struct {
				RevokeUsername string
			}{
				p.RevokeUsername,
			}}
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
	// TODO: implement on redis (and actually implement redis)
	fmt.Printf("Recieved project get online clients request from %s\n", p.SenderID)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"
	res.ServerMessage = response{
		Status: unimplemented,
		Tag:    p.Tag,
		Data:   struct{}{}}
	return res, nil, nil
}

func (p *projectGetOnlineClientsRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.Lookup
type projectLookupRequest struct {
	ProjectIDs []int64
	abstractRequest
}

// this request returns a slice of results for the projects we found, so we need the object that goes in that slice
type projectLookupResult struct {
	FileID      int64
	Name        string
	Permissions map[string](dbfs.ProjectPermission)
}

func (p projectLookupRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	/*
		We could do
			data := make([]interface{}, len(p.ProjectIDs))
		but it seems like poor practice and makes the object oriented side of my brain cry
	*/
	resultData := make([]projectLookupResult, len(p.ProjectIDs))

	var errOut error
	i := 0
	for _, id := range p.ProjectIDs {
		// TODO: see note at modules/dbfs/mysql.go:307
		name, permissions, err := dbfs.MySQLProjectLookup(id, p.SenderID)
		if err != nil {
			errOut = err
		} else {
			resultData[i] = projectLookupResult{
				FileID:      id,
				Name:        name,
				Permissions: permissions}
			i++
		}
	}
	// shrink to cut off remainder left by errors
	resultData = resultData[:i+1]

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	if errOut != nil {
		if len(resultData) == 0 {
			res.ServerMessage = response{
				Status: fail,
				Tag:    p.Tag,
				Data: struct {
					Projects []projectLookupResult
				}{
					resultData,
				}}
		} else {
			res.ServerMessage = response{
				Status: partialfail,
				Tag:    p.Tag,
				Data: struct {
					Projects []projectLookupResult
				}{
					resultData,
				}}
		}
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data: struct {
				Projects []projectLookupResult
			}{
				resultData,
			}}
	}

	//fmt.Printf("Recieved project lookup request from %s\n", p.SenderID)
	return res, nil, nil
}

func (p *projectLookupRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.GetFiles
type projectGetFilesRequest struct {
	ProjectID int64
	abstractRequest
}

type fileLookupResult struct {
	FileID       int64
	Filename     string
	Creator      string
	CreationDate time.Time
	RelativePath string
	ProjectID    int64
	Version      int64
}

func (p projectGetFilesRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	files, err := dbfs.MySQLProjectGetFiles(p.ProjectID)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	if err != nil {
		res.ServerMessage = response{
			Status: fail,
			Tag:    p.Tag,
			Data: struct {
				Files []fileLookupResult
			}{
				make([]fileLookupResult, 0),
			}}

		return res, nil, nil
	}

	resultData := make([]fileLookupResult, len(files))

	i := 0
	var errOut error
	for _, file := range files {
		version, err := dbfs.CBGetFileVersion(file.FileID)
		if err != nil {
			errOut = err
		} else {
			resultData[i] = fileLookupResult{
				FileID:       file.FileID,
				Filename:     file.Filename,
				Creator:      file.Creator,
				CreationDate: file.CreationDate,
				RelativePath: file.RelativePath,
				ProjectID:    file.ProjectID,
				Version:      version}
			i++
		}
	}
	// shrink to cut off remainder left by errors
	resultData = resultData[:i+1]

	if errOut != nil {
		if len(resultData) == 0 {
			res.ServerMessage = response{
				Status: fail,
				Tag:    p.Tag,
				Data: struct {
					Files []fileLookupResult
				}{
					resultData,
				}}
		} else {
			res.ServerMessage = response{
				Status: partialfail,
				Tag:    p.Tag,
				Data: struct {
					Files []fileLookupResult
				}{
					resultData,
				}}
		}
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data: struct {
				Files []fileLookupResult
			}{
				resultData,
			}}
	}

	return res, nil, nil
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
	// TODO: figure out how to subscribe websockets to more/less rabbit sockets
	// SERIOUS ISSUE HERE
	// NOTE: we don't have scope here to either the websocket or rabbit

	fmt.Printf("Recieved project subscribe request from %s\n", p.SenderID)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"
	res.ServerMessage = response{
		Status: unimplemented,
		Tag:    p.Tag,
		Data:   struct{}{}}
	return res, nil, nil
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
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	err := dbfs.MySQLProjectDelete(p.ProjectID, p.SenderID)
	if err != nil {
		if err == dbfs.ErrNoDbChange {
			res.ServerMessage = response{
				Status: fail,
				Tag:    p.Tag,
				Data:   struct{}{}}
		} else {
			res.ServerMessage = response{
				Status: servfail,
				Tag:    p.Tag,
				Data:   struct{}{}}
		}
		not = nil
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data:   struct{}{}}

		not.ServerMessage = notification{
			Resource: p.Resource,
			Method:   p.Method,
			Data: struct {
				DeletedProjectID int64
			}{
				p.ProjectID,
			}}
	}

	return res, not, nil
}

func (p *projectDeleteRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}
