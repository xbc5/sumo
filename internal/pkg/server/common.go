package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/xbc5/sumo/internal/pkg/config"
)

type Server struct {
	checkOrigin TCheckOrigin
	Config      config.Config
}

func (this Server) createHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", this.handleRoot)
	mux.HandleFunc("/ws", this.handleWs)
	return mux
}

func (this Server) Start() {
	err := http.ListenAndServe(this.Config.Server.Address, this.createHandler())
	if err != nil {
		fmt.Errorf("Cannot start server: %s", err) // TODO log error
	}
}

func (this Server) StartTest() *httptest.Server {
	return httptest.NewServer(this.createHandler())
}

func (this *Server) New() *Server {
	build := ServerBuilder{Server: this}
	server := build.CheckOrigin(CheckOrigin).Build()
	return server
}

func (this *Server) NewTest() ServerBuilder {
	return ServerBuilder{Server: this}
}
