package main

import (
	"flag"

	"githun.com/policyd/pkg/acceptor"
	"githun.com/policyd/pkg/handler"
	"githun.com/policyd/pkg/lifecycle"
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
