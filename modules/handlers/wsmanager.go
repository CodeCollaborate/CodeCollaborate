package handlers

import (
	"net/http"
	"sync/atomic"

	"github.com/CodeCollaborate/Server/modules/config"
	"github.com/CodeCollaborate/Server/modules/datahandling"
	"github.com/CodeCollaborate/Server/modules/dbfs"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/gorilla/websocket"
	"github.com/Sirupsen/logrus"
)

/**
 * WSManager handles all WebSocket upgrade requests.
 */

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

	// Run WSSendingHandler in a separate GoRoutine
	sendingRoutineControl := rabbitmq.NewControl()

	// we (probably) want to have 1 publisher per connection to prevent overload. Goroutines are cheap.
	pubCfg := rabbitmq.NewPubConfig(cfg.ServerConfig.Name)

	defer func() {
		// this prevents a channel leak on an unplanned exit
		sendingRoutineControl.Exit <- true
		pubCfg.Control.Exit <- true
		// we want to recover here so that the server doesn't die
		if r := recover(); r != nil {
			// TODO(shapiro): Make sure this gets properly logged.
			// the most likely cause is that we tried to close an already closed channel
			utils.LogError("Recovered from WSManager panic", nil, nil)
		}
	}()

	go rabbitmq.RunPublisher(pubCfg)
	go WSSendingRoutine(wsID, wsConn, sendingRoutineControl)

	// we don't actually need more than 1 datahandler per websocket
	dh := datahandling.DataHandler{
		MessageChan:      pubCfg.Messages,
		SubscriptionChan: sendingRoutineControl.SubChan,
		WebsocketID:      wsID,
		Db:               dbfs.Dbfs,
	}

	for {
		messageType, message, err := wsConn.ReadMessage()
		if err != nil {
			utils.LogError("Failed to read message, terminating connection", err, nil)
			break
		}
		go dh.Handle(messageType, message)
	}
}

// WSSendingRoutine receives messages from the RabbitMq subscriber and passes them to the WebSocket.
func WSSendingRoutine(wsID uint64, wsConn *websocket.Conn, ctrl *rabbitmq.RabbitControl) {
	cfg := config.GetConfig()

	err := rabbitmq.RunSubscriber(
		&rabbitmq.AMQPSubCfg{
			ExchangeName: cfg.ServerConfig.Name,
			QueueID:      wsID,
			Keys:         []string{},
			IsWorkQueue:  false,
			HandleMessageFunc: func(msg rabbitmq.AMQPMessage) error {
				utils.LogDebug("Sending Message", logrus.Fields{
					"Message": string(msg.Message),
				})
				return wsConn.WriteMessage(websocket.TextMessage, msg.Message)
			},
			Control: ctrl,
		},
	)

	// TODO(wongb): Is this really supposed to die if we cannot subscribe?
	utils.LogError("Failed to subscribe to RabbitMQ channel", err, nil)
}
