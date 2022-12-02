package server

import "net/http"

func CheckOrigin(req *http.Request) bool {
	return true
}
