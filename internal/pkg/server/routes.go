package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func (this Server) handleRoot(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("root handled")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func checkOrigin(req *http.Request) bool {
	return true
}

func reader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("WS read error") // TODO log error
		}
		fmt.Printf(string(p))
	}
}

func (this Server) handleWs(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	upgrader.CheckOrigin = checkOrigin
	if err != nil {
		fmt.Printf("WebSocket connection error") // TODO: log error
	}
	reader(conn)
}
