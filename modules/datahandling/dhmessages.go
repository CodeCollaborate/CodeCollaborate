package datahandling

import (
	"encoding/json"
	"errors"

	"github.com/CodeCollaborate/Server/modules/dbfs"
)

/**
 * describes the structs and interfaces which describe messages used in the data handling
 */

// Request should be implemented by all request models.
// Provides standard interface for calling the processing
type request interface {
	process(db dbfs.DBFS, ack func() error) (continuations []dhClosure, err error)
	setAbstractRequest(absReq *abstractRequest)
}

// AbstractRequest is the generic request type
type abstractRequest struct {
	Tag         int64
	Resource    string
	SenderID    string
	SenderToken string
	Method      string
	Timestamp   int64
	Data        json.RawMessage // date is a byte for now because we don't want it to unmarshal it yet
}

// CreateAbstractRequest is the testable parsing into abstractRequests
func createAbstractRequest(jsony []byte) (req *abstractRequest, err error) {
	err = json.Unmarshal(jsony, &req)
	if err != nil {
		return nil, err
	}
	return req, err
}

func commonJSON(req request, absReq *abstractRequest) (request, error) {
	req.setAbstractRequest(absReq)
	rawData := (*absReq).Data
	err := json.Unmarshal(rawData, req)
	return req, err
}

/**
 * Errors
 */

// ErrAuthenticationFailed is thrown when the user does not have the proper access to run a request
var ErrAuthenticationFailed = errors.New("No entries were correctly altered")
