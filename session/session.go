package session

import (
	"bufio"
	"hedis/codec"
	"io"
	"log"
	"net"
)

type Session struct {
	id      int
	state   State
	conn    *net.TCPConn
	pre     *Session
	next    *Session
	list    *SessionList
	message codec.Message
	reader  *bufio.Reader
	writer  *bufio.Writer
}

func NewSession(id int, conn *net.TCPConn, list *SessionList) *Session {
	s := &Session{}

	s.id = id
	s.conn = conn
	s.reader = bufio.NewReader(conn)
	s.writer = bufio.NewWriter(conn)
	s.list = list

	log.Printf("新连接建立: %v\r\n", conn.RemoteAddr())

	return s
}

func (t *Session) ReadLoop() {
	for {
		msg, err := codec.Decode(t.reader)
		if err != nil {
			t.handleError(err)
			return
		}

		str, err := codec.EncodeString(msg)
		if err != nil {
			t.handleError(err)
			return
		}

		log.Print(str)

		if err := codec.Encode(t.writer, codec.NewSimpleString("hello")); err != nil {
			t.handleError(err)
			return
		}
		if err := t.writer.Flush(); err != nil {
			t.handleError(err)
			return
		}
	}
}

func (t *Session) handleError(err error) {
	if err != io.EOF {
		log.Print(err)
		_ = t.conn.Close()
	}
	t.list.Remove(t)
}

func (t *Session) Write(msg codec.Message) error {
	return nil
}
