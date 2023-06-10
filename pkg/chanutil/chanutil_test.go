package chanutil

import (
	"testing"
)

func TestOpenChannel(t *testing.T) {
	tests := []struct {
		err    string
		result bool
		f      func() chan interface{}
	}{
		{"Channel should be closed", false, closedChan},
		{"Channel should be open", true, openChan},
		{"Channel should be open with elements", true, openChanWithElement},
	}
	for _, test := range tests {
		t.Run(test.err, func(t *testing.T) {
			c := test.f()
			if IsChannelOpen(c) != test.result {
				t.Fatal("Failed")
			}
		})
	}
}

func openChan() chan interface{} {
	return make(chan interface{}, 1024)
}

func openChanWithElement() chan interface{} {
	c := openChan()
	c <- nil
	return c
}

func closedChan() chan interface{} {
	c := openChan()
	close(c)
	return c
}
