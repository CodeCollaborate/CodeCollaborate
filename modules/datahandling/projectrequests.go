package datahandling

import (
	"time"

	"strings"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datahandling/messages"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
)

var projectRequestsSetup = false

// TODO(wongb): Create & Use a Project struct

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

	authenticatedRequestMap["Project.GetPermissionConstants"] = func(req *abstractRequest) (request, error) {
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
	if err != nil {
		//if err == project already exists {
		// TODO(shapiro): implement a specific error for this on the mysql.go side
		//}
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusServFail, p.Tag)}}, nil
	}

	res := messages.Response{
		Status: messages.StatusSuccess,
		Tag:    p.Tag,
		Data: struct {
			ProjectID int64
		}{
			ProjectID: projectID,
		},
	}.Wrap()

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
	hasPermission, err := dbfs.PermissionAtLeast(p.SenderID, p.ProjectID, "write", db)
	if err != nil || !hasPermission {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  p.Resource,
			"Method":    p.Method,
			"SenderID":  p.SenderID,
			"ProjectID": p.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, p.Tag)}}, nil
	}

	err = db.MySQLProjectRename(p.ProjectID, p.NewName)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusServFail, p.Tag)}}, err
	}

	res := messages.NewEmptyResponse(messages.StatusSuccess, p.Tag)
	not := messages.Notification{
		Resource:   p.Resource,
		Method:     p.Method,
		ResourceID: p.ProjectID,
		Data: struct {
			NewName string
		}{
			NewName: p.NewName,
		},
	}.Wrap()

	return []dhClosure{toSenderClosure{msg: res}, toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitProjectQueueName(p.ProjectID)}}, nil
}

// Project.GetPermissionConstants
type projectGetPermissionConstantsRequest struct {
	abstractRequest
}

func (p *projectGetPermissionConstantsRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}

func (p projectGetPermissionConstantsRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	res := messages.Response{
		Status: messages.StatusSuccess,
		Tag:    p.Tag,
		Data: struct {
			Constants map[string]int8
		}{
			Constants: config.PermissionsByLabel,
		},
	}.Wrap()

	return []dhClosure{toSenderClosure{msg: res}}, nil
}

// Project.GrantPermissions
type projectGrantPermissionsRequest struct {
	ProjectID       int64
	GrantUsername   string
	PermissionLevel int8
	abstractRequest
}

