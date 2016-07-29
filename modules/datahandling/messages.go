package datahandling

import (
	"encoding/json"

	"errors"

	"github.com/CodeCollaborate/Server/modules/dbfs"
)

/**
 * describes the structs and itnerfaces used in the data handling
 */

// Request should be implemented by all request models.
// Provides standard interface for calling the processing
type request interface {
	process(db dbfs.DBFS) (continuations [](func(dh DataHandler) error), err error)
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
	Data        json.RawMessage
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

// ServerMessageWrapper provides interfaces of messages sent from the server
// This section provides the struct definitions of server replies
type serverMessageWrapper struct {
	Type          string
	Timestamp     int64
	RoutingKey    string
	ServerMessage serverMessage
}

type serverMessage interface {
	serverMessageType() string
}

// Response is the type which is the server responses to the client
type response struct {
	Tag    int64
	Status int
	Data   interface{}
}

func (message response) serverMessageType() string {
	return "Response"
}

// Notification is the type which is the unprompted server messages to clients
type notification struct {
	Resource string
	Method   string
	Data     interface{}
}

func (message notification) serverMessageType() string {
	return "Notification"
}

/**
 * Status codes
 */
// success
const success int = 200
const accepted int = 202

// failure
const fail int = 400
const unauthorized int = 401
const notFound int = 404
const fileVersionOutOfDate int = 409 // (409 = conflict)
const partialfail int = 499

// server failure
const servfail int = 500
const unimplemented = 501
const servpartialfail int = 599

/**
 * Errors
 */

// ErrAuthenticationFailed is thrown when the user does not have the proper access to run a request
var ErrAuthenticationFailed = errors.New("No entries were correctly altered")
