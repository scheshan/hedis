package hedis

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

const bufSize = 16

type requestType int

const (
	requestType_Ready requestType = iota
	requestType_SingleLine
	requestType_Bulk
)

type Session struct {
	prev    *Session
	next    *Session
	conn    *net.TCPConn
	id      int32
	server  *Server
	reqType requestType
	buf     *String
	writer  *bufio.Writer
	closed  bool
}

func NewSession(conn *net.TCPConn) *Session {
	s := new(Session)
	s.conn = conn
	s.writer = bufio.NewWriterSize(s.conn, 1024)
	s.buf = NewEmptyString()

	return s
}

func (t *Session) Server(s *Server) {
	t.server = s
}

func (t *Session) Id(id int32) {
	t.id = id
}

func (t *Session) Read() {
	go t.read()
}

func (t *Session) Close() {
	if t.closed {
		return
	}
	t.closed = true

	t.server.CloseSession(t)

	t.conn.Close()
}

func (t *Session) read() {
	buf := make([]byte, 1024, 1024)

	for {
		n, err := t.conn.Read(buf)
		if err != nil {
			if t.closed {
				return
			}
			if err != io.EOF {
				log.Printf("%s read error: %s", t, err)
			}
			t.Close()
			return
		}

		if n == 0 {
			continue
		}

		t.buf.Append(buf[:n])
		t.processBuffer()
	}
}

func (t *Session) processBuffer() {
	//at beginning, only process inline command
	for t.buf.Len() > 0 {
		i := t.buf.Index("\r\n")
		if i >= 0 {
			buf := t.buf.SliceLength(0, i)
			arr := buf.Split(" ")

			cmd := new(QueryCommand)
			cmd.session = t
			cmd.cmd = arr[0]
			cmd.arg = arr[1:]

			t.buf = t.buf.Slice(i + 2)

			t.server.EnqueueCommand(cmd)
		}
	}
}

func (t *Session) writeAndFlush(data []byte) (n int, err error) {
	n, err = t.writer.Write(data)
	if err != nil {
		return
	}

	err = t.writer.Flush()
	return
}

func (t *Session) String() string {
	return fmt.Sprintf("session: %v", t.id)
}
