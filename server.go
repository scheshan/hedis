package hedis

import (
	"fmt"
	"net"
)

type ServerConfig struct {
	Addr string
}

type Server struct {
	config   *ServerConfig
	listener *net.TCPListener
	running  bool
}

func NewServer(c *ServerConfig) *Server {
	s := new(Server)
	s.config = c

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
		_, err := t.listener.AcceptTCP()

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
