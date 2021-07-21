package session

import (
	"bufio"
	"hedis/codec"
	"net"
)

type Session struct {
	id      int
	conn    *net.TCPConn
	pre     *Session
	next    *Session
	list    *SessionList
	message codec.Message
	reader  *bufio.Reader
	writer  *bufio.Writer
}

func NewSession(id int, conn *net.TCPConn) *Session {
	s := &Session{}

	s.id = id
	s.conn = conn
	s.reader = bufio.NewReader(conn)
	s.writer = bufio.NewWriter(conn)

	return s
}

func (t *Session) ReadLoop() {
	for {
		if t.message == nil {
			msg, err := codec.ReadMessage(t.reader)
			if err != nil {
				t.handleError(err)
				return
			}

			t.message = msg
			continue
		}

		finish, err := t.message.Read(t.reader)
		if err != nil {
			t.handleError(err)
			return
		}

		if finish {
			//TODO process command

			t.message = nil
		}
	}
}

func (t *Session) handleError(err error) {
	t.list.Remove(t)
	_ = t.conn.Close()
}

func (t *Session) Write(msg codec.Message) error {
	return nil
}
