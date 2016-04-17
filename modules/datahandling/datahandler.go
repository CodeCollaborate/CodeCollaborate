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
		return nil
	}
	fmt.Print(req)

	return nil
}
