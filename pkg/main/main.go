package main

import (
	"flag"

	"github.com/policyd/pkg/acceptor"
	"github.com/policyd/pkg/handler"
	"github.com/policyd/pkg/lifecycle"
)

var addr = flag.String("addr", "0.0.0.0", "Listen address")
var port = flag.Int("port", 12345, "Listen port")

func main() {
	flag.Parse()
	lifecycle := lifecycle.New()
	socketHandler := handler.New()
	acceptor, _ := acceptor.New(socketHandler, *addr, *port)
	lifecycle.Manage(acceptor)
	lifecycle.Manage(socketHandler)
	lifecycle.Start()
	lifecycle.Wait()
}
