package rabbitmq

import (
	"encoding/json"
	"errors"

	"github.com/CodeCollaborate/Server/modules/messages"
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
func (r RabbitCommandHandler) HandleCommand(msg AMQPMessage) error {
	var cmd RabbitCommandJSON

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

func (r RabbitCommandHandler) handleSubscribe(cmd RabbitCommandJSON) error {
	var data RabbitQueueData
	err := json.Unmarshal(cmd.Data, &data)
	if err != nil {
		return err
	}

	ch, err := GetChannel()
	if err != nil {
		return err
	}

	msg := messages.NewEmptyResponse(messages.StatusSuccess, cmd.Tag)
	err = BindQueue(ch, RabbitWebsocketQueueName(r.WSID), data.Key, r.ExchangeName)
	if err != nil {
		msg = messages.NewEmptyResponse(messages.StatusFail, cmd.Tag)
	}

	// If no tag, do not send a response
	// This is used in cases where we auto-register a client (ie, for username)
	if cmd.Tag < 0 || r.WSConn == nil {
		return nil
	}

	// Send response
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.WSConn.WriteMessage(websocket.TextMessage, msgJSON)
}

func (r RabbitCommandHandler) handleUnsubscribe(cmd RabbitCommandJSON) error {
	var data RabbitQueueData
	err := json.Unmarshal(cmd.Data, &data)
	if err != nil {
		return err
	}

	ch, err := GetChannel()
	if err != nil {
		return err
	}

	msg := messages.NewEmptyResponse(messages.StatusSuccess, cmd.Tag)
	err = UnbindQueue(ch, RabbitWebsocketQueueName(r.WSID), data.Key, r.ExchangeName)
	if err != nil {
		msg = messages.NewEmptyResponse(messages.StatusFail, cmd.Tag)
	}

	// If no tag, do not send a response
	// This is used in cases where we auto-register a client (ie, for username)
	if cmd.Tag < 0 || r.WSConn == nil {
		return nil
	}

	// Send response
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.WSConn.WriteMessage(websocket.TextMessage, msgJSON)
}
