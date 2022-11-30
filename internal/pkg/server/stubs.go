package server

import (
	"github.com/gorilla/websocket"
)

func CheckOriginStub(conn *websocket.Conn) bool {
	return true
}
