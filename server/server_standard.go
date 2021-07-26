package server

import (
	"errors"
	"hedis/codec"
	"hedis/core"
)

type StandardServer struct {
	*baseServer

	db []*Db
}

func (t *StandardServer) PostInit() {
	t.baseServer.PostInit()

	t.initDb()
	t.initCommands()
}

func (t *StandardServer) PreStart() {
	t.baseServer.PreStart()
}

func (t *StandardServer) PostStart() {
	t.baseServer.PostStart()
}

func (t *StandardServer) PreStop() {
	t.baseServer.PreStop()
}

func (t *StandardServer) SessionCreated(s *Session) {
	t.baseServer.SessionCreated(s)

	t.changeSessionDb(s, 0)
}

func (t *StandardServer) SessionError(s *Session, err error) {
	t.baseServer.SessionError(s, err)
}

func (t *StandardServer) SessionClosed(s *Session) {
	t.baseServer.SessionClosed(s)
}

func NewStandardServer(config *ServerConfig) *StandardServer {
	srv := &StandardServer{}
	srv.baseServer = newBaseServer(config, srv)

	return srv
}

func (t *StandardServer) initDb() {
	dbNum := 16

	t.db = make([]*Db, dbNum)
	for i := 0; i < dbNum; i++ {
		t.db[i] = NewDb()
	}
}

func (t *StandardServer) initCommands() {
	t.addCommand("select", t.CommandSelect)
}

func (t *StandardServer) changeSessionDb(s *Session, db int) error {
	if db < 0 || db >= len(t.db) {
		return errors.New("invalid db")
	}

	s.db = t.db[db]
	return nil
}

//#region connection commands

func (t *StandardServer) CommandSelect(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	dbNum, err := args[0].ToInt()
	if err != nil {
		return codec.NewErrorErr(err)
	}

	err = t.changeSessionDb(s, dbNum)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	return MessageSimpleOK
}

//#endregion