func (p projectGrantPermissionsRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	hasPermission, err := dbfs.PermissionAtLeast(p.SenderID, p.ProjectID, "admin", db)
	if err != nil || !hasPermission {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  p.Resource,
			"Method":    p.Method,
			"SenderID":  p.SenderID,
			"ProjectID": p.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, p.Tag)}}, nil
	}

	p.GrantUsername = strings.ToLower(p.GrantUsername)

	// Prevent users from changing their own permissions
	if p.SenderID == p.GrantUsername {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, p.Tag)}}, nil
	}

	// TODO: Add if User exists check

	requestPerm, err := config.PermissionByLevel(p.PermissionLevel)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, p.Tag)}}, nil
	}

	ownerPerm, err := config.PermissionByLabel("owner")
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusServFail, p.Tag)}}, nil
	}

	if requestPerm.Level == ownerPerm.Level {
		// TODO(shapiro): implement changing ownership
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnimplemented, p.Tag)}}, nil
	}

	err = db.MySQLProjectGrantPermission(p.ProjectID, p.GrantUsername, p.PermissionLevel, p.SenderID)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusServFail, p.Tag)}}, err
	}

	res := messages.NewEmptyResponse(messages.StatusSuccess, p.Tag)
	not := messages.Notification{
		Resource:   p.Resource,
		Method:     p.Method,
		ResourceID: p.ProjectID,
		Data: struct {
			GrantUsername   string
			PermissionLevel int8
		}{
			GrantUsername:   p.GrantUsername,
			PermissionLevel: p.PermissionLevel,
		},
	}.Wrap()

	return []dhClosure{
		toSenderClosure{msg: res},
		toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitProjectQueueName(p.ProjectID)},
		toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitUserQueueName(p.GrantUsername)}}, nil
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
	hasPermission, err := dbfs.PermissionAtLeast(p.SenderID, p.ProjectID, "admin", db)
	if err != nil {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  p.Resource,
			"Method":    p.Method,
			"SenderID":  p.SenderID,
			"ProjectID": p.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, p.Tag)}}, nil
	}

	p.RevokeUsername = strings.ToLower(p.RevokeUsername)

	// allow case where user is removing themselves from a project
	if !hasPermission && p.SenderID != p.RevokeUsername {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, p.Tag)}}, nil
	}

	err = db.MySQLProjectRevokePermission(p.ProjectID, p.RevokeUsername, p.SenderID)

	if err != nil {
		if err == dbfs.ErrNoDbChange {
			ownerPerm, err := config.PermissionByLabel("owner")
			if err != nil {
				return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusServFail, p.Tag)}}, nil
			}

			_, permissions, err := db.MySQLProjectLookup(p.ProjectID, p.SenderID)
			if err != nil {
				return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, p.Tag)}}, err
			}
			for username, lvl := range permissions {
				if lvl.PermissionLevel == ownerPerm.Level && username == p.RevokeUsername {
					// request is trying to remove owner, we can be more specific in errors
					if p.SenderID == username {
						// the owner is trying to remove themselves
						// NOTE: we could do a project delete here? but seems weird
						return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusWrongRequest, p.Tag)}}, err
					}
					return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, p.Tag)}}, err
				}
			}
		}
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusServFail, p.Tag)}}, err
	}

	res := messages.NewEmptyResponse(messages.StatusSuccess, p.Tag)
	not := messages.Notification{
		Resource:   p.Resource,
		Method:     p.Method,
		ResourceID: p.ProjectID,
		Data: struct {
			RevokeUsername string
		}{
			RevokeUsername: p.RevokeUsername,
		},
	}.Wrap()

	unsubscribeCommand := rabbitCommandClosure{
		Command: "Unsubscribe",
		Tag:     -1,
		Key:     rabbitmq.RabbitUserQueueName(p.RevokeUsername),
		Data: rabbitmq.RabbitQueueData{
			Key: rabbitmq.RabbitProjectQueueName(p.ProjectID),
		},
	}

	return []dhClosure{
		toSenderClosure{msg: res},
		toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitProjectQueueName(p.ProjectID)},
		toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitUserQueueName(p.RevokeUsername)},
		unsubscribeCommand}, nil
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
	utils.LogWarn("ProjectGetOnlineClients not implemented", nil)

	return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnimplemented, p.Tag)}}, nil
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
		// it's better to do a cheap lookup and then an expensive one if required than an expensive one every time
		hasPermission, err := dbfs.PermissionAtLeast(p.SenderID, id, "read", db)
		if err != nil || !hasPermission {
			utils.LogError("API permission error", err, utils.LogFields{
				"Resource":  p.Resource,
				"Method":    p.Method,
				"SenderID":  p.SenderID,
				"ProjectID": id,
			})
			errOut = ErrAuthenticationFailed
			continue
		}

		lookupResult, err := projectLookup(p.SenderID, id, db)
		if err != nil {
			errOut = err
		} else {
			resultData[i] = lookupResult
			i++
		}
	}
	// shrink to cut off remainder left by errors
	resultData = resultData[:i]

	if errOut != nil {
		if len(resultData) == 0 {
			res := messages.Response{
				Status: messages.StatusFail,
				Tag:    p.Tag,
				Data: struct {
					Projects []projectLookupResult
				}{
					Projects: resultData,
				},
			}.Wrap()
			return []dhClosure{toSenderClosure{msg: res}}, nil
		}
		res := messages.Response{
			Status: messages.StatusPartialFail,
			Tag:    p.Tag,
			Data: struct {
				Projects []projectLookupResult
			}{
				Projects: resultData,
			},
		}.Wrap()
		return []dhClosure{toSenderClosure{msg: res}}, nil
	}

	res := messages.Response{
		Status: messages.StatusSuccess,
		Tag:    p.Tag,
		Data: struct {
			Projects []projectLookupResult
		}{
			Projects: resultData,
		},
	}.Wrap()

	return []dhClosure{toSenderClosure{msg: res}}, nil
}

func projectLookup(senderID string, projectID int64, db dbfs.DBFS) (projectLookupResult, error) {
	var result projectLookupResult

	name, permissions, err := db.MySQLProjectLookup(projectID, senderID)

	if err != nil {
		return result, err
	}

	result = projectLookupResult{
		ProjectID:   projectID,
		Name:        name,
		Permissions: permissions,
	}

	return result, nil
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
	Version      int64
}

