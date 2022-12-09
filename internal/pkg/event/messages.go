package event

import (
	"context"

	"github.com/google/uuid"
)

type Err struct {
	UUID uuid.UUID `json:"uuid"`
	Msg  string    `json:"msg"`
	Kind string    `json:"-"`
	Err  error     `json:"-"`
	Data any       `json:"data"`
}

type Msg struct {
	UUID uuid.UUID       `json:"uuid"`
	Ctx  context.Context `json:"-"` // not for public consumption
	Data any             `json:"data"`
}

type Sys struct {
	Cmd string
}
