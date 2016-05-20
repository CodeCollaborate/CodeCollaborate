package datahandling

import (
	"fmt"

	"github.com/CodeCollaborate/Server/utils"
)

/**
 * Data Handling logic for the CodeCollaborate Server.
 */

// DataHandler handles the json data received from the WebSocket connection.
type DataHandler struct {
}

// Handle takes the WebSocket Id, MessageType and message in byte-array form,
// processing the data, and updating DB/FS/RabbitMQ as needed.
func (dh DataHandler) Handle(wsID uint64, messageType int, message []byte) error {
	fmt.Printf("Handling Message: %s\n", message)

	req, err := createAbstractRequest(message)
	if err != nil {
		utils.LogOnError(err, "Failed to parse json")
		return err
	}

	// automatically delegates if the request is authenticated or not
	fullR, err := getFullRequest(req)

	if err != nil {
		utils.LogOnError(err, "Failed to construct full request")
		return err
	}

	response, notification, err := fullR.process()

	if err != nil {
		utils.LogOnError(err, "Failed to process request")
	}

	if response != nil {
		// TODO: send on rabbit
	}

	if notification != nil {
		// TODO: send on rabbit
	}

	return err
}

func authenticate(abs abstractRequest) bool {
	fmt.Println("AUTHENTICATION IS NOT IMPLEMENTED YET")
	// TODO: implement this
	return true
}
