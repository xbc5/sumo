package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	CheckOrigin TCheckOrigin
}

func (this Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", this.handleRoot)
	mux.HandleFunc("/ws", this.handleWs)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Errorf("Cannot start server: %s", err)
	}
}

func (this *Server) NewTest() ServerBuilder {
	return ServerBuilder{Server: this}
}
