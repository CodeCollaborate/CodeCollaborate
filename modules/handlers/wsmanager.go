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
	"errors"
	"github.com/kr/pretty"
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

	pubCfg := rabbitmq.NewPubConfig(func(msg rabbitmq.AMQPMessage) {
		// TODO(wongb): Do we need to send errors back to the client on publishing fail? Can we just kill the socket?
		msg.ErrHandler()
	})

	subCfg := &rabbitmq.AMQPSubCfg{
		QueueID:      wsID,
		Keys:         []string{},
		IsWorkQueue:  false,
	}

	pubSubCfg := rabbitmq.NewAMQPPubSubCfg(cfg.ServerConfig.Name, pubCfg, subCfg)

	subCfg.HandleMessageFunc = newAMQPMessageHandler(pubSubCfg, wsConn);

	//defer func() {
	//	// this prevents a channel leak on an unplanned exit
	//	sendingRoutineControl.Exit <- true
	//	pubCfg.Control.Exit <- true
	//	// we want to recover here so that the server doesn't die
	//	if r := recover(); r != nil {
	//		// TODO(shapiro): Make sure this gets properly logged.
	//		// the most likely cause is that we tried to close an already closed channel
	//		utils.LogError("Recovered from WSManager panic", nil, nil)
	//	}
	//}()

	go func() {
		err := rabbitmq.RunPublisher(pubSubCfg)
		if err != nil {
			utils.LogError("Publisher error encountered. Exiting", err, nil)
			close(pubSubCfg.Control.Exit)
		}
	}()
	go func() {
		err := rabbitmq.RunSubscriber(pubSubCfg)
		if err != nil {
			utils.LogError("Subscriber error encountered. Exiting", err, nil)
			close(pubSubCfg.Control.Exit)
		}
	}()

	// we don't actually need more than 1 datahandler per websocket
	dh := datahandling.DataHandler{
		MessageChan:      pubCfg.Messages,
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

func newAMQPMessageHandler(cfg *rabbitmq.AMQPPubSubCfg, wsConn *websocket.Conn) (func(rabbitmq.AMQPMessage) error) {
	return func(msg rabbitmq.AMQPMessage) error {
		switch (msg.ContentType){
		case rabbitmq.ContentType_Msg:
			utils.LogDebug("Sending Message", logrus.Fields{
				"Message": string(msg.Message),
			})
			return wsConn.WriteMessage(websocket.TextMessage, msg.Message)
		case rabbitmq.ContentType_Cmd:
			rch := datahandling.RabbitCommandHandler{
				ExchangeName: cfg.ExchangeName,
				WSConn: wsConn,
				WSID: cfg.SubCfg.QueueID,
			}
			return rch.HandleCommand(msg)
		default:
			err := errors.New("No such ContentType")
			utils.LogError("Invalid ContentType", err, logrus.Fields{
				"AMQPMessage": pretty.Sprint(msg),
			})
			return err
		}
	}
}
