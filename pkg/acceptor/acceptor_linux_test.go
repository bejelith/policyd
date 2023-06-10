//go:build linux

package acceptor

import (
	"fmt"
	"net"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/policyd/pkg/handler"
	"github.com/policyd/pkg/socket"
)

type MockHandler struct {
}

func (m *MockHandler) Handle(conn net.Conn) {
	conn.Write([]byte(string("ok")))
}

func TestHandlerPortConflicts(t *testing.T) {
	socketHandler := handler.New()
	h1, err1 := New(socketHandler, "0.0.0.0", 12345)
	h2, err2 := New(socketHandler, "0.0.0.0", 12345)
	if err1 != nil {
		t.Fatal(err1)
	}
	h1.Start()
	time.Sleep(time.Second)
	defer h1.Stop()
	if err2 == nil {
		defer h2.Stop()
		t.Error("Should have gotten a port conflict error")
	}
}

func TestAcceptor(t *testing.T) {
	handler := &MockHandler{}
	acceptor, err := New(handler, "0.0.0.0", 12345)
	if err != nil {
		t.Fatalf("Failed to bind to socket: %v", err)
	}
	acceptor.Start()
	defer acceptor.Stop()
	fmt.Println("try dial")
	clientfd, err := socket.New("tcp", 12345, syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 2)
	count, readerr := clientfd.Read(buf, time.Second*2)
	if readerr != nil {
		t.Fatalf("%s: %v", reflect.TypeOf(readerr), readerr)
	}

	if string(buf[:count]) != "ok" {
		t.Fatal("Unexpected message")
	}
}
