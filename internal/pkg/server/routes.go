package server

import (
	"fmt"
	"net/http"
)

func (this Server) handleRoot(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("root handled")
}
