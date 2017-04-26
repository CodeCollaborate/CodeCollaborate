package datahandling

import (
	"github.com/CodeCollaborate/Server/modules/datahandling/messages"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
)

var fileRequestsSetup = false
var newFileVersion int64 = 1

// File aggregates information relating to an individual file
// TODO(wongb): Change all responses and notifications to use this struct; add creator and creation date
type File struct {
	FileID       int64
	Filename     string
	RelativePath string
	Version      int64
}

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initFileRequests() {
	if fileRequestsSetup {
		return
	}

	authenticatedRequestMap["File.Create"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(fileCreateRequest), req)
	}

	authenticatedRequestMap["File.Rename"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(fileRenameRequest), req)
	}

	authenticatedRequestMap["File.Move"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(fileMoveRequest), req)
	}

	authenticatedRequestMap["File.Delete"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(fileDeleteRequest), req)
	}

	authenticatedRequestMap["File.Change"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(fileChangeRequest), req)
	}

	authenticatedRequestMap["File.Pull"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(filePullRequest), req)
	}

	fileRequestsSetup = true
}

// File.Create
type fileCreateRequest struct {
	Name         string
	RelativePath string
	ProjectID    int64
	FileBytes    []byte
	abstractRequest
}

func (f *fileCreateRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileCreateRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	hasPermission, err := dbfs.PermissionAtLeast(f.SenderID, f.ProjectID, "write", db)
	if err != nil || !hasPermission {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  f.Resource,
			"Method":    f.Method,
			"SenderID":  f.SenderID,
			"ProjectID": f.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, f.Tag)}}, nil
	}

	fileID, err := db.MySQLFileCreate(f.SenderID, f.Name, f.RelativePath, f.ProjectID)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	err = db.FileWrite(fileID, f.FileBytes)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	err = db.CBInsertNewFile(fileID, newFileVersion, make([]string, 0))
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	res := messages.Response{
		Status: messages.StatusSuccess,
		Tag:    f.Tag,
		Data: struct {
			FileID int64
		}{
			FileID: fileID,
		},
	}.Wrap()
	not := messages.Notification{
		Resource:   f.Resource,
		Method:     f.Method,
		ResourceID: f.ProjectID,
		Data: struct {
			File File
		}{
			File: File{
				FileID:       fileID,
				Filename:     f.Name,
				RelativePath: f.RelativePath,
				Version:      newFileVersion,
			},
		},
	}.Wrap()

	return []dhClosure{toSenderClosure{msg: res}, toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitProjectQueueName(f.ProjectID)}}, nil
}

// File.Rename
type fileRenameRequest struct {
	FileID  int64
	NewName string
	abstractRequest
}

func (f *fileRenameRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileRenameRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	fileMeta, err := db.MySQLFileGetInfo(f.FileID)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	hasPermission, err := dbfs.PermissionAtLeast(f.SenderID, fileMeta.ProjectID, "write", db)
	if err != nil || !hasPermission {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  f.Resource,
			"Method":    f.Method,
			"SenderID":  f.SenderID,
			"ProjectID": fileMeta.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, f.Tag)}}, nil
	}

	err = db.MySQLFileRename(f.FileID, f.NewName)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	res := messages.NewEmptyResponse(messages.StatusSuccess, f.Tag)
	not := messages.Notification{
		Resource:   f.Resource,
		Method:     f.Method,
		ResourceID: f.FileID,
		Data: struct {
			NewName string
		}{
			NewName: f.NewName,
		},
	}.Wrap()

	return []dhClosure{toSenderClosure{msg: res}, toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitProjectQueueName(fileMeta.ProjectID)}}, nil
}

// File.Move
type fileMoveRequest struct {
	FileID  int64
	NewPath string
	abstractRequest
}

func (f *fileMoveRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileMoveRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	fileMeta, err := db.MySQLFileGetInfo(f.FileID)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	hasPermission, err := dbfs.PermissionAtLeast(f.SenderID, fileMeta.ProjectID, "write", db)
	if err != nil || !hasPermission {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  f.Resource,
			"Method":    f.Method,
			"SenderID":  f.SenderID,
			"ProjectID": fileMeta.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, f.Tag)}}, nil
	}

	err = db.MySQLFileMove(f.FileID, f.NewPath)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	res := messages.NewEmptyResponse(messages.StatusSuccess, f.Tag)
	not := messages.Notification{
		Resource:   f.Resource,
		Method:     f.Method,
		ResourceID: f.FileID,
		Data: struct {
			NewPath string
		}{
			NewPath: f.NewPath,
		},
	}.Wrap()

	return []dhClosure{toSenderClosure{msg: res}, toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitProjectQueueName(fileMeta.ProjectID)}}, nil
}

