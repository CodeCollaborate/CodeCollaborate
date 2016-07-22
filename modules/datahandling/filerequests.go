package datahandling

import (
	"time"

	"strconv"

	"github.com/CodeCollaborate/Server/modules/dbfs"
)

var fileRequestsSetup = false
var newFileVersion int64

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

func (f fileCreateRequest) process() ([](func(dh DataHandler) error), error) {
	// TODO (normal/required): check if permission high enough on project
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"
	not.RoutingKey = strconv.FormatInt(f.ProjectID, 10)

	fileID, err := dbfs.MySQLFileCreate(f.SenderID, f.Name, f.RelativePath, f.ProjectID)
	if err != nil {
		res.ServerMessage = response{
			Status: fail,
			Tag:    f.Tag,
			Data:   struct{}{}}
		return accumulate(toSenderCont(res)), nil
	}

	err = dbfs.CBInsertNewFile(fileID, newFileVersion, make([]string, 0))

	if err != nil {
		res.ServerMessage = response{
			Status: servfail,
			Tag:    f.Tag,
			Data:   struct{}{}}
		return accumulate(toSenderCont(res)), nil
	}
	res.ServerMessage = response{
		Status: success,
		Tag:    f.Tag,
		Data: struct {
			FileID int64
		}{
			fileID,
		}}
	not.ServerMessage = notification{
		Resource: f.Resource,
		Method:   f.Method,
		Data: struct {
			FileID       int64
			ProjectID    int64
			Name         string
			RelativePath string
			Version      int64
		}{
			fileID,
			f.ProjectID,
			f.Name,
			f.RelativePath,
			newFileVersion,
		}}
	return accumulate(toSenderCont(res), toChanCont(not)), nil
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

func (f fileRenameRequest) process() ([](func(dh DataHandler) error), error) {
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	res.ServerMessage = response{
		Status: fail,
		Tag:    f.Tag,
		Data:   struct{}{}}

	fileMeta, err := dbfs.MySQLFileGetInfo(f.FileID)
	if err != nil {
		return accumulate(toSenderCont(res)), nil
	}

	not.RoutingKey = strconv.FormatInt(fileMeta.ProjectID, 10)
	// TODO (normal/required): check if permission high enough on project (fileMeta.ProjectID)

	err = dbfs.MySQLFileRename(f.FileID, f.NewName)
	if err != nil {
		return accumulate(toSenderCont(res)), nil
	}

	res.ServerMessage = response{
		Status: success,
		Tag:    f.Tag,
		Data:   struct{}{}}
	not.ServerMessage = notification{
		Resource: f.Resource,
		Method:   f.Method,
		Data: struct {
			FileID  int64
			NewPath string
		}{
			f.FileID,
			f.NewName,
		}}
	return accumulate(toSenderCont(res), toChanCont(not)), nil
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

func (f fileMoveRequest) process() ([](func(dh DataHandler) error), error) {
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	res.ServerMessage = response{
		Status: fail,
		Tag:    f.Tag,
		Data:   struct{}{}}

	fileMeta, err := dbfs.MySQLFileGetInfo(f.FileID)
	if err != nil {
		return accumulate(toSenderCont(res)), nil
	}

	not.RoutingKey = strconv.FormatInt(fileMeta.ProjectID, 10)
	// TODO (normal/required): check if permission high enough on project (fileMeta.ProjectID)

	err = dbfs.MySQLFileMove(f.FileID, f.NewPath)
	if err != nil {
		return accumulate(toSenderCont(res)), nil
	}
	res.ServerMessage = response{
		Status: success,
		Tag:    f.Tag,
		Data:   struct{}{}}
	not.ServerMessage = notification{
		Resource: f.Resource,
		Method:   f.Method,
		Data: struct {
			FileID  int64
			NewPath string
		}{
			f.FileID,
			f.NewPath,
		}}
	return accumulate(toSenderCont(res), toChanCont(not)), nil
}

// File.Delete
type fileDeleteRequest struct {
	FileID int64
	abstractRequest
}

func (f *fileDeleteRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileDeleteRequest) process() ([](func(dh DataHandler) error), error) {
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	res.ServerMessage = response{
		Status: fail,
		Tag:    f.Tag,
		Data:   struct{}{}}

	fileMeta, err := dbfs.MySQLFileGetInfo(f.FileID)
	if err != nil {
		return accumulate(toSenderCont(res)), err
	}

	not.RoutingKey = strconv.FormatInt(fileMeta.ProjectID, 10)
	// TODO (normal/required): check if permission high enough on project (fileMeta.ProjectID)

	err = dbfs.MySQLFileDelete(f.FileID)
	if err != nil {
		return accumulate(toSenderCont(res)), err
	}

	err = dbfs.CBDeleteFile(f.FileID)
	if err != nil {
		return accumulate(toSenderCont(res)), err
	}

	err = dbfs.FileDelete(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)
	if err != nil {
		return accumulate(toSenderCont(res)), err
	}

	res.ServerMessage = response{
		Status: success,
		Tag:    f.Tag,
		Data:   struct{}{}}
	not.ServerMessage = notification{
		Resource: f.Resource,
		Method:   f.Method,
		Data: struct {
			FileID int64
		}{
			f.FileID,
		}}
	return accumulate(toSenderCont(res), toChanCont(not)), nil
}

// File.Change
type fileChangeRequest struct {
	FileID          int64
	Changes         []string
	BaseFileVersion int64
	abstractRequest
}

func (f *fileChangeRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileChangeRequest) process() ([](func(dh DataHandler) error), error) {
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	res.ServerMessage = response{
		Status: fail,
		Tag:    f.Tag,
		Data:   struct{}{}}

	fileMeta, err := dbfs.MySQLFileGetInfo(f.FileID)
	if err != nil {
		return accumulate(toSenderCont(res)), err
	}

	not.RoutingKey = strconv.FormatInt(fileMeta.ProjectID, 10)
	// TODO (normal/required): check if permission high enough on project (fileMeta.ProjectID)

	// TODO (normal/required): verify changes are valid changes
	version, err := dbfs.CBAppendFileChange(f.FileID, f.BaseFileVersion, f.Changes)
	if err != nil {
		return accumulate(toSenderCont(res)), err
	}

	res.ServerMessage = response{
		Status: success,
		Tag:    f.Tag,
		Data: struct {
			FileVersion int64
		}{
			version,
		}}

	not.ServerMessage = notification{
		Resource: f.Resource,
		Method:   f.Method,
		Data: struct {
			FileID          int64
			BaseFileVersion int64
			FileVersion     int64
			Changes         []string
		}{
			f.FileID,
			f.BaseFileVersion,
			version,
			f.Changes,
		}}

	return accumulate(toSenderCont(res), toChanCont(not)), nil
}

// File.Pull
type filePullRequest struct {
	FileID int64
	abstractRequest
}

func (f *filePullRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f filePullRequest) process() ([](func(dh DataHandler) error), error) {
	// TODO (normal/required): check if permission high enough on project

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	res.ServerMessage = response{
		Status: fail,
		Tag:    f.Tag,
		Data:   struct{}{}}

	fileMeta, err := dbfs.MySQLFileGetInfo(f.FileID)
	if err != nil {
		return accumulate(toSenderCont(res)), err
	}

	rawFile, err := dbfs.FileRead(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)
	if err != nil {
		return accumulate(toSenderCont(res)), err
	}

	changes, err := dbfs.CBGetFileChanges(f.FileID)
	if err != nil {
		return accumulate(toSenderCont(res)), err
	}

	res.ServerMessage = response{
		Status: success,
		Tag:    f.Tag,
		Data: struct {
			FileBytes []byte
			Changes   []string
		}{
			*rawFile,
			changes,
		}}

	return accumulate(toSenderCont(res)), nil
}
