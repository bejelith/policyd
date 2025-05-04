package chanutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenChannel(t *testing.T) {
	tests := []struct {
		err    string
		result bool
		f      func() chan interface{}
	}{
		{"Channel should be closed", true, closedChan},
		{"Channel should be open", false, openChan},
		{"Channel should be open with elements", false, openChanWithElement},
	}
	for _, test := range tests {
		t.Run(test.err, func(t *testing.T) {
			c := test.f()
			assert.Equal(t, test.result, IsClosed(c))
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
