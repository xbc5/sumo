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
	Evt         *event.EvtChans
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
	go this.listenForEvents() // blocking
}

func (this *Server) Stop() {
	this.Evt.Sys <- event.Sys{Cmd: event.StopCmd}
}

func (this *Server) StartTest() *httptest.Server {
	go this.listenForEvents() // blocking
	return httptest.NewServer(this.createHandler())
}

func (this *Server) New() *Server {
	this.Config = config.GetConfig().Server
	build := ServerBuilder{Server: this}
	server := build.CheckOrigin(CheckOrigin).Build()
	return server
}

func (this *Server) NewTest() *ServerBuilder {
	this.Config = config.GetConfig().Server
	s := &ServerBuilder{Server: this}
	return s
}
