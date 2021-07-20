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

func (t *Session) ReadLoop() {
	for {
		reads, err := t.conn.Read(t.buffer)
		if err != nil {
			if err == net.ErrClosed {
				t.list.Remove(t)
			} else {
				//TODO close the connection
			}

			return
		}

		if reads > 0 {
			if err = t.readData(reads); err != nil {
				//TODO close the connection
				return
			}
		}
	}
}

func (t *Session) readData(bytes int) error {
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
