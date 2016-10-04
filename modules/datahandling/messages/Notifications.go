package messages

import "time"

// Notification is the type which is the unprompted server messages to clients
type Notification struct {
	Resource   string
	Method     string
	ResourceID int64
	Data       interface{}
}

// Wrap builds the server message wrapper for this Notification struct
func (message Notification) Wrap() *ServerMessageWrapper {
	return &ServerMessageWrapper{
		Timestamp:     time.Now().Unix(),
		Type:          "Notification",
		ServerMessage: message,
	}
}

//
