package server

import "github.com/xbc5/sumo/internal/pkg/event"

type ServerBuilder struct {
	Server      *Server
	checkOrigin bool
}

func (this *ServerBuilder) CheckOrigin(fn TCheckOrigin) *ServerBuilder {
	this.Server.checkOrigin = fn
	this.checkOrigin = true
	return this
}

func (this ServerBuilder) Build() *Server {
	if this.Server.Evt == nil {
		evt := &event.EvtChans{}
		this.Server.Evt = evt.New()
	}

	if !this.checkOrigin {
		this.Server.checkOrigin = CheckOriginStub
	}

	return this.Server
}
