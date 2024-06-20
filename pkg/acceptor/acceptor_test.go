package acceptor

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var addr = "localhost:12345"

type MockHandler struct {
}

func (m *MockHandler) Handle(conn net.Conn) {
	if _, err := conn.Write([]byte("action=ok")); err != nil {
		conn.Close()
	}
}

func setupListener(t testing.TB, addr string) net.Listener {
	t.Helper()
	l := net.ListenConfig{}
	if addr == "" {
		addr = "localhost:0"
	}
	listener, err := l.Listen(context.TODO(), "tcp", addr)
	if err != nil {
		t.Fatal("failed to open listener", err)
	}
	t.Cleanup(func() { listener.Close() })
	return listener
}

func TestAcceptor(t *testing.T) {
	handler := &MockHandler{}
	listener := setupListener(t, "")
	acceptor, err := New(handler, listener)
	assert.Nil(t, err, "Failed to bind to socket", err)

	acceptor.Start()
	t.Cleanup(acceptor.Stop)

	clientFd, err := net.Dial("tcp", listener.Addr().String())
	assert.Nil(t, err)

	buf := make([]byte, 1024)
	_, err = clientFd.Read(buf)
	assert.Nil(t, err)
	assert.Regexp(t, "^action=ok.*", string(buf))
}

type listenerDouble struct {
	AcceptFunc func() (net.Conn, error)
}

func (l *listenerDouble) Accept() (net.Conn, error) {
	return l.AcceptFunc()
}

func (l *listenerDouble) Addr() net.Addr {
	return nil
}

func (l *listenerDouble) Close() error {
	return nil
}

func TestDoubleStart(t *testing.T) {
	handler := &MockHandler{}
	listener := listenerDouble{}
	acceptCount := 0
	mu := sync.Mutex{}
	mu.Lock()
	listener.AcceptFunc = func() (net.Conn, error) {
		acceptCount += 1
		mu.Lock()
		defer mu.Unlock()
		_, c := net.Pipe()
		return c, nil
	}
	acceptor, err := New(handler, &listener)
	assert.Nil(t, err, "Failed to bind to socket", err)
	acceptor.Start()
	acceptor.Start()
	<-time.After(time.Millisecond)
	mu.Unlock()
	acceptor.Stop()
	assert.Equal(t, 2, acceptCount, "Only two acceptor loop should have been started")
}

func BenchmarkAccept(b *testing.B) {
	handler := &MockHandler{}

	acceptor, err := New(handler, setupListener(b, ""))
	if err != nil {
		b.Fatalf("Failed to bind to socket: %v", err)
	}

	acceptor.Start()
	defer acceptor.Stop()
	b.SetParallelism(3)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			clientFd, err := net.Dial("tcp", addr)
			if err != nil {
				b.Fatal(err)
			}
			clientFd.Close()
		}

	})
}
