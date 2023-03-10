package lifecycle

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Lifecycle interface {
	Manage(Managed)
	Start()
	Wait()
}

type lifecycle struct {
	objects []Managed
	signals chan os.Signal
	started bool
	done    chan interface{}
	lock    sync.Mutex
}

func New() Lifecycle {
	l := &lifecycle{
		[]Managed{},
		make(chan os.Signal, 1),
		false,
		make(chan interface{}),
		sync.Mutex{},
	}
	return l
}

func (l *lifecycle) Manage(object Managed) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.started {
		return
	}
	l.objects = append(l.objects, object)
}

func (l *lifecycle) Start() {
	if l.started {
		return
	}
	signal.Notify(l.signals,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	l.started = true
	go l.manage()
}

func (l *lifecycle) Wait() {
	<-l.done
}

func (l *lifecycle) manage() {
	l.lock.Lock()
	defer l.lock.Unlock()
	for _, obj := range l.objects {
		obj.Start()
	}
	s := <-l.signals
	fmt.Printf("Signal %s received. Terminating\n", s)
	for _, obj := range l.objects {
		obj.Stop()
	}
	close(l.done)
}
