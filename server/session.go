package server

import (
	"bufio"
	"errors"
	"hedis/codec"
	"hedis/core"
	"io"
	"log"
	"net"
	"time"
)

type SessionCloseFunc func(s *Session)

const (
	SessionFlagPubSub   = 1
	SessionFlagBlocking = 2
)

type Session struct {
	id           int
	conn         *net.TCPConn
	server       Server
	db           *Db
	auth         bool
	pre          *Session
	next         *Session
	reader       *bufio.Reader
	writer       *bufio.Writer
	closeFunc    SessionCloseFunc
	messages     chan codec.Message
	running      bool
	flag         int
	subscription *core.Hash
}

func (t *Session) canProcessMessage() bool {
	if t.flag&SessionFlagPubSub == SessionFlagPubSub {
		return false
	}
	if t.flag&SessionFlagBlocking == SessionFlagBlocking {
		return false
	}

	return true
}

func (t *Session) Id() int {
	return t.id
}

func (t *Session) Server() Server {
	return t.server
}

func (t *Session) Db() *Db {
	return t.db
}

func (t *Session) SetDb(db *Db) {
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
	s.messages = make(chan codec.Message, 1024)
	s.subscription = core.NewHash()

	log.Printf("新连接建立: %v", conn.RemoteAddr())

	return s
}

func (t *Session) handleError(err error) {
	if err == io.EOF {
		log.Printf("连接主动断开: %v", t.conn.RemoteAddr())
	} else {
		log.Printf("连接读取发生错误: %v", err)
	}

	_ = t.Close()
}

func (t *Session) processResponse() {
	for {
		select {
		case msg := <-t.messages:
			t.Write(msg)
		case <-time.After(20 * time.Second):
			continue
		}
	}
}

func (t *Session) processRequest() {
	for {
		msg, err := codec.Decode(t.reader)
		if err != nil {
			t.handleError(err)
			return
		}
		if t.flag&SessionFlagPubSub == SessionFlagPubSub {
			str := msg.Command().String()

			if str != "subscribe" && str != "psubscribe" && str != "unsubscribe" && str != "punsubscribe" && str != "ping" && str != "quit" {
				continue
			}
		}

		if t.flag&SessionFlagBlocking == SessionFlagBlocking {
			continue
		}

		t.server.QueueCommand(t, msg.Command(), msg.Args())
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

func (t *Session) StartLoop() {
	t.running = true

	go t.processRequest()
	go t.processResponse()
}

func (t *Session) QueueMessage(msg codec.Message) {
	t.messages <- msg
}

func (t *Session) SelectDb(dbNum int) error {
	db, err := t.Server().Db(dbNum)
	if err != nil {
		return err
	}
	if db == nil {
		return errors.New("Invalid db num")
	}

	t.db = db
	return nil
}

func (t *Session) Close() error {
	t.running = false

	err := t.conn.Close()

	if t.closeFunc != nil {
		t.closeFunc(t)
	}

	return err
}
