package handler

import (
	"bufio"
	"net"
	"sync"

	"github.com/policyd/pkg/chanutil"
	"github.com/policyd/pkg/lifecycle"
)

type ManagedHandler interface {
	lifecycle.Managed
	Handle(net.Conn)
}

type SocketHandler struct {
	messageHandler MessageHandler
	waitGroup      *sync.WaitGroup
	done           chan interface{}
}

func New() *SocketHandler {
	return &SocketHandler{handleMessage, &sync.WaitGroup{}, make(chan interface{})}
}

func (h *SocketHandler) handle(con net.Conn) {
	h.waitGroup.Add(1)
	defer h.waitGroup.Done()
	lines := make([]string, 1024)
	var newMessage = true
	var line string
	var err error
	reader := bufio.NewReader(con)
	writer := bufio.NewWriter(con)
	for {
		if chanutil.IsChannelOpen(h.done) && newMessage {
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
			response := h.messageHandler(1, lines)
			writer.WriteString(response.Resp)
			lines = []string{}
		} else {
			lines = append(lines, line)
		}
	}
}

func (h *SocketHandler) Handle(con net.Conn) {
	go h.handle(con)
}

func (h *SocketHandler) Start() {

}

func (h *SocketHandler) Stop() {
	close(h.done)
	h.waitGroup.Wait()
}
