package datahandling

import (
	"encoding/json"
	"os"
	"time"

	"github.com/CodeCollaborate/Server/modules/rabbitmq"
)

var hostname, _ = os.Hostname()

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
		RoutingKey:  rabbitmq.RabbitQueueName(dh.WebsocketID),
		ContentType: cont.msg.Type,
		Persistent:  false,
		Message:     msgJSON,
	}
	return nil
}

type toRabbitChannelClosure struct {
	msg        *serverMessageWrapper
	routingKey int64
}

// toRabbitChannelClosure.call is the function that will forward a server message to a channel based on the given routing key
func (cont toRabbitChannelClosure) call(dh DataHandler) error {
	msgJSON, err := json.Marshal(cont.msg)
	if err != nil {
		return err
	}
	dh.MessageChan <- rabbitmq.AMQPMessage{
		Headers:     make(map[string]interface{}),
		RoutingKey:  rabbitmq.RabbitProjectQueueName(cont.routingKey),
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

	// I (joel) don't believe we actually have a way to know here if this subscribe throws an error
	dh.SubscriptionChan <- rabbitmq.Subscription{
		Channel:     cont.key,
		IsSubscribe: true,
	}

	//if err != nil {
	//	res.ServerMessage = response{
	//		Status: fail,
	//		Tag:    p.Tag,
	//		Data:   struct{}{}}
	//} else {
	res.ServerMessage = response{
		Status: success,
		Tag:    cont.tag,
		Data:   struct{}{}}
	//}
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

	// I (joel) don't believe we actually have a way to know here if this subscribe throws an error
	dh.SubscriptionChan <- rabbitmq.Subscription{
		Channel:     cont.key,
		IsSubscribe: false,
	}

	//if err != nil {
	//	res.ServerMessage = response{
	//		Status: fail,
	//		Tag:    p.Tag,
	//		Data:   struct{}{}}
	//} else {
	res.ServerMessage = response{
		Status: success,
		Tag:    cont.tag,
		Data:   struct{}{}}
	//}
	err := toSenderClosure{msg: res}.call(dh) // go ahead and send from here
	return err
}
