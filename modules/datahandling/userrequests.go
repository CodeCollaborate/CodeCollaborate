package datahandling

import (
	"fmt"

	"github.com/CodeCollaborate/Server/modules/dbfs"
)

var userRequestsSetup = false

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initUserRequests() {
	if userRequestsSetup {
		return
	}

	unauthenticatedRequestMap["User.Register"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(userRegisterRequest), req)
	}

	unauthenticatedRequestMap["User.Login"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(userLoginRequest), req)
	}

	authenticatedRequestMap["User.Lookup"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(userLookupRequest), req)
	}

	authenticatedRequestMap["User.Projects"] = func(req *abstractRequest) (request, error) {
		return commonJSON(new(userProjectsRequest), req)
	}

	userRequestsSetup = true
}

// User.Register
type userRegisterRequest struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string
	abstractRequest
}

func (f *userRegisterRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f userRegisterRequest) process(db dbfs.DBFS) ([]dhClosure, error) {

	newUser := dbfs.UserMeta{
		Username:  f.Username,
		FirstName: f.FirstName,
		LastName:  f.LastName,
		Email:     f.Email,
		Password:  f.Password}

	// TODO (non-immediate/required): password validation

	err := db.MySQLUserRegister(newUser)

	if err != nil {
		if err == dbfs.ErrNoDbChange {
			return []dhClosure{toSenderClosure{msg: newEmptyResponse(notFound, f.Tag)}}, err
		}
		return []dhClosure{toSenderClosure{msg: newEmptyResponse(fail, f.Tag)}}, err
	}
	return []dhClosure{toSenderClosure{msg: newEmptyResponse(success, f.Tag)}}, err
}

// User.Login
type userLoginRequest struct {
	Username string
	Password string
	abstractRequest
}

func (f *userLoginRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f userLoginRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	// TODO (non-immediate/required): implement login logic

	fmt.Printf("Received login request from %s. Login logic not implemented yet.\n", f.Username)

	res := response{
		Status: success,
		Tag:    f.Tag,
		Data: struct {
			Token string
		}{
			Token: "TEST_TOKEN",
		},
	}.wrap()

	return []dhClosure{toSenderClosure{msg: res}}, nil
}

// User.Lookup
type userLookupRequest struct {
	Usernames []string
	abstractRequest
}

func (f *userLookupRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f userLookupRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	users := make([]dbfs.UserMeta, len(f.Usernames))
	index := 0
	var erro error
	for _, username := range f.Usernames {
		usr, err := db.MySQLUserLookup(username)
		if err != nil {
			erro = err
		} else {
			users[index] = usr
			index++
		}
	}
	// shrink as needed
	users = users[:index]

	if len(users) < 0 {
		return []dhClosure{toSenderClosure{msg: newEmptyResponse(fail, f.Tag)}}, erro
	} else if erro != nil {
		// at least 1 value failed
		// return what we can but
		// tell the client whatever they don't get back failed
		res := response{
			Status: partialfail,
			Tag:    f.Tag,
			Data: struct {
				Users []dbfs.UserMeta
			}{
				Users: users,
			},
		}.wrap()
		return []dhClosure{toSenderClosure{msg: res}}, erro
	}

	res := response{
		Status: success,
		Tag:    f.Tag,
		Data: struct {
			Users []dbfs.UserMeta
		}{
			Users: users,
		},
	}.wrap()

	return []dhClosure{toSenderClosure{msg: res}}, erro
}

// User.Projects
type userProjectsRequest struct {
	abstractRequest
}

func (f *userProjectsRequest) setAbstractRequest(req *abstractRequest) {
	f.abstractRequest = *req
}

func (f userProjectsRequest) process(db dbfs.DBFS) ([]dhClosure, error) {
	projects, err := db.MySQLUserProjects(f.SenderID)
	if err != nil {
		res := response{
			Status: partialfail,
			Tag:    f.Tag,
			Data: struct {
				Projects []dbfs.ProjectMeta
			}{
				Projects: projects,
			},
		}.wrap()
		return []dhClosure{toSenderClosure{msg: res}}, err
	}

	res := response{
		Status: success,
		Tag:    f.Tag,
		Data: struct {
			Projects []dbfs.ProjectMeta
		}{
			Projects: projects,
		},
	}.wrap()

	return []dhClosure{toSenderClosure{msg: res}}, err
}
