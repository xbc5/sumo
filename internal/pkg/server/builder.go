package server

import "net/http"

type ServerBuilder struct {
	Server      *Server
	checkOrigin bool
	handleRoot  bool
}

func (this *ServerBuilder) CheckOrigin(fn TCheckOrigin) *ServerBuilder {
	this.Server.checkOrigin = fn
	this.checkOrigin = true
	return this
}

func (this *ServerBuilder) HandleRoot(fn http.HandlerFunc) *ServerBuilder {
	this.Server.handleRoot = fn
	this.handleRoot = true
	return this
}

func (this ServerBuilder) Build() *Server {
	if !this.checkOrigin {
		this.Server.checkOrigin = CheckOriginStub
	}
	if !this.handleRoot {
		this.Server.handleRoot = HandleRootStub
	}
	return this.Server
}
