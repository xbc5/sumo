package server

type ServerBuilder struct {
	Server      *Server
	checkOrigin bool
}

func (this *ServerBuilder) CheckOrigin(fn TCheckOrigin) *ServerBuilder {
	this.Server.CheckOrigin = fn
	this.checkOrigin = true
	return this
}

func (this ServerBuilder) Build() *Server {
	if !this.checkOrigin {
		this.Server.CheckOrigin = CheckOriginStub
	}
	return this.Server
}
