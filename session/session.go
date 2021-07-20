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

	}
}

func (t *Session) handleError(err error) {
	t.list.Remove(t)
}

func (t *Session) ReadData() error {
	ind := 0

	for ind < bytes {
		data := t.buffer[ind:bytes]

		if t.state == nil {
			message, err := codec.ReadMessage(data)
			if err != nil {
				return err
			}

			t.state = message
			ind++
			continue
		}

		read, complete, err := t.state.Read(data)
		if err != nil {
			return err
		}

		ind += read
		if complete {
			//TODO send command

			t.state = nil
		}
	}

	return nil
}

func (t *Session) Write(msg codec.Message) error {
	return nil
}
