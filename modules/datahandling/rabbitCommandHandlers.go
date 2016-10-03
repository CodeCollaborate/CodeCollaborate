package datahandling

import (
	"encoding/json"
	"errors"

	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/gorilla/websocket"
)

// RabbitCommandHandler handles all rabbit commands (sub/unsub)
type RabbitCommandHandler struct {
	WSConn       *websocket.Conn
	WSID         uint64
	ExchangeName string
}

// HandleCommand handles an individual command
func (r RabbitCommandHandler) HandleCommand(msg rabbitmq.AMQPMessage) error {
	var cmd rabbitCommand

	err := json.Unmarshal(msg.Message, &cmd)
	if err != nil {
		return err
	}

	switch cmd.Command {
	case "Subscribe":
		return r.handleSubscribe(cmd)
	case "Unsubscribe":
		return r.handleUnsubscribe(cmd)
	default:
		err := errors.New("Invalid rabbit command given")
		utils.LogError("Invalid rabbit command given", err, utils.LogFields{
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

	msg := newEmptyResponse(success, cmd.Tag)
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

	msg := newEmptyResponse(success, cmd.Tag)
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
