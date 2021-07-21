package session

import (
	"bufio"
	"fmt"
	"hedis/codec"
	"log"
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

func NewSession(id int, conn *net.TCPConn, list *SessionList) *Session {
	s := &Session{}

	s.id = id
	s.conn = conn
	s.reader = bufio.NewReader(conn)
	s.writer = bufio.NewWriter(conn)
	s.list = list

	return s
}

func (t *Session) ReadLoop() {
	for {
		msg, err := codec.ReadMessage(t.reader)
		if err != nil {
			t.handleError(err)
			return
		}

		if err = msg.Read(t.reader); err != nil {
			t.handleError(err)
			return
		}

		fmt.Println("命令可以执行了")
		fmt.Println(msg.String())
	}
}

func (t *Session) handleError(err error) {
	log.Print(err)
	t.list.Remove(t)
	_ = t.conn.Close()
}

func (t *Session) Write(msg codec.Message) error {
	return nil
}
