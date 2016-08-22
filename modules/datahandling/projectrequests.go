package datahandling

import (
	"fmt"
	"strconv"
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

	authenticatedRequestMap["Project.Unsubscribe"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(projectUnsubscribeRequest), req)
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

func (p projectCreateRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	projectID, err := db.MySQLProjectCreate(p.SenderID, p.Name)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Response"

	if err != nil {
		//if err == project already exists {
		// TODO: implement a specific error for this on the mysql.go side
		//}

		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   struct{ ProjectID int64 }{ProjectID: -1},
		}
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data: struct {
				ProjectID int64
			}{
				ProjectID: projectID,
			}}
	}

	return []dhClosure{toSenderClosure{msg: res}}, nil
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

func (p projectRenameRequest) process(db dbfs.DBFS) ([]dhClosure, error) {

	// TODO: check if permission high enough on project

	err := db.MySQLProjectRename(p.ProjectID, p.NewName)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Response"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"
	not.RoutingKey = strconv.FormatInt(p.ProjectID, 10)

	if err != nil {
		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   struct{}{}}

		return []dhClosure{toSenderClosure{msg: res}}, err
	}
	res.ServerMessage = response{
		Status: success,
		Tag:    p.Tag,
		Data:   struct{}{}}
	not.ServerMessage = notification{
		Resource:   p.Resource,
		Method:     p.Method,
		ResourceID: p.ProjectID,
		Data: struct {
			NewName string
		}{
			NewName: p.NewName,
		}}

	return []dhClosure{toSenderClosure{msg: res}, toRabbitChannelClosure{msg: not}}, nil
}

// Project.GetPermissionConstants
type projectGetPermissionConstantsRequest struct {
	abstractRequest
}

func (p *projectGetPermissionConstantsRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

func (p projectGetPermissionConstantsRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	// TODO (non-immediate/required): figure out how we want to do projectGetPermissionConstantsRequest
	fmt.Printf("Recieved project get permissions constants request from %s\n", p.SenderID)
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Response"
	res.ServerMessage = response{
		Status: unimplemented,
		Tag:    p.Tag,
		Data:   struct{}{}}
	return []dhClosure{toSenderClosure{msg: res}}, nil
}

// Project.GrantPermissions
type projectGrantPermissionsRequest struct {
	ProjectID       int64
	GrantUsername   string
	PermissionLevel int
	abstractRequest
}

func (p projectGrantPermissionsRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	// TODO: check if permission high enough on project

	err := db.MySQLProjectGrantPermission(p.ProjectID, p.GrantUsername, p.PermissionLevel, p.SenderID)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Response"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"
	not.RoutingKey = strconv.FormatInt(p.ProjectID, 10)

	if err != nil {
		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   struct{}{}}

		return []dhClosure{toSenderClosure{msg: res}}, err
	}

	res.ServerMessage = response{
		Status: success,
		Tag:    p.Tag,
		Data:   struct{}{}}
	not.ServerMessage = notification{
		Resource:   p.Resource,
		Method:     p.Method,
		ResourceID: p.ProjectID,
		Data: struct {
			GrantUsername   string
			PermissionLevel int
		}{
			GrantUsername:   p.GrantUsername,
			PermissionLevel: p.PermissionLevel,
		}}

	return []dhClosure{toSenderClosure{msg: res}, toRabbitChannelClosure{msg: not}}, nil
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

func (p projectRevokePermissionsRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	// TODO: check if permission high enough on project
	err := db.MySQLProjectRevokePermission(p.ProjectID, p.RevokeUsername, p.SenderID)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Response"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"
	not.RoutingKey = strconv.FormatInt(p.ProjectID, 10)

	if err != nil {
		res.ServerMessage = response{
			Status: servfail,
			Tag:    p.Tag,
			Data:   struct{}{}}

		return []dhClosure{toSenderClosure{msg: res}}, err
	}

	res.ServerMessage = response{
		Status: success,
		Tag:    p.Tag,
		Data:   struct{}{}}
	not.ServerMessage = notification{
		Resource:   p.Resource,
		Method:     p.Method,
		ResourceID: p.ProjectID,
		Data: struct {
			RevokeUsername string
		}{
			RevokeUsername: p.RevokeUsername,
		}}

	return []dhClosure{toSenderClosure{msg: res}, toRabbitChannelClosure{msg: not}}, nil
}

func (p *projectRevokePermissionsRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.GetOnlineClients
type projectGetOnlineClientsRequest struct {
	ProjectID int64
	abstractRequest
}

func (p projectGetOnlineClientsRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	// TODO: implement on redis (and actually implement redis)
	fmt.Printf("Recieved project get online clients request from %s\n", p.SenderID)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Response"
	res.ServerMessage = response{
		Status: unimplemented,
		Tag:    p.Tag,
		Data:   struct{}{}}
	return []dhClosure{toSenderClosure{msg: res}}, nil
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
	ProjectID   int64
	Name        string
	Permissions map[string](dbfs.ProjectPermission)
}

func (p projectLookupRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
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
		name, permissions, err := db.MySQLProjectLookup(id, p.SenderID)
		if err != nil {
			errOut = err
		} else {
			resultData[i] = projectLookupResult{
				ProjectID:   id,
				Name:        name,
				Permissions: permissions}
			i++
		}
	}
	// shrink to cut off remainder left by errors
	resultData = resultData[:i]

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Response"

	if errOut != nil {
		if len(resultData) == 0 {
			res.ServerMessage = response{
				Status: fail,
				Tag:    p.Tag,
				Data: struct {
					Projects []projectLookupResult
				}{
					Projects: resultData,
				}}
		} else {
			res.ServerMessage = response{
				Status: partialfail,
				Tag:    p.Tag,
				Data: struct {
					Projects []projectLookupResult
				}{
					Projects: resultData,
				}}
		}
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data: struct {
				Projects []projectLookupResult
			}{
				Projects: resultData,
			}}
	}

	//fmt.Printf("Received project lookup request from %s\n", p.SenderID)
	return []dhClosure{toSenderClosure{msg: res}}, nil
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

func (p projectGetFilesRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	files, err := db.MySQLProjectGetFiles(p.ProjectID)

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Response"

	if err != nil {
		res.ServerMessage = response{
			Status: fail,
			Tag:    p.Tag,
			Data: struct {
				Files []fileLookupResult
			}{
				Files: make([]fileLookupResult, 0),
			}}

		return []dhClosure{toSenderClosure{msg: res}}, nil
	}

	resultData := make([]fileLookupResult, len(files))

	i := 0
	var errOut error
	for _, file := range files {
		version, err := db.CBGetFileVersion(file.FileID)
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
	resultData = resultData[:i]

	if errOut != nil {
		if len(resultData) == 0 {
			res.ServerMessage = response{
				Status: fail,
				Tag:    p.Tag,
				Data: struct {
					Files []fileLookupResult
				}{
					Files: resultData,
				}}
		} else {
			res.ServerMessage = response{
				Status: partialfail,
				Tag:    p.Tag,
				Data: struct {
					Files []fileLookupResult
				}{
					Files: resultData,
				}}
		}
	} else {
		res.ServerMessage = response{
			Status: success,
			Tag:    p.Tag,
			Data: struct {
				Files []fileLookupResult
			}{
				Files: resultData,
			}}
	}

	return []dhClosure{toSenderClosure{msg: res}}, nil
}

func (p *projectGetFilesRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.Subscribe
type projectSubscribeRequest struct {
	ProjectID int64
	abstractRequest
}

func (p projectSubscribeRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	subscribeClosure := rabbitChannelSubscribeClosure{
		key: strconv.FormatInt(p.ProjectID, 10),
		tag: p.Tag,
	}
	return []dhClosure{subscribeClosure}, nil
}

func (p *projectSubscribeRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.Unsubscribe
type projectUnsubscribeRequest struct {
	ProjectID int64
	abstractRequest
}

func (p projectUnsubscribeRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	unsubscribeClosure := rabbitChannelUnsubscribeClosure{
		key: strconv.FormatInt(p.ProjectID, 10),
		tag: p.Tag,
	}
	return []dhClosure{unsubscribeClosure}, nil
}

func (p *projectUnsubscribeRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

// Project.Delete
type projectDeleteRequest struct {
	ProjectID int64
	abstractRequest
}

func (p projectDeleteRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Response"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"
	not.RoutingKey = strconv.FormatInt(p.ProjectID, 10)

	err := db.MySQLProjectDelete(p.ProjectID, p.SenderID)
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

		return []dhClosure{toSenderClosure{msg: res}}, err
	}
	res.ServerMessage = response{
		Status: success,
		Tag:    p.Tag,
		Data:   struct{}{}}

	not.ServerMessage = notification{
		Resource:   p.Resource,
		Method:     p.Method,
		ResourceID: p.ProjectID,
		Data:       struct{}{}}

	return []dhClosure{toSenderClosure{msg: res}, toRabbitChannelClosure{msg: not}}, nil
}

func (p *projectDeleteRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}
