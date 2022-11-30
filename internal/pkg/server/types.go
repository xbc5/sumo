package server

import (
	"net/http"
)

type TCheckOrigin func(req *http.Request) bool
