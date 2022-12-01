package server

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:embed www/index.html
var homePage string

func (this Server) handleRoot(res http.ResponseWriter, req *http.Request) {
	fmt.Printf(homePage)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (this Server) handleWs(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	upgrader.CheckOrigin = this.checkOrigin
	if err != nil {
		fmt.Printf("WebSocket connection error") // TODO: log error
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("WS read error") // TODO log error
			break
		}
		fmt.Printf(string(p))
	}
}
