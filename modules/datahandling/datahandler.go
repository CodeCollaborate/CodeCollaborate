package datahandling

import "fmt"

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
	return nil
}
