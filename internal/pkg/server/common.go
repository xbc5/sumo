package server

import (
	"fmt"
	"net/http"

	"github.com/xbc5/sumo/internal/pkg/config"
)

type Server struct {
	CheckOrigin TCheckOrigin
	Config      config.Config
}

func (this Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", this.handleRoot)
	mux.HandleFunc("/ws", this.handleWs)

	err := http.ListenAndServe(this.Config.Server.Address, mux)
	if err != nil {
		fmt.Errorf("Cannot start server: %s", err)
	}
}

func (this *Server) NewTest() ServerBuilder {
	return ServerBuilder{Server: this}
}
