package datahandling

import (
	"encoding/json"
	"fmt"
)

var projectRequestsSetup = false

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initProjectRequests() {
	if projectRequestsSetup {
		return
	}

	requestMap["ProjectLookup"] = func(req AbstractRequest) (Request, error) {
		p := new(projectLookupRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	projectRequestsSetup = true
}

// Project.Create

// Project.Lookup
type projectLookupRequest struct {
	ProjectIds []int64
	AbstractRequest
}

func (p projectLookupRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Println(p.ProjectIds)
	return nil, nil, nil
}
