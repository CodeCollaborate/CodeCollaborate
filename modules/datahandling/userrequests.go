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

	unauthenticatedRequestMap["UserRegister"] = func(req AbstractRequest) (Request, error) {
		p := new(userRegisterRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	unauthenticatedRequestMap["UserLogin"] = func(req AbstractRequest) (Request, error) {
		p := new(userLoginRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["UserLookup"] = func(req AbstractRequest) (Request, error) {
		p := new(userLookupRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	requestMap["UserProjects"] = func(req AbstractRequest) (Request, error) {
		p := new(userProjectsRequest)
		p.AbstractRequest = req
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
	AbstractRequest
}

func (p userRegisterRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Printf("Recieved register request from %s\n", p.Username)
	return nil, nil, nil
}

// User.Login
type userLoginRequest struct {
	Username string
	Password string
	AbstractRequest
}

func (p userLoginRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Printf("Recieved login request from %s\n", p.Username)
	return nil, nil, nil
}

// User.Lookup
type userLookupRequest struct {
	Usernames []int64
	AbstractRequest
}

func (p userLookupRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Printf("Recieved user lookup request from %s\n", p.SenderID)
	return nil, nil, nil
}

// User.Projects
type userProjectsRequest struct {
	AbstractRequest
}

func (p userProjectsRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Printf("Recieved user projects request from %s\n", p.SenderID)
	return nil, nil, nil
}
