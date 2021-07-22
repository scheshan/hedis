package server

import (
	"hedis/core"
	"log"
	"net"
)

type Server interface {
	Run() error
	Stop() error
	QueueCommand(session *Session, name *core.String, args []*core.String)
}

type StandardServer struct {
	cm       *Commands
	config   *ServerConfig
	listener net.Listener
	session  *Session
	running  bool
	clientId int
	requests chan *CommandContext
	db       []*Db
}

func NewStandard(c *ServerConfig) Server {
	server := &StandardServer{}
	server.config = c
	server.requests = make(chan *CommandContext, 102400)
	server.initCommands()

	return server
}

func (t *StandardServer) initCommands() {
	cm := NewCommands()
	cm.add("ping", CommandPing)
	cm.add("quit", CommandQuit)
	cm.add("echo", CommandEcho)

	t.cm = cm
}

func (t *StandardServer) initDb() {
	dbSize := 16

	t.db = make([]*Db, dbSize)
	for i := 0; i < dbSize; i++ {
		t.db[i] = NewDb()
	}
}

func (t *StandardServer) accept() {
	for t.running {
		conn, err := t.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		t.clientId++
		s := NewSession(t.clientId, conn.(*net.TCPConn), t)
		s.SetCloseFunc(t.onSessionClose)

		s.SetNext(t.session)
		if t.session != nil {
			t.session.SetPre(s)
		}
		t.session = s

		s.StartLoop()
	}
}

func (t *StandardServer) processRequest() {
	for t.running {
		ctx := <-t.requests

		msg := ctx.command(ctx.session, ctx.args)

		if msg != nil {
			ctx.session.QueueMessage(msg)
		}
	}
}

func (t *StandardServer) onSessionClose(s *Session) {
	if s.Pre() != nil {
		s.Pre().SetNext(s.Next())
	}

	if s.Next() != nil {
		s.Next().SetPre(s.Pre())
	}

	if s == t.session {
		t.session = t.session.Next()
	}

	s.SetPre(nil)
	s.SetNext(nil)
}

func (t *StandardServer) Run() error {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:6379")
	if err != nil {
		return err
	}

	var listener net.Listener
	listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	t.listener = listener

	t.running = true

	go t.accept()
	go t.processRequest()
	return nil
}

func (t *StandardServer) Stop() error {
	panic("implement me")
}

func (t *StandardServer) QueueCommand(session *Session, name *core.String, args []*core.String) {
	ctx := &CommandContext{}
	ctx.session = session
	ctx.name = name
	ctx.args = args
	ctx.command = t.cm.Get(ctx.name)

	t.requests <- ctx
}
