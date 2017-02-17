package handlers

import (
	"errors"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/kr/pretty"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
)

/**
 * WSManager handles all WebSocket upgrade requests.
 */

const outboundMessageQueueBufferSize = 32

// Counter for unique ID of WebSockets Connections. Unique to hostname.
var atomicIDCounter uint64

// Define WebSocket Upgrader that ignores origin; there is never going to be a referral source.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewWSConn accepts a HTTP Upgrade request, creating a new websocket connection.
// Once a WebSocket connection is created, will setup the Receiving and Sending routines,
// then
func NewWSConn(responseWriter http.ResponseWriter, request *http.Request) {
	// Receive and upgrade request
	if request.URL.Path != "/ws/" {
		http.Error(responseWriter, "Not found", 404)
		return
	}
	if request.Method != "GET" {
		http.Error(responseWriter, "Method not allowed", 405)
		return
	}
	wsConn, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		utils.LogError("Failed to upgrade connection", err, nil)
		return
	}
	defer wsConn.Close()
	cfg := config.GetConfig()

	// Generate unique ID for this websocket
	wsID := atomic.AddUint64(&atomicIDCounter, 1)

	pubCfg := rabbitmq.NewPubConfig(func(msg rabbitmq.AMQPMessage) {
		// TODO(wongb): Do we need to send errors back to the client on publishing fail? Can we just kill the socket?
		msg.ErrHandler()
	}, outboundMessageQueueBufferSize)

	defer close(pubCfg.Messages)

	subCfg := &rabbitmq.AMQPSubCfg{
		QueueName:   rabbitmq.LocalWebsocketName(wsID),
		Keys:        []string{},
		IsWorkQueue: false,
	}

	pubSubCfg := rabbitmq.NewAMQPPubSubCfg(cfg.ServerConfig.Name, pubCfg, subCfg)

	subCfg.HandleMessageFunc = newClientMessageHandler(pubSubCfg, wsConn)

	go func() {
		err := rabbitmq.RunPublisher(pubSubCfg)
		if err != nil {
			utils.LogError("Publisher error encountered. Exiting", err, nil)
			pubSubCfg.Control.Shutdown()
		}
	}()
	go func() {
		err := rabbitmq.RunSubscriber(pubSubCfg)
		if err != nil {
			utils.LogError("Subscriber error encountered. Exiting", err, nil)
			pubSubCfg.Control.Shutdown()
		}
	}()

	pubSubCfg.Control.Ready.Wait()

msgLoop:
	for {
		select {
		case <-pubSubCfg.Control.Exit:
			break msgLoop
		default:
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseMessage, websocket.CloseNoStatusReceived) {
					break msgLoop
				}
				utils.LogError("Failed to read message, terminating connection", err, nil)
				pubSubCfg.Control.Shutdown()
				break msgLoop
			}

			err = WorkerEnqueue(message, wsID)
			if err != nil {
				// FIXME: retry?
			}
		}
	}
}

func newClientMessageHandler(cfg *rabbitmq.AMQPPubSubCfg, wsConn *websocket.Conn) func(rabbitmq.AMQPMessage) error {
	queueName := cfg.SubCfg.QueueName

	return func(msg rabbitmq.AMQPMessage) error {
		err := msg.Ack() // ack early b/c regardless the outcome here we don't want to re-enqueue
		utils.LogError("Error Ack'ing RabbitMQ message", err, utils.LogFields{
			"Message": string(msg.Message),
		})
		// I don't know what to do with this error, it can only happen if we disconnect from rabbit,
		// so it means we have bigger issues
		// FIXME: is this fatal?

		switch msg.ContentType {
		case rabbitmq.ContentTypeMsg:
			// If notification with self as origin, early-out; ignore our own notifications.
			if val, ok := msg.Headers["MessageType"]; ok && val == "Notification" {
				if val, ok := msg.Headers["Origin"]; ok && val == queueName {
					return nil
				}
			}

			utils.LogDebug("Sending Message", utils.LogFields{
				"Message": string(msg.Message),
			})
			return wsConn.WriteMessage(websocket.TextMessage, msg.Message)
		case rabbitmq.ContentTypeCmd:
			rch := rabbitmq.RabbitCommandHandler{
				ExchangeName: cfg.ExchangeName,
				WSConn:       wsConn,
				QueueName:    cfg.SubCfg.QueueName,
			}
			return rch.HandleCommand(msg)
		default:
			err := errors.New("Unnable to process RabbitMQ message type")
			utils.LogError("Invalid ContentType", err, utils.LogFields{
				"AMQPMessage": pretty.Sprint(msg),
			})
			return err
		}
	}
}
