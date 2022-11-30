package server

import "github.com/gorilla/websocket"

type TCheckOrigin func(conn *websocket.Conn) bool
