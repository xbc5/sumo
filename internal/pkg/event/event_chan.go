package event

type EvtChans struct {
	NewReq chan Msg
	NewRes chan Msg
	Err    chan Err
	Sys    chan Sys
}

func (this *EvtChans) New() *EvtChans {
	this.NewReq = make(chan Msg)
	this.NewRes = make(chan Msg)
	this.Err = make(chan Err)
	return this
}
