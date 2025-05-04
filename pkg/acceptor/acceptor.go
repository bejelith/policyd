package acceptor

import (
	log "log/slog"
	"net"
	"sync"
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
	wg *sync.WaitGroup
	handler    Handler
	listener   net.Listener
	loops      int
	started    bool
	stop       chan struct{}
}

func New(handler Handler, listener net.Listener) (*Acceptor, error) {
	return &Acceptor{
		wg: &sync.WaitGroup{},
		listener:   listener,
		handler:    handler,
		stop:       make(chan struct{}),
		loops:      defaultLoops,
	}, nil
}

func (a *Acceptor) accept() {
	for range a.loops {
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			for {
				select {
				case <-a.stop:
					return
				default:
					conn, err := a.listener.Accept()
					if err != nil {
						log.Error("accept", "error", err)
						return
					}
					go a.handler.Handle(conn)
				}
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
	close(a.stop)
	a.listener.Close()
	a.wg.Wait()
}
