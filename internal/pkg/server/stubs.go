package server

import (
	"context"
	"net/http"

	"github.com/xbc5/sumo/internal/pkg/event"
)

func CheckOriginStub(req *http.Request) bool {
	return true
}

type JsonStub struct {
	Kind string
}

func NewOkEvtStub() event.IEvt[any] {
	evt := &event.Evt[any]{}

	// Short circuit the request and response events:
	// when it receives a request, immediately push it to
	// response handler.
	evt.Sub(WsRecv, func(ctx context.Context, msg any) {
		evt.Pub(ctx, ResReady, msg)
	})

	return evt
}
