package rabbitmq

import "encoding/json"

// RabbitCommandJSON represents the RabbitCommand with the Data Struct left unparsed
type RabbitCommandJSON struct {
	Command string
	Tag     int64
	Data    json.RawMessage
}

// RabbitCommandStruct represents the RabbitCommand with the Data Struct included
type RabbitCommandStruct struct {
	Command string
	Tag     int64
	Data    interface{}
}

// RabbitQueueData represents the data needed to identify a specific RabbitMQ queue.
type RabbitQueueData struct {
	Key string
}
