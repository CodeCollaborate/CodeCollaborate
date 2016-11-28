package messages

import "time"

// ServerMessageWrapper provides interfaces of messages sent from the server
type ServerMessageWrapper struct {
	Type          string
	Timestamp     int64
	ServerMessage ServerMessage
}

// ServerMessage is the interface of all messages that the server sends to the client (Responses + Notifications)
type ServerMessage interface {
	Wrap() *ServerMessageWrapper
}

// Response is the type which is the server responses to the client
type Response struct {
	Tag    int64
	Status int
	Data   interface{}
}

// Wrap builds the server message wrapper for this Response struct
func (message Response) Wrap() *ServerMessageWrapper {
	return &ServerMessageWrapper{
		Timestamp:     time.Now().Unix(),
		Type:          "Response",
		ServerMessage: message,
	}
}

// NewEmptyResponse creates a new response, with the specified status and tag.
// The generated response contains no data, and is already wrapped.
func NewEmptyResponse(status int, tag int64) *ServerMessageWrapper {
	return Response{
		Status: status,
		Tag:    tag,
		Data:   struct{}{},
	}.Wrap()
}

/**
 * Status codes
 */

// StatusSuccess represents a successful outcome
const StatusSuccess int = 200

// StatusWrongRequest represents the case where a request was made incorrectly called in favor of the correct one
const StatusWrongRequest int = 301

// StatusFail represents a outcome that failed to process
const StatusFail int = 400

// StatusUnauthorized represents a outcome that could not be authenticated
const StatusUnauthorized int = 401

// StatusNotFound represents a state in which the specified resource was not found
const StatusNotFound int = 404

// StatusVersionOutOfDate represents a state in which the client has an outdated version of the resource
const StatusVersionOutOfDate int = 409 // (409 = conflict)

// StatusPartialFail represents a partial failure in processing the request
const StatusPartialFail int = 499

// StatusServFail represents an internal failure in processing.
const StatusServFail int = 500

// StatusUnimplemented represents a called method that has not yet been implemented
const StatusUnimplemented = 501

// StatusServPartialFail represents an internal failure in processing part of the request.
const StatusServPartialFail int = 599
