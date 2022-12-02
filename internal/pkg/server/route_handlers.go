package server

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"sync"

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
	wlock   ctxKey = ctxKey("writeMu")
)

func (this *Server) subWs() {
	this.wsSubId = this.Evt.Sub(ResReady, func(ctx context.Context, msg any) {
		conn := ctx.Value(connkey).(*websocket.Conn)
		mu := ctx.Value(wlock).(*sync.RWMutex)

		// the event system uses goroutines; there may be several concurrent pub events
		// calling this code. The Gorilla docs state that apps have locking responsibility.
		// This mutex is scoped to one socket via ctx.
		mu.Lock()
		err := conn.WriteJSON(msg)
		mu.Unlock()

		if err != nil {
			// errors when it can't get NextWriter, Encode, or Close
			fmt.Printf("Cannot write JSON response\n") // TODO: log error
		}
	})
}

// Open a new socket. It emits the raw received messages via the server.WsRecv event.
// You must do type assertions, and validate the data. It passes the Context object
// into the published event.
func (this *Server) handleWs(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		fmt.Printf("Cannot upgrade WebSocket connection") // TODO: log error
		return
	}
	upgrader.CheckOrigin = this.checkOrigin

	ctx := context.WithValue(req.Context(), connkey, conn)

	// this handler initites a new ws; only one writer per connection; we may have multiple
	// connections, so scope the mutex to this connection
	var mu sync.RWMutex
	ctx = context.WithValue(ctx, wlock, &mu)

	for {
		var data interface{}

		err := conn.ReadJSON(&data) // there only ever one reader, no need to lock
		if err != nil {
			// close messages are errors
			break
		}

		if data == nil {
			fmt.Printf("Nil data on socket\n") // TODO log error
			continue
		}

		this.Evt.Pub(ctx, WsRecv, data)
	}
}
