package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/policyd/pkg/types"
)

var message = `request=smtpd_access_policy
protocol_state=RCPT
protocol_name=SMTP
helo_name=some.domain.tld
queue_id=8045F2AB23
sender=foo@bar.tld
recipient=bar@foo.tld
recipient_count=0
client_address=1.2.3.4
client_name=another.domain.tld
reverse_client_name=another.domain.tld
instance=123.456.7
sasl_method=plain
sasl_username=you
sasl_sender=
size=12345
ccert_subject=solaris9.porcupine.org
ccert_issuer=Wietse+20Venema
ccert_fingerprint=C2:9D:F4:87:71:73:73:D9:18:E7:C2:F3:C1:DA:6E:04
encryption_protocol=TLSv1/SSLv3
encryption_cipher=DHE-RSA-AES256-SHA
encryption_keysize=256
etrn_domain=
stress=
ccert_pubkey_fingerprint=68:B3:29:DA:98:93:E3:40:99:C7:D8:AD:5C:B9:C9:40
client_port=1234
policy_context=submission
server_address=10.3.2.1
server_port=54321
compatibility_level=major.minor.patch
mail_version=3.8.0`

func TestConnHandler(t *testing.T) {
	server, client := net.Pipe()
	handler := NewConnHandler()
	messages := make(chan types.Message, 128)
	handler.messageHandler = func(_ types.RequestID, m types.Message) *types.Response {
		messages <- m
		return types.OkResponse
	}
	handler.Start()
	handler.Handle(server)

	tests := []struct {
		kv               map[string]string
		messageDelimiter string
	}{{
		kv:               map[string]string{"a": "b", "1": "2"},
		messageDelimiter: "\n",
	}, {
		kv:               map[string]string{"a": "b", "1": "2"},
		messageDelimiter: "",
	}}
	for _, test := range tests {
		_, err := client.Write([]byte(buildPayload(test.kv, test.messageDelimiter)))
		if err != nil {
			t.Fatal(err)
		}
		select {
		case message := <-messages:
			for k, v := range test.kv {
				if message[k] != v {
					t.Fatalf("expected \"%s\" but received \"%s\"", v, message[k])
				}
			}
			// Read postfix response
			response, err := readResponse(client)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.HasPrefix(response, "action=") {
				t.Fatalf("unexpected response \"%s\"", response)
			}
		case <-time.After(time.Second * 2):
			if test.messageDelimiter == "\n" {
				t.Fatalf("should not timeout when delimiter is \"%s\"", test.messageDelimiter)
			}
		}
	}
}

func TestReceiveEmptyMessages(t *testing.T) {
	client, server := net.Pipe()
	handler := NewConnHandler()
	messages := make(chan types.Message)
	handler.messageHandler = func(_ types.RequestID, m types.Message) *types.Response {
		messages <- m
		return types.OkResponse
	}
	handler.Handle(server)
	for range 10 {
		_, err := client.Write([]byte("\n"))
		if err != nil {
			t.Fatal(err)
		}
		select {
		case m := <-messages:
			if len(m) != 0 {
				t.Fatal("message should be empty")
			}
			_, _ = readResponse(client)
		case <-time.After(time.Second):
			t.Fatal("timeout reading message")
		}
	}
}

func BenchmarkConnectionHandler(b *testing.B) {
	handler := NewConnHandler()
	messages := make(chan types.Message, b.N)
	handler.messageHandler = func(ri types.RequestID, m types.Message) *types.Response {
		messages <- m
		return types.OkResponse
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c, s := net.Pipe()

			handler.Handle(s)
			c.Write([]byte("a=b\n\n"))
			_, err := readResponse(c)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

}

func readResponse(c net.Conn) (string, error) {
	r := bufio.NewReader(c)
	response, err := r.ReadString('\n')
	if err != nil {
		return response, err
	}
	_, err = r.ReadString('\n')
	return response, err

}

func buildPayload(m map[string]string, d string) string {
	b := bytes.NewBufferString("")
	for k, v := range m {
		b.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}
	b.WriteString(d)
	return b.String()
}
