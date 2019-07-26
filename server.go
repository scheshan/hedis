package hedis

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type ServerConfig struct {
	Addr string
}

type QueryCommand struct {
	session *Session
	cmd     *String
	arg     []*String
}

type Server struct {
	config   *ServerConfig
	listener *net.TCPListener
	running  bool
	id       *int32
	clients  *List
	mu       *sync.Mutex
	cmdChan  chan *QueryCommand
}

func NewServer(c *ServerConfig) *Server {
	s := new(Server)
	s.config = c
	s.id = new(int32)
	s.mu = new(sync.Mutex)
	s.clients = NewList()
	s.cmdChan = make(chan *QueryCommand, 1024)

	return s
}

func (t *Server) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", t.config.Addr)
	if err != nil {
		return err
	}

	t.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	t.running = true

	go t.listen()
	go t.processCommand()
	return nil
}

func (t *Server) Stop() error {
	t.running = false
	return t.listener.Close()
}

func (t *Server) listen() {
	for t.running {
		conn, err := t.listener.AcceptTCP()

		if err != nil {
			fmt.Println(err)
			return
		}

		session := NewSession(conn)
		session.Server(t)
		id := atomic.AddInt32(t.id, 1)
		session.Id(id)

		t.clients.AddLast(session)

		log.Printf("%s connected", session)

		session.Read()
	}
}

func (t *Server) CloseSession(s *Session) {
	t.mu.Lock()
	defer t.mu.Unlock()

	log.Printf("%s disconnected", s)

	t.clients.Remove(s)
}

func (t *Server) EnqueueCommand(cmd *QueryCommand) {
	t.cmdChan <- cmd
}

func (t *Server) processCommand() {
	for {
		c := <-t.cmdChan

		if c.cmd.Equals("exit") {
			c.session.Close()
			continue
		}

		log.Printf("command: %s, args: %s", c.cmd, c.arg)
	}
}
