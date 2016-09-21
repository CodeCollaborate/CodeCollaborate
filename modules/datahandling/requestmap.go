package datahandling

import (
	"errors"
)

/**
 * provides the pseudo-factory map for looking up the associated request
 */

// flag to disable authentication for testing purposes
var disableAuth = false

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
	if _, contains := unauthenticatedRequestMap[(*req).Resource+"."+(*req).Method]; contains {
		// unauthenticated request
		return unauthenticatedRequest(req)
	}

	// authenticated request
	if disableAuth || authenticate(*req) == nil {
		return authenticatedRequest(req)
	}
	return nil, ErrAuthenticationFailed
}

// authenticatedRequest returns fully parsed Request from the given authenticated AbstractRequest
func authenticatedRequest(req *abstractRequest) (request, error) {
	constructor := authenticatedRequestMap[(*req).Resource+"."+(*req).Method]
	if constructor == nil {
		err := errors.New("The function for " + req.Resource + "." + req.Method + " does not exist in the authenticated map.")
		return nil, err
	}
	request, err := constructor(req)
	return request, err
}

// unauthenticatedRequest returns fully parsed Request from the given unauthenticated AbstractRequest
func unauthenticatedRequest(req *abstractRequest) (request, error) {
	constructor := unauthenticatedRequestMap[(*req).Resource+"."+(*req).Method]
	if constructor == nil {
		err := errors.New("The function for " + req.Resource + "." + req.Method + " does not exist in the unauthenticated map.")
		return nil, err
	}
	request, err := constructor(req)
	return request, err
}
