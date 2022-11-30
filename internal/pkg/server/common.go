package server

import (
	"fmt"
	"net/http"
)

type Server struct{}

func (this Server) Start() {
	http.HandleFunc("/", this.handleRoot)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Errorf("Cannot start server: %s", err)
	}
}
