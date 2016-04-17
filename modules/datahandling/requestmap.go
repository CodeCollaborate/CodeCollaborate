package datahandling

import (
	"encoding/json"
	"errors"
	"fmt"
)

/**
 * requestmap.go provides the pseudo-factory map for looking up the associated request
 */

// pseudo
var requestMap = make(map[string](func(req AbstractRequest)(Request, error)))

func init() {
	requestMap["ProjectLookup"] = func(req AbstractRequest) (Request, error) {
		p := new(ProjectLookupRequest)
		p.AbstractRequest = req
		rawData := req.Data
		err := json.Unmarshal(rawData, &p)
		return p, err
	}
}

func GetRequestMap(req AbstractRequest) (Request, error) {
	constructor := requestMap[req.Resource + req.Method]
	if (constructor == nil) {
		err := errors.New("The function for the given request does not exist in the map.")
		return  nil, err
	}
	request, err := constructor(req)
	return request, err
}

type ProjectLookupRequest struct {
	ProjectId []uint64
	AbstractRequest
}

func (p ProjectLookupRequest) Process() (err error) {
	fmt.Println(p.ProjectId)
	return nil
}