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
	res     chan Message
}

func NewSession(conn *net.TCPConn) *Session {
	s := new(Session)
	s.conn = conn
	s.writer = bufio.NewWriterSize(s.conn, 1024)
	s.buf = NewEmptyString()
	s.res = make(chan Message, 1024)

	go s.read()
	go s.processReply()

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
	close(t.res)
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

func (t *Session) processReply() {
	for msg := range t.res {
		_, err := msg.Write(t.writer)
		if err != nil {
			log.Printf("session: %s write failed: %v", t, err)
			t.Close()
		}
		if err = t.writer.Flush(); err != nil {
			log.Printf("session: %s flush failed: %v", t, err)
			t.Close()
		}
	}
}

func (t *Session) Reply(msg Message) {
	t.res <- msg
}

func (t *Session) String() string {
	return fmt.Sprintf("session: %v", t.id)
}
