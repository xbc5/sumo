package server

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
	if !this.checkOrigin {
		this.Server.checkOrigin = CheckOriginStub
	}
	return this.Server
}
