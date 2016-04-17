package datahandling

import (
	"encoding/json"
	"fmt"
)

var fileRequestsSetup = false

// initProjectRequests populates the requestMap from requestmap.go with the appropriate constructors for the project methods
func initFileRequests() {
	if fileRequestsSetup {
		return
	}

	requestMap["FileCreate"] = func(req AbstractRequest) (Request, error) {
		p := new(fileCreateRequest)
		p.AbstractRequest = req
		rawData := p.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}

	fileRequestsSetup = true
}

// File.Create
type fileCreateRequest struct {
	Name         string
	RelativePath string
	ProjectID    string
	FileBytes    []byte
	AbstractRequest
}

func (p fileCreateRequest) Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error) {
	fmt.Println(p.Name)
	return nil, nil, nil
}
