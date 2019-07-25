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

type Server struct {
	config   *ServerConfig
	listener *net.TCPListener
	running  bool
	id       *int32
	head     *Session
	mu       *sync.Mutex
}

func NewServer(c *ServerConfig) *Server {
	s := new(Server)
	s.config = c
	s.id = new(int32)
	s.mu = new(sync.Mutex)

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

		if t.head != nil {
			t.head.prev = session
		}
		session.next = t.head
		t.head = session

		log.Printf("%s connected", session)

		session.Read()
	}
}

func (t *Server) CloseSession(s *Session) {
	t.mu.Lock()
	defer t.mu.Unlock()

	log.Printf("%s disconnected", s)

	if s.prev == nil {
		//head
		t.head = s.next
		s.next = nil
	} else {
		prev := s.prev
		prev.next = s.next

		s.prev = nil
		s.next = nil
	}
}
