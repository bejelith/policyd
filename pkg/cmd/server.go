package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	log "log/slog"

	"github.com/policyd/pkg/acceptor"
	"github.com/policyd/pkg/handler"
	"github.com/policyd/pkg/lifecycle"
)

var addr = flag.String("addr", "0.0.0.0:12345", "Listen address")
var logLevel = log.LevelInfo

func main() {
	flag.Func("loglevel", "logevel value", setLevel)
	flag.Parse()
	log.SetDefault(log.New(log.NewTextHandler(os.Stdout, &log.HandlerOptions{
		Level: logLevel,
		AddSource: true,
	})))

	//plugin.Register(&quota.Plugin{})

	lc := lifecycle.New()
	connHandler := handler.NewConnHandler()
	listener, err := setupListener(*addr)
	if err != nil {
		log.Error("failed to open port at address", "address", *addr, "error", err)
	}
	defer listener.Close()
	acceptor, _ := acceptor.New(connHandler, listener)
	lc.Register(acceptor)
	lc.Register(connHandler)

	lc.Start()
	lc.Wait()
}

func setLevel(s string) error {
	var level log.Level = log.LevelInfo
	err := level.UnmarshalJSON([]byte(s))
	if err != nil {
		logLevel = level
	}
	return err
}

func setupListener(addr string) (net.Listener, error) {
	l := net.ListenConfig{}
	listener, err := l.Listen(context.TODO(), "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to open listener: %v", err)
	}
	return listener, nil
}
