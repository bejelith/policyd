//go:build linux

package acceptor

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
)

type Handler interface {
	Handle(net.Conn)
}
type Managed interface {
	Start()
	Stop()
}

type Acceptor struct {
	goroutines *sync.WaitGroup
	isRunning  *atomic.Bool
	handler    Handler
	listenFd   int
	epollFd    int
}

func New(hander Handler, address string, port int) (*Acceptor, error) {
	listenFd, _ := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	syscall.SetNonblock(listenFd, true)

	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP(address).To4())

	if binderr := syscall.Bind(listenFd, &addr); binderr != nil {
		return nil, binderr
	}
	if listenerr := syscall.Listen(listenFd, 10); listenerr != nil {
		return nil, listenerr
	}

	epollfd, _ := syscall.EpollCreate1(0)

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(listenFd)
	syscall.EpollCtl(epollfd, syscall.EPOLL_CTL_ADD, listenFd, &event)
	started := &atomic.Bool{}
	started.Store(false)
	return &Acceptor{
		goroutines: &sync.WaitGroup{},
		isRunning:  started,
		listenFd:   listenFd,
		epollFd:    epollfd,
		handler:    hander,
	}, nil
}

func (a *Acceptor) accept() {
	a.goroutines.Add(1)
	defer syscall.Close(a.epollFd)
	defer syscall.Close(a.listenFd)
	fmt.Println("for1")
	for a.isRunning.Load() {
		var events = make([]syscall.EpollEvent, 1024)
		n, err := syscall.EpollWait(a.epollFd, events, 1)
		if err != nil {
			fmt.Println("EpollWait: " + err.Error())
		} else {
			for i := 0; i < n; i++ {
				event := events[i]
				if event.Fd == int32(a.listenFd) {
					newFd, _, _ := syscall.Accept(a.listenFd)
					conn, err := fdToConn(newFd)
					if err != nil {
						fmt.Println("Accept()" + err.Error())
					}
					a.handler.Handle(conn)
				}
			}
		}
		runtime.Gosched()
	}
	fmt.Println("acceptor terminating")
	a.goroutines.Done()
}

func (a *Acceptor) Start() {
	a.isRunning.Store(true)
	go a.accept()
}

func (a *Acceptor) Stop() {
	a.isRunning.Store(false)
	a.goroutines.Wait()
	fmt.Println("terminated")
}

func fdToConn(fd int) (net.Conn, error) {
	f := os.NewFile(uintptr(fd), "")
	defer f.Close()
	return net.FileConn(f)
}
