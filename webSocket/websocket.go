package webSocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
)

type WsJsonResponse struct {
	Message string `json:"message"`
}

type WsConnection struct {
	*websocket.Conn
}

type WsPayload struct {
	Pnu     string       `json:"pnu"`
	Code    string       `json:"code"`
	Mail    string       `json:"mail"`
	Message string       `json:"message"`
	Conn    WsConnection `json:"-"`
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	wsChan = make(chan WsPayload)

	clients = make(map[WsConnection][]string)
	/*
		define clients
		map[WsConnection][]string = {
			*conn: []string{
				"PNU"
				"CODE"
			}
		}
	*/
)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Websocket Closed")
		ws.Close()
	}
	fmt.Println("OK Client Connecting")
	conn := WsConnection{ws}
	clients[conn] = []string{"", ""}
	response := WsJsonResponse{
		Message: "Connection Establish",
	}
	err = ws.WriteJSON(response)

	go ListenForWs(&conn)

}

func ListenForWs(conn *WsConnection) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error", fmt.Sprintf("%v", r))
			conn.Close()
			delete(clients, *conn)
			fmt.Println("Connection Closed.")
		}
	}()
	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err == nil {

			payload.Conn = *conn
			if payload.Code == "" {
				var result string
				for i := 0; i < 10; i++ {
					index := rand.Intn(len(charset))
					result += string(charset[index])
				}
				response := WsJsonResponse{
					Message: "CODE=" + result,
				}
				clients[*conn] = []string{payload.Pnu, result}
				conn.WriteJSON(response)
				payload.Code = result
			} else if clients[*conn][0] == "" && clients[*conn][1] == "" {
				clients[*conn] = []string{payload.Pnu, payload.Code}
			}
			fmt.Println("Connection Receive Pnu:", payload.Pnu, " Code:", payload.Code, " Mail:", payload.Mail, " Message:", payload.Message)
			wsChan <- payload

		}
	}
}

func ListenToWsChan() {
	var response WsJsonResponse
	for {
		e := <-wsChan
		response.Message = e.Message
		broadcastToClient(response, e.Pnu, e.Code, e.Conn)
	}
}

func broadcastToClient(response WsJsonResponse, pnu string, code string, conn WsConnection) {
	for client, clientInfo := range clients {
		if clientInfo[0] == pnu && clientInfo[1] == code && client != conn {
			err := client.WriteJSON(&response)
			if err != nil {
				_ = client.Close()
				fmt.Println("Connection Closed.")
				delete(clients, client)
			}
		}
	}
}
