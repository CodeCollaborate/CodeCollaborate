package datahandling

import "encoding/json"

type rabbitCommand struct {
	Command string
	Tag     int64
	Data    json.RawMessage
}
type rabbitQueueData struct {
	Key string
}