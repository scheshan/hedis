package server

import (
	"bufio"
	"hedis/codec"
	"io"
	"log"
	"net"
)

type SessionCloseFunc func(s *Session)

type SessionState int

const (
	StateNonAuth SessionState = iota
	StateOpen
	StateError
	StateClosed
)

type Session struct {
	id        int
	conn      *net.TCPConn
	server    Server
	db        int
	auth      bool
	pre       *Session
	next      *Session
	reader    *bufio.Reader
	writer    *bufio.Writer
	state     SessionState
	closeFunc SessionCloseFunc
}

func (t *Session) Id() int {
	return t.id
}

func (t *Session) Server() Server {
	return t.server
}

func (t *Session) Db() int {
	return t.db
}

func (t *Session) SetDb(db int) {
	t.db = db
}

func (t *Session) Auth() bool {
	return t.auth
}

func (t *Session) SetAuth(auth bool) {
	t.auth = auth
}

func (t *Session) Pre() *Session {
	return t.pre
}

func (t *Session) SetPre(pre *Session) {
	t.pre = pre
}

func (t *Session) Next() *Session {
	return t.next
}

func (t *Session) SetNext(next *Session) {
	t.next = next
}

func (t *Session) Reader() *bufio.Reader {
	return t.reader
}

func (t *Session) Writer() *bufio.Writer {
	return t.writer
}

func (t *Session) State() SessionState {
	return t.state
}

func (t *Session) SetState(state SessionState) {
	t.state = state
}

func (t *Session) CloseFunc() SessionCloseFunc {
	return t.closeFunc
}

func (t *Session) SetCloseFunc(closeFunc SessionCloseFunc) {
	t.closeFunc = closeFunc
}

func NewSession(id int, conn *net.TCPConn, server Server) *Session {
	s := &Session{}

	s.id = id
	s.conn = conn
	s.reader = bufio.NewReader(conn)
	s.writer = bufio.NewWriter(conn)
	s.server = server

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

		t.Write(codec.NewSimpleString("hello"))
	}
}

func (t *Session) handleError(err error) {
	if err != io.EOF {
		log.Print(err)
		_ = t.Close()
	}
}

func (t *Session) Write(msg codec.Message) {
	if err := codec.Encode(t.writer, msg); err != nil {
		t.handleError(err)
		return
	}

	if err := t.writer.Flush(); err != nil {
		t.handleError(err)
		return
	}
}

func (t *Session) Close() error {
	err := t.conn.Close()

	if t.closeFunc != nil {
		t.closeFunc(t)
	}

	return err
}
