package handlers

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"github.com/CodeCollaborate/CodeCollaborate/modules/datahandling"
	"sync/atomic"
	"github.com/CodeCollaborate/CodeCollaborate/modules/rabbitmq"
	"github.com/CodeCollaborate/CodeCollaborate/utils"
)

var atomicIdCounter uint64 = 0
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options


func NewWSConn(responseWriter http.ResponseWriter, request *http.Request) {
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
	var wsId uint64 = atomic.AddUint64(&atomicIdCounter, 1)
	go WSSendingRoutine(wsId, wsConn)

	defer wsConn.Close()
	// defer managers.WebSocketDisconnected(wsConn)

	for {
		// messageType, message, err := wsConn.ReadMessage()
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			fmt.Println("WebSocket failed to read message, exiting handler\n")
			break
		}
		dh := datahandling.DataHandler{}
		dh.Handle(wsId, message)
	}
}

func WSSendingRoutine(wsId uint64, wsConn *websocket.Conn) {

	ch, messages, err := rabbitmq.RunSubscriber(wsId)
	if err != nil {
		utils.LogOnError(err, "Failed to subscribe")
		return
	}
	defer ch.Close()

	for message := range messages {
		wsConn.WriteMessage(websocket.TextMessage, message.Body)
	}
}