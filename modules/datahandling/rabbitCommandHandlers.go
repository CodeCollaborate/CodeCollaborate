package datahandling

import (
	"github.com/gorilla/websocket"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"encoding/json"
	"github.com/CodeCollaborate/Server/utils"
	"errors"
	"github.com/Sirupsen/logrus"
)

type RabbitCommandHandler struct {
	WSConn       *websocket.Conn
	WSID         uint64
	ExchangeName string
}

func (r RabbitCommandHandler) HandleCommand(msg rabbitmq.AMQPMessage) error {
	var cmd rabbitCommand

	err := json.Unmarshal(msg.Message, &cmd)
	if err != nil {
		return err
	}

	switch(cmd.Command){
	case "Subscribe":
		return r.handleSubscribe(cmd)
	case "Unsubscribe":
		return r.handleUnsubscribe(cmd)
	default:
		err := errors.New("Invalid rabbit command given")
		utils.LogError("Invalid rabbit command given", err, logrus.Fields{
			"RabbitCommand": cmd.Command,
		})
		return err
	}
}

func (r RabbitCommandHandler) handleSubscribe(cmd rabbitCommand) error {
	var data rabbitQueueData
	err := json.Unmarshal(cmd.Data, &data)
	if err != nil {
		return err
	}

	ch, err := rabbitmq.GetChannel()
	if err != nil {
		return err
	}

	var msg *serverMessageWrapper = newEmptyResponse(success, cmd.Tag)
	err = rabbitmq.BindQueue(ch, rabbitmq.RabbitWebsocketQueueName(r.WSID), data.Key, r.ExchangeName)
	if err != nil {
		msg = newEmptyResponse(fail, cmd.Tag)
	}

	// If no tag, do not send a response
	// This is used in cases where we auto-register a client (ie, for username)
	if cmd.Tag < 0 {
		return nil
	}

	// Send response
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.WSConn.WriteMessage(websocket.TextMessage, msgJSON)
}

func (r RabbitCommandHandler) handleUnsubscribe(cmd rabbitCommand) error {
	var data rabbitQueueData
	err := json.Unmarshal(cmd.Data, &data)
	if err != nil {
		return err
	}

	ch, err := rabbitmq.GetChannel()
	if err != nil {
		return err
	}

	var msg *serverMessageWrapper = newEmptyResponse(success, cmd.Tag)
	err = rabbitmq.UnbindQueue(ch, rabbitmq.RabbitWebsocketQueueName(r.WSID), data.Key, r.ExchangeName)
	if err != nil {
		msg = newEmptyResponse(fail, cmd.Tag)
	}

	// If no tag, do not send a response
	// This is used in cases where we auto-register a client (ie, for username)
	if cmd.Tag < 0 {
		return nil
	}

	// Send response
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.WSConn.WriteMessage(websocket.TextMessage, msgJSON)
}