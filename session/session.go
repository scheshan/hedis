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
	list   *SessionList
	state  codec.Message
}

func NewSession(id int, conn *net.TCPConn) *Session {
	s := &Session{}

	s.id = id
	s.conn = conn
	s.buffer = make([]byte, 40960, 40960)

	return s
}

func (t *Session) ReadLoop() error {
	for {
		reads, err := t.conn.Read(t.buffer)
		if err == net.ErrClosed {
			t.list.Remove(t)
			return err
		}

		if reads > 0 {

		}
	}
}

func (t *Session) Write(msg codec.Message) error {
	return nil
}
