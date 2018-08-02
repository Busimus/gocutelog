package gocutelog

import (
	"bytes"
	"fmt"
	"net"
	"testing"
)

func TestWriting(t *testing.T) {
	cases := []struct{ len, msg []byte }{
		{[]byte("\x00\x00\x00\x00"), []byte("")},
		{[]byte("\x00\x00\x00\x02"), []byte("{}")},
		{[]byte("\x00\x00\x00\x82"), []byte(`{"name": "MyServer.ReqHandler", "level": "debug", "created": 1528702099, "msg": "User registered", "username": "bob", "id": 13525}`)},
	}
	server, client := net.Pipe()
	w := LogWriter{conn: client, Format: "json", connecting: false}
	reader := func() {
		var length, msg []byte
		for _, tc := range cases {
			length = make([]byte, len(tc.len))
			msg = make([]byte, len(tc.msg))
			n, err := server.Read(length)
			if bytes.Compare(length, tc.len) != 0 || err != nil {
				if err != nil {
					fmt.Printf("read %d, err: \"%s\"\n", n, err.Error())
				}
				fmt.Printf("got:      %q\n", length)
				fmt.Printf("expected: %q\n", tc.len)
				t.Fail()
			}
			n, err = server.Read(msg)
			if bytes.Compare(msg, tc.msg) != 0 || err != nil {
				if err != nil {
					fmt.Printf("read %d, err: \"%s\"\n", n, err.Error())
				}
				fmt.Printf("got:      %q\n", msg)
				fmt.Printf("expected: %q\n", tc.msg)
				t.Fail()
			}
		}
	}
	go reader()
	for _, tc := range cases {
		w.Write(tc.msg)
	}
}
