package types

import (
	"bufio"
	"bytes"
	"strings"
)

type RequestID int64

type Message map[string]string

const lineSeparator = "="

func (m Message) Load(b []byte) {
	r := bufio.NewReader(bytes.NewReader(b))
	for l, _, e := r.ReadLine(); e == nil; l, _, e = r.ReadLine() {
		if len(l) > 0 {
			tokens := strings.SplitN(string(l), lineSeparator, 2)
			if len(tokens) == 2 {
				m[tokens[0]] = tokens[1]
			}
		}
	}
}

type Response struct {
	Resp     string
	Failures []error
}

var OkResponse = &Response{Resp: "ok"}
