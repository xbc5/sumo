package server

import "github.com/xbc5/sumo/internal/pkg/event"

type ServerBuilder struct {
	Server      *Server
	checkOrigin bool
	evt         bool
}

func (this *ServerBuilder) CheckOrigin(fn TCheckOrigin) *ServerBuilder {
	this.Server.checkOrigin = fn
	this.checkOrigin = true
	return this
}

func (this *ServerBuilder) Evt(evt event.IEvt[any]) *ServerBuilder {
	this.Server.Evt = evt
	this.evt = true
	return this
}

func (this ServerBuilder) Build() *Server {
	if !this.checkOrigin {
		this.Server.checkOrigin = CheckOriginStub
	}

	if !this.evt {
		this.Server.Evt = NewOkEvtStub()
	}
	this.Server.subWs()
	return this.Server
}