func (p projectGetFilesRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	hasPermission, err := dbfs.PermissionAtLeast(p.SenderID, p.ProjectID, "read", db)
	if err != nil || !hasPermission {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  p.Resource,
			"Method":    p.Method,
			"SenderID":  p.SenderID,
			"ProjectID": p.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, p.Tag)}}, nil
	}

	files, err := db.MySQLProjectGetFiles(p.ProjectID)
	if err != nil {
		res := messages.Response{
			Status: messages.StatusFail,
			Tag:    p.Tag,
			Data: struct {
				Files []fileLookupResult
			}{
				Files: make([]fileLookupResult, 0),
			},
		}.Wrap()

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
				Version:      version}
			i++
		}
	}
	// shrink to cut off remainder left by errors
	resultData = resultData[:i]

	if errOut != nil {
		if len(resultData) == 0 {
			res := messages.Response{
				Status: messages.StatusFail,
				Tag:    p.Tag,
				Data: struct {
					Files []fileLookupResult
				}{
					Files: resultData,
				},
			}.Wrap()
			return []dhClosure{toSenderClosure{msg: res}}, nil
		}
		res := messages.Response{
			Status: messages.StatusPartialFail,
			Tag:    p.Tag,
			Data: struct {
				Files []fileLookupResult
			}{
				Files: resultData,
			},
		}.Wrap()
		return []dhClosure{toSenderClosure{msg: res}}, nil
	}
	res := messages.Response{
		Status: messages.StatusSuccess,
		Tag:    p.Tag,
		Data: struct {
			Files []fileLookupResult
		}{
			Files: resultData,
		},
	}.Wrap()

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
	hasPermission, err := dbfs.PermissionAtLeast(p.SenderID, p.ProjectID, "read", db)
	if err != nil || !hasPermission {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  p.Resource,
			"Method":    p.Method,
			"SenderID":  p.SenderID,
			"ProjectID": p.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, p.Tag)}}, nil
	}

	cmdClosure := rabbitCommandClosure{
		Command: "Subscribe",
		Tag:     p.Tag,
		Data: rabbitmq.RabbitQueueData{
			Key: rabbitmq.RabbitProjectQueueName(p.ProjectID),
		},
	}
	return []dhClosure{cmdClosure}, nil
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
	cmdClosure := rabbitCommandClosure{
		Command: "Unsubscribe",
		Tag:     p.Tag,
		Data: rabbitmq.RabbitQueueData{
			Key: rabbitmq.RabbitProjectQueueName(p.ProjectID),
		},
	}
	return []dhClosure{cmdClosure}, nil
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
	hasPermission, err := dbfs.PermissionAtLeast(p.SenderID, p.ProjectID, "owner", db)
	if err != nil {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  p.Resource,
			"Method":    p.Method,
			"SenderID":  p.SenderID,
			"ProjectID": p.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, p.Tag)}}, nil
	}

	if !hasPermission {
		hasCurrentProjectPermission, err := dbfs.PermissionAtLeast(p.SenderID, p.ProjectID, "read", db)
		if err != nil {
			utils.LogError("API permission error", err, utils.LogFields{
				"Resource":  p.Resource,
				"Method":    p.Method,
				"SenderID":  p.SenderID,
				"ProjectID": p.ProjectID,
			})
			return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, p.Tag)}}, nil
		}

		if hasCurrentProjectPermission {
			// replace this delete request with a self revoke permissions request
			realRequest := projectRevokePermissionsRequest{
				ProjectID:       p.ProjectID,
				RevokeUsername:  p.SenderID,
				abstractRequest: p.abstractRequest,
			}
			return realRequest.process(db)
		}

		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, p.Tag)}}, nil
	}

	err = db.MySQLProjectDelete(p.ProjectID, p.SenderID)
	if err != nil {
		if err == dbfs.ErrNoDbChange {
			return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, p.Tag)}}, err
		}
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusServFail, p.Tag)}}, err

	}

	res := messages.NewEmptyResponse(messages.StatusSuccess, p.Tag)
	not := messages.Notification{
		Resource:   p.Resource,
		Method:     p.Method,
		ResourceID: p.ProjectID,
		Data:       struct{}{},
	}.Wrap()

	return []dhClosure{toSenderClosure{msg: res}, toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitProjectQueueName(p.ProjectID)}}, nil
}

func (p *projectDeleteRequest) setAbstractRequest(req *abstractRequest) {
	p.abstractRequest = *req
}
