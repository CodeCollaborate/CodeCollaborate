package datahandling

import (
	"encoding/json"
	"time"

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
		ContentType: cont.msg.Type,
		Persistent:  false,
		Message:     msgJSON,
	}
	return nil
}

type toRabbitChannelClosure struct {
	msg       *serverMessageWrapper
	projectID int64
}

// toRabbitChannelClosure.call is the function that will forward a server message to a channel based on the given routing key
func (cont toRabbitChannelClosure) call(dh DataHandler) error {
	msgJSON, err := json.Marshal(cont.msg)
	if err != nil {
		return err
	}
	dh.MessageChan <- rabbitmq.AMQPMessage{
		Headers:     make(map[string]interface{}),
		RoutingKey:  rabbitmq.RabbitProjectQueueName(cont.projectID),
		ContentType: cont.msg.Type,
		Persistent:  false,
		Message:     msgJSON,
	}
	return nil
}

type rabbitChannelSubscribeClosure struct {
	key string
	tag int64
}

func (cont rabbitChannelSubscribeClosure) call(dh DataHandler) error {
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().Unix()
	res.Type = "Response"

	// TODO(shapiro): find a way to tell the client if the subscription errored
	dh.SubscriptionChan <- rabbitmq.Subscription{
		Channel:     cont.key,
		IsSubscribe: true,
	}

	res.ServerMessage = response{
		Status: success,
		Tag:    cont.tag,
		Data:   struct{}{},
	}
	err := toSenderClosure{msg: res}.call(dh) // go ahead and send from here
	return err

}

type rabbitChannelUnsubscribeClosure struct {
	key string
	tag int64
}

func (cont rabbitChannelUnsubscribeClosure) call(dh DataHandler) error {
	res := new(serverMessageWrapper)
	res.Timestamp = time.Now().Unix()
	res.Type = "Response"

	// TODO(shapiro): find a way to tell the client if the subscription errored
	dh.SubscriptionChan <- rabbitmq.Subscription{
		Channel:     cont.key,
		IsSubscribe: false,
	}

	res.ServerMessage = response{
		Status: success,
		Tag:    cont.tag,
		Data:   struct{}{},
	}
	err := toSenderClosure{msg: res}.call(dh) // go ahead and send from here
	return err
}
