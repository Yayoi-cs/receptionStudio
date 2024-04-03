package webSocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type WsJsonResponse struct {
	Message string `json:"message"`
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("OK Client Connecting")
	conn := WsConnection{ws}
	clients[conn] = "user1"
	response := WsJsonResponse{
		Message: "Hello World",
	}
	err = ws.WriteJSON(response)
	go ListenForWs(&conn)

}

func ListenForWs(conn *WsConnection) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error", fmt.Sprintf("%v", r))
		}
	}()
	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChan() {
	var response WsJsonResponse
	for {
		e := <-wsChan
		response.Message = e.Message
		broadcastToAll(response)
	}
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(&response.Message)
		if err != nil {
			_ = client.Close()
			delete(clients, client)
		}
	}
}

type WsConnection struct {
	*websocket.Conn
}

type WsPayload struct {
	Message string       `json:"message"`
	Conn    WsConnection `json:"-"`
}

var (
	wsChan = make(chan WsPayload)

	clients = make(map[WsConnection]string)
)
