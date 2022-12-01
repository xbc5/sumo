package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/xbc5/sumo/internal/pkg/config"
	"github.com/xbc5/sumo/internal/pkg/event"
)

type Server struct {
	checkOrigin TCheckOrigin
	Config      config.Server
	Evt         event.IEvt[any]

	// the ID used to unsubscribe the WebSocket listener
	wsSubId int
}

func (this *Server) createHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", this.handleRoot)
	mux.HandleFunc("/ws", this.handleWs)
	return mux
}

func (this *Server) Start() {
	addr := this.Config.Address
	fmt.Printf("Starting server on: %q\n", addr) // TODO log info
	err := http.ListenAndServe(addr, this.createHandler())
	fmt.Printf("Err %s\n", err) // TODO log info
	if err != nil {
		fmt.Errorf("Cannot start server: %s", err) // TODO log error
	}
}

func (this Server) Stop() {
	// no need to do conn.Close() here, ws handlers use defer conn.Close()
	this.Evt.Unsub(this.wsSubId) // errors is "id not found", so what?
}

func (this *Server) StartTest() *httptest.Server {
	return httptest.NewServer(this.createHandler())
}

func (this *Server) New() *Server {
	this.Config = config.GetConfig().Server
	this.subWs() // make responses to WebSocket
	build := ServerBuilder{Server: this}
	server := build.CheckOrigin(CheckOrigin).Build()
	return server
}

func (this *Server) NewTest() ServerBuilder {
	this.Config = config.GetConfig().Server
	return ServerBuilder{Server: this}
}
