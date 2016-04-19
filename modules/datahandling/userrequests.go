package datahandling

import (
	"fmt"
)

var userRequestsSetup = false

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initUserRequests() {
	if userRequestsSetup {
		return
	}

	unauthenticatedRequestMap["User.Register"] = func(req abstractRequest) (request, error) {
		return commonJson(new(userRegisterRequest), req)
	}

	unauthenticatedRequestMap["User.Login"] = func(req abstractRequest) (request, error) {
		return commonJson(new(userLoginRequest), req)
	}

	authenticatedRequestMap["User.Lookup"] = func(req abstractRequest) (request, error) {
		return commonJson(new(userLookupRequest), req)
	}

	authenticatedRequestMap["User.Projects"] = func(req abstractRequest) (request, error) {
		return commonJson(new(userProjectsRequest), req)
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

func (f *userRegisterRequest)setAbstractRequest(req abstractRequest) {
	f.abstractRequest = req
}

func (p userRegisterRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved register request from %s\n", p.Username)
	return nil, nil, nil
}

// User.Login
type userLoginRequest struct {
	Username string
	Password string
	abstractRequest
}

func (f *userLoginRequest)setAbstractRequest(req abstractRequest) {
	f.abstractRequest = req
}

func (p userLoginRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved login request from %s\n", p.Username)
	return nil, nil, nil
}

// User.Lookup
type userLookupRequest struct {
	Usernames []string
	abstractRequest
}

func (f *userLookupRequest)setAbstractRequest(req abstractRequest) {
	f.abstractRequest = req
}

func (p userLookupRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved user lookup request from %s\n", p.SenderID)
	return nil, nil, nil
}

// User.Projects
type userProjectsRequest struct {
	abstractRequest
}

func (f *userProjectsRequest)setAbstractRequest(req abstractRequest) {
	f.abstractRequest = req
}

func (p userProjectsRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved user projects request from %s\n", p.SenderID)
	return nil, nil, nil
}
