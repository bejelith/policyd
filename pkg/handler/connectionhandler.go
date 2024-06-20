package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"sync"

	log "log/slog"

	"github.com/policyd/pkg/lifecycle"
	"github.com/policyd/pkg/types"
)

type ManagedHandler interface {
	lifecycle.Managed
	Handle(net.Conn)
}

type ConnHandler struct {
	messageHandler MessageHandler
	waitGroup      *sync.WaitGroup
	done           chan interface{}
}

func NewConnHandler() *ConnHandler {
	return &ConnHandler{
		messageHandler: defaultMessageHandler,
		waitGroup:      &sync.WaitGroup{},
		done:           make(chan interface{}),
	}
}

func (h *ConnHandler) handle(con net.Conn) {
	defer h.waitGroup.Done()
	defer con.Close()

	log.Info("new connection", "src", con.RemoteAddr())
	lines := bytes.NewBuffer(make([]byte, 0, 2048))
	reader := bufio.NewReaderSize(con, 2048)
	writer := bufio.NewWriterSize(con, 1)
	var terminate = false
	for {
		select {
		case <-h.done:
			terminate = true
		default:
			line, err := reader.ReadString('\n')
			if line == "\n" {
				m := types.Message{}
				m.Load(lines.Bytes())
				lines.Reset()
				result := h.messageHandler(1, m)
				if _, err := writer.WriteString(fmt.Sprintf("action=%s\n\n", result.Resp)); err != nil {
					log.Error("failed to reply", "error", err)
					return
				}
				writer.Flush()
				if terminate {
					return
				}
			} else {
				lines.WriteString(line)
			}
			if err != nil {
				log.Error("read error", "src", con.RemoteAddr(), "error", err)
				return
			}
		}
	}
}

func (h *ConnHandler) Handle(con net.Conn) {
	h.waitGroup.Add(1)
	go h.handle(con)
}

func (h *ConnHandler) Start() {

}

func (h *ConnHandler) Stop() {
	close(h.done)
}

func (h *ConnHandler) Wait() {
	h.waitGroup.Wait()
}
