//go:build linux

package acceptor

import (
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"syscall"

	"githun.com/policyd/pkg/handler"
)

type Acceptor interface {
	Start()
	Stop()
	Wait()
}
type acceptor struct {
	done     chan interface{}
	started  atomic.Bool
	handler  handler.SocketHandler
	listenFd int
	epollFd  int
}

func New(hander handler.SocketHandler, address string, port int) (Acceptor, error) {
	listenFd, _ := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	syscall.SetNonblock(listenFd, true)

	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP(address).To4())

	syscall.Bind(listenFd, &addr)
	syscall.Listen(listenFd, 10)

	epollfd, _ := syscall.EpollCreate1(0)

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(listenFd)
	syscall.EpollCtl(epollfd, syscall.EPOLL_CTL_ADD, listenFd, &event)
	started := atomic.Bool{}
	started.Store(false)
	return &acceptor{
		done:     make(chan interface{}),
		started:  started,
		listenFd: listenFd,
		epollFd:  epollfd,
		handler:  hander,
	}, nil
}

func (a *acceptor) accept() {
	if a.started.Load() {
		return
	}
	defer syscall.Close(a.epollFd)
	defer syscall.Close(a.listenFd)

	for {
		if ischannelClosed(a.done) {
			return
		}
		var events = make([]syscall.EpollEvent, 1024)
		n, err := syscall.EpollWait(a.epollFd, events, -1)
		if err != nil {
			fmt.Println("Epoll: " + err.Error())
		} else {
			for i := 0; i < n; i++ {
				event := events[i]
				if event.Fd == int32(a.listenFd) {
					newFd, _, _ := syscall.Accept(a.listenFd)
					conn, err := fdToConn(newFd)
					if err != nil {
						fmt.Println("a" + err.Error())
					}
					a.handler.Handle(conn)
				}
			}
		}
	}
}

func (a *acceptor) Start() {
	go a.accept()
}

func (a *acceptor) Stop() {

}

func (a *acceptor) Wait() {
	<-a.done
}

func fdToConn(fd int) (net.Conn, error) {
	f := os.NewFile(uintptr(fd), "")
	defer f.Close()
	return net.FileConn(f)
}

func ischannelClosed(c chan interface{}) bool {
	select {
	case _, ok := <-c:
		return ok
	default:
		return false
	}
}
