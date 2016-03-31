package handlers

import (
	"fmt"
	"github.com/CodeCollaborate/Server/modules/datahandling"
	"github.com/CodeCollaborate/Server/modules/rabbitmq"
	"github.com/CodeCollaborate/Server/utils"
	"github.com/gorilla/websocket"
	"net/http"
	"sync/atomic"
)

var atomicIdCounter uint64 = 0
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

/**
 * RabbitMq manager for CodeCollaborate Server.
 * @author: Austin Fahsl and Benedict Wong
 */

/**
 * Create a new WebSocket connection given a http request.
 */
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
		fmt.Println("Failed to upgrade connection: %s\n", err)
		return
	}
	defer wsConn.Close()

	// Generate unique ID for this websocket
	var wsId uint64 = atomic.AddUint64(&atomicIdCounter, 1)

	// Run WSSendingHandler in a separate GoRoutine
	go WSSendingRoutine(wsId, wsConn)

	for {
		// messageType, message, err := wsConn.ReadMessage()
		messageType, message, err := wsConn.ReadMessage()
		if err != nil {
			fmt.Println("WebSocket failed to read message, exiting handler\n")
			break
		}
		dh := datahandling.DataHandler{}
		dh.Handle(wsId, messageType, message)
	}
}

/**
 * Receives messages from the RabbitMq subscriber and passes them to the WebSocket.
 */
func WSSendingRoutine(wsId uint64, wsConn *websocket.Conn) {

	ch, messages, err := rabbitmq.RunSubscriber(
		rabbitmq.QueueConfig{
			ExchangeName: "CodeCollaborate",
			QueueId:      wsId,
			Keys:         []string{},
			IsWorkQueue:  false,
		},
	)
	if err != nil {
		utils.LogOnError(err, "Failed to subscribe")
		return
	}
	defer ch.Close()

	for message := range messages {
		wsConn.WriteMessage(websocket.TextMessage, message.Body)
	}
}
