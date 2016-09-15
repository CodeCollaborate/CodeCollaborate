package datahandling

import (
	"fmt"

	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
)

/**
 * Data Handling logic for the CodeCollaborate Server.
 */

// DataHandler handles the json data received from the WebSocket connection.
type DataHandler struct {
	MessageChan      chan<- rabbitmq.AMQPMessage
	SubscriptionChan chan<- rabbitmq.Subscription
	WebsocketID      uint64
	Db               dbfs.DBFS
}

// Handle takes the WebSocket Id, MessageType and message in byte-array form,
// processing the data, and updating DB/FS/RabbitMQ as needed.
func (dh DataHandler) Handle(messageType int, message []byte) error {
	fmt.Printf("Handling Message: %s\n", message)

	req, err := createAbstractRequest(message)
	if err != nil {
		utils.LogOnError(err, "Failed to parse json")
		return err
	}

	// automatically determines if the request is authenticated or not
	fullRequest, err := getFullRequest(req)

	var closures []dhClosure

	if err != nil {
		// TODO(shapiro): create response and notification factory
		if err == ErrAuthenticationFailed {
			utils.LogOnError(err, "User not logged in")
			closures = []dhClosure{toSenderClosure{msg: newEmptyResponse(unauthorized, req.Tag)}}
		} else {
			utils.LogOnError(err, "Failed to construct full request")
			closures = []dhClosure{toSenderClosure{msg: newEmptyResponse(unimplemented, req.Tag)}}
		}
	} else {
		closures, err = fullRequest.process(dh.Db)
		if err != nil {
			utils.LogOnError(err, "Failed to handle process request")
			// TODO: forward error message onto client? (or at least inform that error occurred)
		}
	}

	for _, closure := range closures {
		err := closure.call(dh)
		if err != nil {
			utils.LogOnError(err, "Failed to complete continuation")
		}
	}

	return err
}

func authenticate(abs abstractRequest) bool {
	fmt.Println("AUTHENTICATION IS NOT IMPLEMENTED YET")
	// TODO (non-immediate/required): implement user authentication
	return true
}
