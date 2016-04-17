package datahandling

import (
	"encoding/json"
	"github.com/CodeCollaborate/Server/utils"
)

/**
 * Interfaces.go describes the structs and itnerfaces used in the data handling
 */

// ProcessorInterface should be implemented by all request models.
// Provides standard interface for calling the processing.
type Request interface {
	Process() (err error)
}

// Interface which defines the different data blocks
type Data interface {
	UnmarshalData(req *AbstractRequest)
}

// generic request type
type AbstractRequest struct {
	Tag uint64
	Resource string
	SenderId string
	SenderToken string
	Method string
	Time uint64
	Data json.RawMessage // date is a byte for now because we don't want it to unmarshal it yet
}

// testable parsing into abstractRequests
func CreateAbstractRequest(jsony []byte) (req AbstractRequest, err error) {
	err = json.Unmarshal(jsony, &req)
	if err != nil {
		utils.LogOnError(err, "Failed to parse json")
		return
	}
	return
}



// Provides interfaces of messages sent from the server
type ServerMessage interface {
	Send()
}

type Response struct {

}

type Notification struct {

}
