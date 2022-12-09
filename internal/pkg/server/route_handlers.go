package server

import (
	"context"
	_ "embed"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/xbc5/sumo/internal/pkg/errs"
	"github.com/xbc5/sumo/internal/pkg/event"
)

//go:embed www/index.html
var homePage []byte

func (this Server) handleRoot(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	_, err := res.Write([]byte(homePage))
	if err != nil {
		this.Evt.Err <- event.Err{
			Err:  err,
			Msg:  "Cannot write reponse for /",
			Kind: errs.HTTP,
			Data: event.HTTPErr{Req: req},
		}
	}
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

func (this *Server) respond(msg event.Msg) (err error) {
	conn := msg.Ctx.Value(connkey).(*websocket.Conn)
	mu := msg.Ctx.Value(wlock).(*sync.RWMutex)

	// Gorilla expects you to manage concurrent writes to socket.
	// Each request creates a session. The lifetime of that session
	// relates to one client, and thus we create the mutex when the
	// session starts, and pass it via the context.
	mu.Lock()
	err = conn.WriteJSON(msg)
	mu.Unlock()

	if err != nil {
		// errors when it can't get NextWriter, Encode, or Close
		this.Evt.Err <- event.Err{
			UUID: msg.UUID,
			Err:  err,
			Msg:  "Cannot write JSON response",
			Kind: errs.WebSocket,
		}
	}

	return
}

func (this *Server) listenForEvents() {
	for {
		select {
		case msg := <-this.Evt.NewRes:
			this.respond(msg)
		case msg := <-this.Evt.Sys:
			if msg.Cmd == event.StopCmd {
				return
			}
		}
	}
}

// Open a new socket. It emits the raw received messages via the server.WsRecv event.
// You must do type assertions, and validate the data. It passes the Context object
// into the published event.
func (this *Server) handleWs(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		this.Evt.Err <- event.Err{
			Err:  err,
			Msg:  "Cannot upgrade WebSocket connection",
			Kind: errs.WebSocket,
		}
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
			// "close messages" are treated as errors
			// TODO: log close messages and ordinary errors differently: info and err
			// also emit an error message to the client if (for example) invalid JSON
			break
		}

		if data == nil {
			this.Evt.Err <- event.Err{
				Err:  err,
				Msg:  "Nil data on socket",
				Kind: errs.WebSocket,
			}
			continue
		}

		id, uErr := uuid.NewRandom() // v4; v1 includes MAC (privacy issue)
		if uErr != nil {
			panic("Cannot create UUID")
		}
		msg := event.Msg{
			UUID: id,
			Ctx:  ctx,
			Data: data,
		}
		this.Evt.NewReq <- msg
	}
}
