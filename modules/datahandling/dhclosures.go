package datahandling

import (
	"encoding/json"

	"github.com/CodeCollaborate/Server/modules/rabbitmq"
)

type dhClosure interface {
	call(dh DataHandler) error
}

type toSenderClosure struct {
	msg *serverMessageWrapper
}

// toSenderClosure.call is the function that will forward a server message back to the client
func (cont toSenderClosure) call(dh DataHandler) error {
	msgJSON, err := json.Marshal(cont.msg)
	if err != nil {
		return err
	}

	dh.MessageChan <- rabbitmq.AMQPMessage{
		Headers:     make(map[string]interface{}),
		RoutingKey:  rabbitmq.RabbitWebsocketQueueName(dh.WebsocketID),
		ContentType: rabbitmq.ContentType_Msg,
		Persistent:  false,
		Message:     msgJSON,
	}
	return nil
}

type toRabbitChannelClosure struct {
	msg *serverMessageWrapper
	key string
}

// toRabbitChannelClosure.call is the function that will forward a server message to a channel based on the given routing key
func (cont toRabbitChannelClosure) call(dh DataHandler) error {
	msgJSON, err := json.Marshal(cont.msg)
	if err != nil {
		return err
	}
	dh.MessageChan <- rabbitmq.AMQPMessage{
		Headers:     make(map[string]interface{}),
		RoutingKey:  cont.key,
		ContentType: rabbitmq.ContentType_Msg,
		Persistent:  false,
		Message:     msgJSON,
	}
	return nil
}

type rabbitCommandClosure struct {
	Command string
	Tag     int64
	Data    interface{}
}

// toRabbitChannelClosure.call is the function that will forward a server message to a channel based on the given routing key
func (cont rabbitCommandClosure) call(dh DataHandler) error {
	msgJSON, err := json.Marshal(cont)
	if err != nil {
		return err
	}
	dh.MessageChan <- rabbitmq.AMQPMessage{
		Headers:     make(map[string]interface{}),
		RoutingKey:  rabbitmq.RabbitWebsocketQueueName(dh.WebsocketID),
		ContentType: rabbitmq.ContentType_Cmd,
		Persistent:  false,
		Message: msgJSON,
	}
	return nil
}