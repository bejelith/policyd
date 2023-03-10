package handler

import (
	"bufio"
	"net"
	"sync"

	"githun.com/policyd/pkg/chanutil"
	"githun.com/policyd/pkg/lifecycle"
)

type SocketHandler interface {
	lifecycle.Managed
	Handle(net.Conn)
}

type socketHandler struct {
	messageHandler MessageHandler
	waitGroup      *sync.WaitGroup
	done           chan interface{}
}

func New(handler MessageHandler) SocketHandler {
	return &socketHandler{handler, &sync.WaitGroup{}, make(chan interface{})}
}

func (h *socketHandler) handle(con net.Conn) {
	h.waitGroup.Add(1)
	defer h.waitGroup.Done()
	lines := make([]string, 1024)
	var newMessage = true
	var line string
	var err error
	reader := bufio.NewReader(con)
	writer := bufio.NewWriter(con)
	for {
		if chanutil.IschannelClosed(h.done) && newMessage {
			con.Close()
			return
		}
		line, err = reader.ReadString('\n')
		if err != nil {
			return
		}
		newMessage = false
		if line == "" {
			newMessage = true
			writer.WriteString(h.messageHandler.Handle(lines))
			lines = []string{}
		} else {
			lines = append(lines, line)
		}
	}
}

func (h *socketHandler) Handle(con net.Conn) {
	go h.handle(con)
}

func (h *socketHandler) Start() {

}

func (h *socketHandler) Stop() {
	close(h.done)
	h.waitGroup.Wait()
}
