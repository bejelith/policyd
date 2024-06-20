package acceptor

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	log "log/slog"
)

var defaultLoops = 2

type Handler interface {
	Handle(net.Conn)
}

type Managed interface {
	Start()
	Stop()
}

type Acceptor struct {
	goroutines *sync.WaitGroup
	handler    Handler
	listener   net.Listener
	loops      int
	started    bool
	stop       *atomic.Bool
}

func New(handler Handler, listener net.Listener) (*Acceptor, error) {
	return &Acceptor{
		goroutines: &sync.WaitGroup{},
		listener:   listener,
		handler:    handler,
		stop:       &atomic.Bool{},
		loops:      defaultLoops,
	}, nil
}

func (a *Acceptor) accept() {
	for i := 0; i < a.loops; i++ {
		a.goroutines.Add(1)
		go func() {
			defer a.goroutines.Done()
			for !a.stop.Load() {
				conn, err := a.listener.Accept()
				if err != nil {
					log.Error("accept", "error", err)
					return
				}
				go a.handler.Handle(conn)
			}
		}()
	}
}

func (a *Acceptor) Start() {
	if a.started {
		return
	}
	a.started = true
	go a.accept()
}

func (a *Acceptor) Stop() {
	a.stop.Store(true)
	a.listener.Close()
	a.goroutines.Wait()
	fmt.Println("terminated")
}
