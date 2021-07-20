package server

import (
	"hedis/config"
	"net"
)

type StandardServer struct {
	config   *config.Config
	listener net.Listener
}

func NewStandard(c *config.Config) Server {
	server := &StandardServer{}
	server.config = c

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

	go t.accept()
	return nil
}

func (t *StandardServer) accept() {

}

func (t *StandardServer) Stop() error {
	panic("implement me")
}
