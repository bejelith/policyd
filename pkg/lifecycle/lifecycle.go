package lifecycle

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "log/slog"
)

type Managed interface {
	Start()
	Stop()
}

type Lifecycle interface {
	Managed
	Register(Managed)
	Wait()
}

type LifecycleService struct {
	objects []Managed
	signals chan os.Signal
	started bool
	done    chan interface{}
	lock    sync.Mutex
}

func New() *LifecycleService {
	l := &LifecycleService{
		[]Managed{},
		make(chan os.Signal, 1),
		false,
		make(chan interface{}),
		sync.Mutex{},
	}
	return l
}

func (l *LifecycleService) Register(object Managed) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.started {
		return
	}
	l.objects = append(l.objects, object)
}

func (l *LifecycleService) Start() {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.started {
		return
	}
	l.started = true
	signal.Notify(l.signals,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for _, obj := range l.objects {
		obj.Start()
	}
	go l.handleSignals()
}

func (l *LifecycleService) Wait() {
	<-l.done
}

func (l *LifecycleService) Stop() {
	l.lock.Lock()
	defer l.lock.Unlock()
	for _, o := range l.objects {
		o.Stop()
	}
	close(l.done)
}

func (l *LifecycleService) handleSignals() {
	s := <-l.signals
	log.Info("Signal received. Terminating\n", "signal", s)
	l.Stop()
}
