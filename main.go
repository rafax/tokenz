package main

import (
	"github.com/rafax/tokenz/server"
	"github.com/rafax/tokenz/token"
)

var (
	bindTo string = ":8080"
)

func main() {
	handlers := map[string]token.Handler{"b64": token.NewBase64Handler(), "mem": token.NewMemoryHandler()}
	s := server.NewServer(bindTo, handlers)
	s.Start()
}
