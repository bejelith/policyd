package handler

import (
	"bufio"
	"net"
)

type SocketHandler interface {
	Handle(net.Conn)
}

type socketHandler struct {
	messageHandler MessageHandler
}

func (h *socketHandler) handle(con net.Conn) {
	lines := make([]string, 1024)
	var line string
	var err error
	reader := bufio.NewReader(con)
	line, err = reader.ReadString('\n')
	lines = append(lines, line)
	for err == nil {
		line, err = reader.ReadString('\n')
		if line == "" {
			h.messageHandler.Handle(lines)
			lines = []string{}
		} else {
			lines = append(lines, line)
		}
	}
}

func (h *socketHandler) Handle(con net.Conn) {
	go h.handle(con)
}

func NewSocketHandler() SocketHandler {
	return &socketHandler{}
}

type MessageHandler interface {
	Handle([]string)
}
