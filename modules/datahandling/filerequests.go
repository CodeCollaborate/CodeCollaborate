package datahandling

import (
	"fmt"
	"time"

	"github.com/CodeCollaborate/Server/modules/dbfs"
)

var fileRequestsSetup = false
var newFileVersion int

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

func (f fileCreateRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	// TODO: check if permission high enough on project
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	fileID, err := dbfs.MySQLFileCreate(f.SenderID, f.Name, f.RelativePath, f.ProjectID)
	if err != nil {
		res.ServerMessage = response{
			Status: fail,
			Tag:    f.Tag,
			Data:   struct{}{}}

		return res, nil, nil
	}

	err = dbfs.CBInsertNewFile(fileID, newFileVersion, make([]string, 0))

	if err != nil {
		res.ServerMessage = response{
			Status: servfail,
			Tag:    f.Tag,
			Data:   struct{}{}}
		not = nil
	} else {
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
	}

	return res, not, nil
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

func (f fileRenameRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	// TODO: check if permission high enough on project
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	err := dbfs.MySQLFileRename(f.FileID, f.NewName)
	if err != nil {
		res.ServerMessage = response{
			Status: fail,
			Tag:    f.Tag,
			Data:   struct{}{}}
		not = nil
	} else {
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
	}

	return res, not, nil
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

func (f fileMoveRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	// TODO: check if permission high enough on project
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	not := new(serverMessageWrapper)
	not.Timestamp = res.Timestamp
	not.Type = "Notification"

	err := dbfs.MySQLFileMove(f.FileID, f.NewPath)
	if err != nil {
		res.ServerMessage = response{
			Status: fail,
			Tag:    f.Tag,
			Data:   struct{}{}}
		not = nil
	} else {
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
	}

	return res, not, nil
}

// File.Delete
type fileDeleteRequest struct {
	FileID int64
	abstractRequest
}

func (f *fileDeleteRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileDeleteRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	// TODO: check if permission high enough on project
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
		return res, nil, nil
	}

	err = dbfs.MySQLFileDelete(f.FileID)
	if err != nil {
		return res, nil, nil
	}

	err = dbfs.CBDeleteFile(f.FileID)
	if err != nil {
		return res, nil, nil
	}

	fmt.Println("not actually deleting " + fileMeta.Filename + " yet")
	// TODO: delete from filesystem with
	//err = dbfs.FileDelete(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)
	if err != nil {
		return res, nil, nil
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

	return res, not, nil
}

// File.Change
type fileChangeRequest struct {
	FileID      int64
	Changes     []string
	FileVersion int64
	abstractRequest
}

func (f *fileChangeRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f fileChangeRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	// TODO: check if permission high enough on project

	// TODO: decide if file version stays in request (it shouldn't!!)
	// TODO: change to increment version inside cb.bucket.MutateIn call
	// TODO: have it return new version in response

	fmt.Printf("Recieved file change request from %s\n", f.SenderID)
	return nil, nil, nil
}

// File.Pull
type filePullRequest struct {
	FileID int64
	abstractRequest
}

func (f *filePullRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f filePullRequest) process() (*serverMessageWrapper, *serverMessageWrapper, error) {
	// TODO: check if permission high enough on project

	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().UnixNano()
	res.Type = "Responce"

	res.ServerMessage = response{
		Status: fail,
		Tag:    f.Tag,
		Data:   struct{}{}}

	fileMeta, err := dbfs.MySQLFileGetInfo(f.FileID)
	if err != nil {
		return res, nil, nil
	}

	rawFile, err := dbfs.FileRead(fileMeta.RelativePath, fileMeta.Filename, fileMeta.ProjectID)
	if err != nil {
		return res, nil, nil
	}

	changes, err := dbfs.CBGetFileChanges(f.FileID)
	if err != nil {
		return res, nil, nil
	}

	res.ServerMessage = response{
		Status: success,
		Tag:    f.Tag,
		Data: struct {
			FileBytes []byte
			Changes   []string
		}{
			rawFile,
			changes,
		}}

	return res, nil, nil
}
