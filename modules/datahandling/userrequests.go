package datahandling

import (
	"encoding/json"
	"fmt"
)

var userRequestsSetup = false

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initUserRequests() {
	if userRequestsSetup {
		return
	}

	unauthenticatedRequestMap["UserRegister"] = func(req abstractRequest) (request, error) {
		p := new(userRegisterRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	unauthenticatedRequestMap["UserLogin"] = func(req abstractRequest) (request, error) {
		p := new(userLoginRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["UserLookup"] = func(req abstractRequest) (request, error) {
		p := new(userLookupRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["UserProjects"] = func(req abstractRequest) (request, error) {
		p := new(userProjectsRequest)
		p.abstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
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

func (p userLoginRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved login request from %s\n", p.Username)
	return nil, nil, nil
}

// User.Lookup
type userLookupRequest struct {
	Usernames []int64
	abstractRequest
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

func (p userProjectsRequest) process() (response *serverMessageWrapper, notification *serverMessageWrapper, err error) {
	// TODO
	fmt.Printf("Recieved user projects request from %s\n", p.SenderID)
	return nil, nil, nil
}
