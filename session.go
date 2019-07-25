package hedis

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

const bufSize = 16

type Session struct {
	prev   *Session
	next   *Session
	conn   *net.TCPConn
	id     int32
	rw     *bufio.ReadWriter
	server *Server
	buf    []byte
}

func NewSession(conn *net.TCPConn) *Session {
	s := new(Session)
	s.conn = conn

	reader := bufio.NewReaderSize(s.conn, bufSize)
	writer := bufio.NewWriterSize(s.conn, bufSize)
	s.rw = bufio.NewReadWriter(reader, writer)

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
	t.buf = nil
	t.server.CloseSession(t)

	t.conn.Close()
}

func (t *Session) read() {
	for {
		line, p, err := t.rw.ReadLine()

		if err != nil {
			if err != io.EOF {
				log.Printf("%s read error: %s", t, err)
			}
			t.Close()
			return
		}

		if len(line) > 0 {
			if t.buf == nil {
				t.buf = make([]byte, len(line))
			}

			t.buf = append(t.buf, line...)
		}

		if !p && t.buf != nil && len(t.buf) > 0 {
			buf := t.buf
			t.buf = nil
			_, err := t.writeAndFlush(buf)
			if err != nil {
				log.Printf("%s closed due to write error: %s", t, err)
				t.Close()
			}
		}
	}
}

func (t *Session) writeAndFlush(data []byte) (n int, err error) {
	n, err = t.rw.Write(data)
	if err != nil {
		return
	}

	err = t.rw.Flush()
	return
}

func (t *Session) String() string {
	return fmt.Sprintf("session: %v", t.id)
}
