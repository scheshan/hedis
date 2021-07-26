package server

import (
	"bufio"
	"fmt"
	"hedis/codec"
	"hedis/core"
	"net"
)

type SessionCloseFunc func(s *Session)

type SessionState int

const (
	SessionFlagPubSub   = 1
	SessionFlagBlocking = 2
)

const (
	SessionStateConnected SessionState = iota
	SessionStateClosed
)

type Session struct {
	id           int
	conn         net.Conn
	server       Server
	state        SessionState
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
	listBlocking *core.Hash
}

func (t *Session) String() string {
	return fmt.Sprintf("session-%d[%s]", t.id, t.conn.RemoteAddr().String())
}
