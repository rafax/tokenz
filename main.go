package main

import (
	"github.com/rafax/tokenz/server"
	"github.com/rafax/tokenz/token"
)

var (
	bindTo string = ":8080"
)

func main() {
	s := server.NewServer(token.NewBase64Handler(), token.NewMemoryHandler(), bindTo)
	s.Start()
}
