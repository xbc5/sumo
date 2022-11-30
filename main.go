package main

import (
	"fmt"

	"github.com/xbc5/sumo/internal/pkg/server"
)

func main() {
	s := server.Server{}
	fmt.Printf("Starting")
	s.Start()
}
