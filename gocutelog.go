// Package gocutelog makes it possible to send log records from Go logging
// libraries to a cutelog instance without having to manually manage
// a socket connection.
//
// Function NewWriter returns a struct that implements io.Writer interface so
// it can be used as output by libraries like zerolog, zap, onelog, logrus, etc.
//
// Just like cutelog itself, this package is meant to be used only during
// development, so performance or reliability are not the focus here.
//
// Example usage with zerolog:
//		w := gocutelog.NewWriter("localhost:19996", "json")
//		l := zerolog.New(w)
//		l.Info().Msg("Hello world from zerolog!")
package gocutelog // import "github.com/busimus/gocutelog"

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// LogWriter manages the socket connection to cutelog. It implements
// io.Writer interface, as well as zapcore.WriteSyncer and zap.Sink.
type LogWriter struct {
	Addr   string
	Format string

	conn       net.Conn
	connecting bool
	sendLock   sync.Mutex
}

// NewWriter creates a prepared LogWriter that is ready to be used.
// If it can connect within the first 100 milliseconds then no records will
// be dropped in the beginning.
//
// Argument addr specifies the address of the cutelog instance (e.g. "localhost:19996"),
// and format specifies serialization format that will be signaled (e.g. json, msgpack, cbor).
func NewWriter(addr string, format string) *LogWriter {
	l := &LogWriter{Addr: addr, Format: format}
	connected := make(chan struct{})
	go func() {
		l.connect()
		close(connected)
	}()
	select {
	case <-connected:
	case <-time.After(100 * time.Millisecond):
	}
	return l
}

// connect tries to connect to the cutelog instance until it succeeds.
func (l *LogWriter) connect() {
	if l.connecting {
		return
	}
	l.connecting = true
	var err error
	for l.conn == nil {
		l.conn, err = net.Dial("tcp", l.Addr)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		err = l.sendMsg([]byte(fmt.Sprintf("!!cutelog!!format=%s", l.Format)))
		if err != nil {
			l.Close()
		}
	}
	l.connecting = false
}

// Write sends the encoded record if there is an established connection.
func (l *LogWriter) Write(msg []byte) (n int, err error) {
	defer func() {
		if v := recover(); v != nil {
			errMsg := fmt.Sprintf("gocutelog.Write panicked with value: %v", v)
			err = errors.New(errMsg)
		}
	}()

	if l.conn == nil || l.connecting {
		if l.connecting == false {
			go l.connect()
		}
		return len(msg), nil
	}
	return len(msg), l.sendMsg(msg)
}

// Sync is here to satisfy zapcore.WriteSyncer. Since there is no buffer
// this function does nothing.
func (l *LogWriter) Sync() (err error) {
	return
}

// Close closes the connection to cutelog.
func (l *LogWriter) Close() (err error) {
	if l.conn != nil {
		err = l.conn.Close()
		l.conn = nil
	}
	return
}

func (l *LogWriter) sendMsg(msg []byte) (err error) {
	l.sendLock.Lock()
	defer l.sendLock.Unlock()
	l.conn.SetWriteDeadline(time.Now().Add(time.Second * 2))
	err = binary.Write(l.conn, binary.BigEndian, uint32(len(msg)))
	if err != nil {
		l.Close()
		go l.connect()
		return
	}
	_, err = l.conn.Write(msg)
	l.conn.SetWriteDeadline(time.Time{})
	if err != nil {
		l.Close()
		go l.connect()
		return
	}
	return
}
