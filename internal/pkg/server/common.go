package server

import (
	"fmt"
	"net/http"
)

type Server struct{}

func (this Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", this.handleRoot)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Errorf("Cannot start server: %s", err)
	}
}