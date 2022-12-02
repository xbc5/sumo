package main

import (
	"github.com/xbc5/sumo/internal/pkg/server"
)

func main() {
	s := server.Server{}
	s.New().Start()
}
