package datahandling

import (
	"errors"
)

/**
 * requestmap.go provides the pseudo-factory map for looking up the associated request
 */

// map to lookup authenticated api functions
var authenticatedRequestMap = make(map[string](func(req *abstractRequest) (request, error)))

// map to lookup unauthenticated api functions
var unauthenticatedRequestMap = make(map[string](func(req *abstractRequest) (request, error)))

func init() {
	initProjectRequests()
	initUserRequests()
	initFileRequests()
}

func getFullRequest(req *abstractRequest) (request, error) {
	if (*req).SenderToken == "" {
		// unauthenticated request
		return unauthenticatedRequest(req)
	}
	// authenticated request
	if authenticate(*req) {
		return authenticatedRequest(req)
	}

	return nil, errors.New("Cannot authenticate user")
}

// authenticatedRequest returns fully parsed Request from the given authenticated AbstractRequest
func authenticatedRequest(req *abstractRequest) (request, error) {
	constructor := authenticatedRequestMap[(*req).Resource+"."+(*req).Method]
	if constructor == nil {
		err := errors.New("The function for the given request does not exist in the authenticated map.")
		return nil, err
	}
	request, err := constructor(req)
	return request, err
}

// unauthenticatedRequest returns fully parsed Request from the given unauthenticated AbstractRequest
func unauthenticatedRequest(req *abstractRequest) (request, error) {
	constructor := unauthenticatedRequestMap[(*req).Resource+"."+(*req).Method]
	if constructor == nil {
		err := errors.New("The function for the given request does not exist in the unauthenticated map.")
		return nil, err
	}
	request, err := constructor(req)
	return request, err
}
