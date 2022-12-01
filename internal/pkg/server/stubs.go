package server

import (
	"net/http"
)

func CheckOriginStub(req *http.Request) bool {
	return true
}
