//go:build linux

package socket

import (
	"fmt"
	"net"
	"sync"
	"time"

	"syscall"
)

type NonBlockFd struct {
	fd      int
	epollfd int
	readbuf []byte
	lock    *sync.Mutex
}

func (fd *NonBlockFd) Close() error {
	return syscall.Close(fd.fd)
}

func (n *NonBlockFd) Write(buf []byte) (int, error) {
	return syscall.Write(n.fd, buf)
}

func (nonBlockFd *NonBlockFd) Read(buf []byte, timeout time.Duration) (int, error) {
	start := time.Now()
	i := 0
	for i < len(buf) {
		if i < len(nonBlockFd.readbuf) {
			buf[i] = nonBlockFd.readbuf[i]
			i++
		}
		if start.Add(time.Duration(timeout)).Before(time.Now()) {
			return -1, fmt.Errorf("timeout")
		}
		time.Sleep(time.Millisecond)
	}
	nonBlockFd.lock.Lock()
	defer nonBlockFd.lock.Unlock()
	nonBlockFd.readbuf = nonBlockFd.readbuf[i-1:]
	return i, nil
}

func New(address string, port int, socketType int, socketOptions int, proto int) (*NonBlockFd, error) {
	fd, err := syscall.Socket(socketType, socketOptions, proto)
	if err != nil {
		return nil, err
	}

	addr := &syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP(address).To4())

	err = syscall.Connect(fd, addr)
	if err != nil {
		return nil, err
	}
	syscall.SetNonblock(fd, true)

	epollfd, _ := syscall.EpollCreate1(0)

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(fd)
	syscall.EpollCtl(epollfd, syscall.EPOLL_CTL_ADD, fd, &event)
	nonblockfd := &NonBlockFd{fd, epollfd, []byte{}, &sync.Mutex{}}
	go nonblockfd.poll()
	return nonblockfd, nil
}

func (n *NonBlockFd) readfd(fd int32) error {
	n.lock.Lock()
	defer n.lock.Unlock()
	buffer := make([]byte, 4096)
	s, err := syscall.Read(int(fd), buffer)
	if err != nil {
		return err
	}
	n.readbuf = append(n.readbuf, buffer[0:s]...)
	return nil
}

func (nonBlockFd *NonBlockFd) poll() {
	for {
		var events = make([]syscall.EpollEvent, 1024)
		n, err := syscall.EpollWait(nonBlockFd.epollfd, events, 1)
		if err != nil {
			fmt.Println("EpollWait: " + err.Error())
			return
		} else {
			for i := 0; i < n; i++ {
				event := events[i]
				if nonBlockFd.readfd(event.Fd) != nil {
					fmt.Println("Socket error")
					syscall.Close(int(event.Fd))
					syscall.Close(nonBlockFd.epollfd)
					return
				}
			}
		}
	}
}
