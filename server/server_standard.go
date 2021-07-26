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
	t.addCommand("set", t.CommandSet)
	t.addCommand("get", t.CommandGet)
	t.addCommand("getset", t.CommandGetSet)
	t.addCommand("getdel", t.CommandGetDel)
	t.addCommand("strlen", t.CommandStrLen)
	t.addCommand("append", t.CommandAppend)
	t.addCommand("incr", t.CommandIncr)
	t.addCommand("decr", t.CommandDecr)
	t.addCommand("incrby", t.CommandIncrBy)
	t.addCommand("decrby", t.CommandDecrBy)
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

//#region string commands

func (t *StandardServer) CommandSet(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]
	v := args[1]

	_, obj, find, err := s.db.GetStringOrCreate(k, v)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	if find {
		obj.value = v
	}

	return codec.NewInteger(1)
}

func (t *StandardServer) CommandGet(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	str, _, _, err := s.db.GetString(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	return codec.NewBulkStr(str)
}

func (t *StandardServer) CommandGetSet(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]
	v := args[1]

	str, obj, find, err := s.db.GetString(k)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	if !find {
		return codec.NewBulkStr(nil)
	}

	obj.value = v
	return codec.NewBulkStr(str)
}

func (t *StandardServer) CommandGetDel(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]

	str, _, _, err := s.db.GetString(k)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	s.db.Remove(k)
	return codec.NewBulkStr(str)
}

func (t *StandardServer) CommandStrLen(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]

	str, _, _, err := s.db.GetString(k)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	if str == nil {
		return codec.NewInteger(0)
	}
	return codec.NewInteger(str.Len())
}

func (t *StandardServer) CommandAppend(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]
	v := args[1]

	str, _, find, err := s.db.GetStringOrCreate(k, v)
	if err != nil {
		return codec.NewErrorErr(err)
	}
	if find {
		str.AppendStr(v)
	}

	return codec.NewInteger(str.Len())
}

func (t *StandardServer) commandIncrBy(s *Session, key *core.String, num int) codec.Message {
	str, _, _, err := s.db.GetStringOrCreate(key, core.NewStringStr("0"))
	if err != nil {
		return codec.NewErrorErr(err)
	}

	res, err := str.Incr(num)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	return codec.NewInteger(res)
}

func (t *StandardServer) CommandIncr(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]

	return t.commandIncrBy(s, key, 1)
}

func (t *StandardServer) CommandDecr(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]

	return t.commandIncrBy(s, key, -1)
}

func (t *StandardServer) CommandIncrBy(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	arg := args[1]

	num, err := arg.ToInt()
	if err != nil {
		return codec.NewErrorErr(err)
	}

	return t.commandIncrBy(s, key, num)
}

func (t *StandardServer) CommandDecrBy(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	arg := args[1]

	num, err := arg.ToInt()
	if err != nil {
		return codec.NewErrorErr(err)
	}

	return t.commandIncrBy(s, key, -num)
}

//endregion
