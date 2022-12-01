package server

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:embed www/index.html
var homePage []byte

func (this Server) handleRoot(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(homePage))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ctxKey string

const (
	connkey ctxKey = ctxKey("conn")
)

func (this *Server) subWs() {
	this.wsSubId = this.Evt.Sub(ResReady, func(ctx context.Context, msg any) {
		conn := ctx.Value(connkey).(*websocket.Conn)

		// the event system uses goroutines; the docs say apps have locking responsibility
		this.wsWLock.Lock()
		err := conn.WriteJSON(msg)
		this.wsWLock.Unlock()

		if err != nil {
			// errors when it can't get NextWriter, Encode, or Close
			fmt.Printf("Cannot write JSON response") // TODO: log error
		}
	})
}

func (this *Server) handleWs(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		fmt.Printf("Cannot upgrade WebSocket connection") // TODO: log error
		return
	}
	defer conn.Close()
	upgrader.CheckOrigin = this.checkOrigin

	ctx := context.WithValue(req.Context(), connkey, conn)

	for {
		var data *any

		err := conn.ReadJSON(data) // there only ever one reader, no need to lock
		if err != nil {
			fmt.Printf("WS read error") // TODO log error
			break
		}

		this.Evt.Pub(ctx, WsRecv, *data)
	}
}
