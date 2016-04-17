package datahandling

import (
	"errors"
)

/**
 * requestmap.go provides the pseudo-factory map for looking up the associated request
 */

// map to lookup (authenticated?) api functions
var requestMap = make(map[string](func(req AbstractRequest) (Request, error)))

func init() {
	initProjectRequests()
	initUserRequests()
	initFileRequests()
}

// GetRequestMap returns fully parsed Request from the given AbstractRequest
func GetRequestMap(req AbstractRequest) (Request, error) {
	constructor := requestMap[req.Resource+req.Method]
	if constructor == nil {
		err := errors.New("The function for the given request does not exist in the map.")
		return nil, err
	}
	request, err := constructor(req)
	return request, err
}
