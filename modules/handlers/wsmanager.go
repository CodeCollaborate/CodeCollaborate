package handlers

import (
	"errors"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datahandling"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/gorilla/websocket"
	"github.com/kr/pretty"
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

	// TODO: Send data blob

	// Generate unique ID for this websocket
	wsID := atomic.AddUint64(&atomicIDCounter, 1)

	pubCfg := rabbitmq.NewPubConfig(func(msg rabbitmq.AMQPMessage) {
		// TODO(wongb): Do we need to send errors back to the client on publishing fail? Can we just kill the socket?
		msg.ErrHandler()
	}, outboundMessageQueueBufferSize)

	subCfg := &rabbitmq.AMQPSubCfg{
		QueueID:     wsID,
		Keys:        []string{},
		IsWorkQueue: false,
	}

	pubSubCfg := rabbitmq.NewAMQPPubSubCfg(cfg.ServerConfig.Name, pubCfg, subCfg)

	subCfg.HandleMessageFunc = newAMQPMessageHandler(wsID, pubSubCfg, wsConn)

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

	// we don't actually need more than 1 datahandler per websocket
	dh := datahandling.DataHandler{
		MessageChan: pubCfg.Messages,
		WebsocketID: wsID,
		Db:          dbfs.Dbfs,
	}

	// Waitgroup to make sure channel is closed at appropriate time.
	dhCompleted := &sync.WaitGroup{}

loop:
	for {
		select {
		case <-pubSubCfg.Control.Exit:
			break loop
		default:
			messageType, message, err := wsConn.ReadMessage()
			if err != nil {
				utils.LogError("Failed to read message, terminating connection", err, nil)
				pubSubCfg.Control.Shutdown()
				break loop
			}

			dhCompleted.Add(1)
			go dh.Handle(messageType, message, dhCompleted)
		}
	}

	// Wait for all datahandlers to complete before closing channel
	dhCompleted.Wait()
	close(pubCfg.Messages)
}

func newAMQPMessageHandler(websocketID uint64, cfg *rabbitmq.AMQPPubSubCfg, wsConn *websocket.Conn) func(rabbitmq.AMQPMessage) error {
	queueName := rabbitmq.RabbitWebsocketQueueName(websocketID)

	return func(msg rabbitmq.AMQPMessage) error {
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
				WSID:         cfg.SubCfg.QueueID,
			}
			return rch.HandleCommand(msg)
		default:
			err := errors.New("No such ContentType")
			utils.LogError("Invalid ContentType", err, utils.LogFields{
				"AMQPMessage": pretty.Sprint(msg),
			})
			return err
		}
	}
}
