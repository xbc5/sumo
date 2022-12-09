package event

import (
	"net/http"
)

type HTTPErr struct {
	Req *http.Request `json:"-"`
}
