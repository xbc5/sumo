package server

import (
	"net/http"
)

func CheckOriginStub(req *http.Request) bool {
	return true
}

func HandleRootStub(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("response stub"))
}
