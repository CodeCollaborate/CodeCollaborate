package datahandling

import "fmt"

/**
 * Handle the data received by the WebSocket connection.
 * @author: Austin Fahsl and Benedict Wong
 */

type DataHandler struct {
}

/**
 * Handle the data received by the WebSocket connection.
 */
func (dh DataHandler) Handle(wsId uint64, messageType int, message []byte) error {
	fmt.Printf("Handling Message: %s\n", message)
	return nil
}
