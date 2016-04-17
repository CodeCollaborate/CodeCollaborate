package datahandling

import (
	"encoding/json"

	"github.com/CodeCollaborate/Server/utils"
)

/**
 * Interfaces.go describes the structs and itnerfaces used in the data handling
 */

// Request should be implemented by all request models.
// Provides standard interface for calling the processing
type Request interface {
	Process() (response *ServerMessageWrapper, notification *ServerMessageWrapper, err error)
}

// AbstractRequest is the generic request type
type AbstractRequest struct {
	Tag         int64
	Resource    string
	SenderID    string
	SenderToken string
	Method      string
	Timestamp   int64
	Data        json.RawMessage // date is a byte for now because we don't want it to unmarshal it yet
}

// CreateAbstractRequest is the testable parsing into abstractRequests
func createAbstractRequest(jsony []byte) (req AbstractRequest, err error) {
	err = json.Unmarshal(jsony, &req)
	if err != nil {
		utils.LogOnError(err, "Failed to parse json")
		return req, err
	}
	return req, err
}

// ServerMessageWrapper provides interfaces of messages sent from the server
// This section provides the struct definitions of server replies
type ServerMessageWrapper struct {
	Type          string
	Timestamp     int64
	ServerMessage serverMessage
}

type serverMessage interface {
	serverMessageType() string
}

// Response is the type which is the server responses to the client
type Response struct {
	Tag    int64
	Status int
	Data   interface{}
}

func (message Response) serverMessageType() string {
	return "Response"
}

// Notification is the type which is the unprompted server messages to clients
type Notification struct {
	Resource string
	Method   string
	Data     interface{}
}

func (message Notification) serverMessageType() string {
	return "Notification"
}
