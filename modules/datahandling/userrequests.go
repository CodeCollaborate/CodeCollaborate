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

	requestMap["UserLookup"] = func(req AbstractRequest) (Request, error) {
		p := new(userLookupRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	userRequestsSetup = true
}

// User.Lookup
type userLookupRequest struct {
	Usernames []int64
	AbstractRequest
}

func (p userLookupRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Println(p.Usernames)
	return nil, nil, nil
}
