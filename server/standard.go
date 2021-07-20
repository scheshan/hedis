package server

import (
	"hedis/config"
	"hedis/session"
	"log"
	"net"
)

type StandardServer struct {
	config   *config.Config
	listener net.Listener
	sessions *session.SessionList
	running  bool
	clientId int
}

func NewStandard(c *config.Config) Server {
	server := &StandardServer{}
	server.config = c
	server.sessions = session.NewSessionList()

	return server
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
	return nil
}

func (t *StandardServer) accept() {
	for t.running {
		conn, err := t.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		t.clientId++
		session := session.NewSession(t.clientId, conn.(*net.TCPConn))
		t.sessions.AddLast(session)
	}
}

func (t *StandardServer) Stop() error {
	panic("implement me")
}