// File.Delete
type fileDeleteRequest struct {
	FileID int64
	abstractRequest
}

func (f *fileDeleteRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileDeleteRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	fileMeta, err := db.MySQLFileGetInfo(f.FileID)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	hasPermission, err := dbfs.PermissionAtLeast(f.SenderID, fileMeta.ProjectID, "write", db)
	if err != nil || !hasPermission {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  f.Resource,
			"Method":    f.Method,
			"SenderID":  f.SenderID,
			"ProjectID": fileMeta.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, f.Tag)}}, nil
	}

	err = db.MySQLFileDelete(f.FileID)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	err = db.FileDelete(f.FileID)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	err = db.CBDeleteFile(f.FileID)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	res := messages.NewEmptyResponse(messages.StatusSuccess, f.Tag)
	not := messages.Notification{
		Resource:   f.Resource,
		Method:     f.Method,
		ResourceID: f.FileID,
		Data:       struct{}{},
	}.Wrap()

	return []dhClosure{
		toSenderClosure{msg: res},
		toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitProjectQueueName(fileMeta.ProjectID)},
	}, nil
}

// File.Change
type fileChangeRequest struct {
	FileID  int64
	Changes string
	abstractRequest
}

func (f *fileChangeRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileChangeRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	// This has to be before the CouchBase append, to make sure that the the two databases are kept in sync.
	// Specifically, this prevents CouchBase from incrementing a version number without the notifications being sent out.
	fileMeta, err := db.MySQLFileGetInfo(f.FileID)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	hasPermission, err := dbfs.PermissionAtLeast(f.SenderID, fileMeta.ProjectID, "write", db)
	if err != nil || !hasPermission {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  f.Resource,
			"Method":    f.Method,
			"SenderID":  f.SenderID,
			"ProjectID": fileMeta.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, f.Tag)}}, nil
	}

	// TODO (normal/optional): verify changes are valid changes
	changes, version, missing, numchanges, err := db.CBAppendFileChange(fileMeta, f.Changes)
	if err != nil {
		if err == dbfs.ErrVersionOutOfDate {
			return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusVersionOutOfDate, f.Tag)}}, err
		} else if err == dbfs.ErrResourceNotFound {
			return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusNotFound, f.Tag)}}, err
		}
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	res := messages.Response{
		Status: messages.StatusSuccess,
		Tag:    f.Tag,
		Data: struct {
			FileVersion    int64
			Changes        string
			MissingPatches []string
		}{
			FileVersion:    version,
			Changes:        changes,
			MissingPatches: missing,
		},
	}.Wrap()
	not := messages.Notification{
		Resource:   f.Resource,
		Method:     f.Method,
		ResourceID: f.FileID,
		Data: struct {
			FileVersion int64
			Changes     string
		}{
			FileVersion: version,
			Changes:     changes,
		},
	}.Wrap()

	// Trigger scrunching if longer than maxBufferLength
	if numchanges > dbfs.MaxBufferLength {
		go func() {
			db.ScrunchFile(fileMeta)
		}()
	}

	return []dhClosure{toSenderClosure{msg: res}, toRabbitChannelClosure{msg: not, key: rabbitmq.RabbitProjectQueueName(fileMeta.ProjectID)}}, nil
}

// File.Pull
type filePullRequest struct {
	FileID int64
	abstractRequest
}

func (f *filePullRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f filePullRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	fileMeta, err := db.MySQLFileGetInfo(f.FileID)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	hasPermission, err := dbfs.PermissionAtLeast(f.SenderID, fileMeta.ProjectID, "read", db)
	if err != nil || !hasPermission {
		utils.LogError("API permission error", err, utils.LogFields{
			"Resource":  f.Resource,
			"Method":    f.Method,
			"SenderID":  f.SenderID,
			"ProjectID": fileMeta.ProjectID,
		})
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusUnauthorized, f.Tag)}}, nil
	}

	rawFile, changes, err := db.PullFile(fileMeta)
	if err != nil {
		return []dhClosure{toSenderClosure{msg: messages.NewEmptyResponse(messages.StatusFail, f.Tag)}}, err
	}

	res := messages.Response{
		Status: messages.StatusSuccess,
		Tag:    f.Tag,
		Data: struct {
			FileBytes []byte
			Changes   []string
		}{
			FileBytes: rawFile,
			Changes:   changes,
		},
	}.Wrap()

	return []dhClosure{toSenderClosure{msg: res}}, nil
}
