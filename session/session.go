package session

import (
	"hedis/codec"
	"net"
)

type Session struct {
	id     int
	conn   *net.TCPConn
	pre    *Session
	next   *Session
	buffer []byte
}

func NewSession(id int, conn *net.TCPConn) *Session {
	s := &Session{}

	s.id = id
	s.conn = conn
	s.buffer = make([]byte, 40960)

	return s
}

func (t *Session) ReadLoop() error {
	for {

	}
}

func (t *Session) Write(msg codec.Message) error {
	return nil
}
