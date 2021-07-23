package server

import (
	"hedis/codec"
	"hedis/core"
)

type CommandContext struct {
	session *Session
	name    *core.String
	args    []*core.String
	command Command
}

type Command func(s *Session, args []*core.String) codec.Message

var MessageErrorInvalidArgNum = codec.NewErrorString("Invalid arg num")
var MessageInvalidObjectType = codec.NewErrorString("Invalid object type")

func CommandPing(s *Session, args []*core.String) codec.Message {
	var msg codec.Message
	if len(args) == 1 {
		msg = codec.NewSimpleStr(args[0])
	} else {
		msg = codec.NewSimpleString("pong")
	}

	return msg
}

func CommandQuit(s *Session, args []*core.String) codec.Message {
	msg := codec.NewSimpleString("ok")
	s.Write(msg)
	s.Close()

	return nil
}

func CommandEcho(s *Session, args []*core.String) codec.Message {
	var msg codec.Message
	if len(args) != 1 {
		msg = MessageErrorInvalidArgNum
	} else {
		msg = codec.NewBulkStr(args[0])
	}

	return msg
}

/**  string commands start  **/

func CommandSet(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	db, err := s.Server().Db(s.db)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	k := args[0]
	v := args[1]

	obj, find := db.Get(k)
	if find {
		if obj.objType != ObjectTypeString {
			return MessageInvalidObjectType
		}

		obj.value = v
		return codec.NewInteger(1)
	}

	obj = NewObject(ObjectTypeString, v)
	if err := db.Put(k, obj); err != nil {
		return codec.NewErrorErr(err)
	}

	return codec.NewInteger(1)
}

func CommandGet(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	db, err := s.Server().Db(s.db)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	k := args[0]
	obj, find := db.Get(k)
	if !find {
		return codec.NewBulkStr(nil)
	}
	if obj.objType != ObjectTypeString {
		return MessageInvalidObjectType
	}

	str := obj.value.(*core.String)
	return codec.NewBulkStr(str)
}

func CommandGetSet(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	db, err := s.Server().Db(s.db)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	k := args[0]
	v := args[1]

	obj, find := db.Get(k)
	if !find {
		return codec.NewBulkStr(nil)
	}
	if obj.objType != ObjectTypeString {
		return MessageInvalidObjectType
	}

	ov := obj.value.(*core.String)
	obj.value = v
	return codec.NewBulkStr(ov)
}

func CommandGetDel(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	db, err := s.Server().Db(s.db)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	k := args[0]
	obj, find := db.Get(k)
	if !find {
		return codec.NewBulkStr(nil)
	}
	if obj.objType != ObjectTypeString {
		return MessageInvalidObjectType
	}

	db.Remove(k)
	str := obj.value.(*core.String)

	return codec.NewBulkStr(str)
}

func CommandStrLen(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	db, err := s.Server().Db(s.db)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	k := args[0]
	obj, find := db.Get(k)
	if !find {
		return codec.NewInteger(0)
	}
	if obj.objType != ObjectTypeString {
		return MessageInvalidObjectType
	}

	str := obj.value.(*core.String)
	return codec.NewInteger(str.Len())
}

func CommandAppend(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	db, err := s.Server().Db(s.db)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	k := args[0]
	v := args[1]

	var str *core.String

	obj, find := db.Get(k)
	if !find {
		str = v
		obj = NewObject(ObjectTypeString, str)
		if err = db.Put(k, obj); err != nil {
			return codec.NewErrorErr(err)
		}
	} else {
		if obj.objType != ObjectTypeString {
			return MessageInvalidObjectType
		}

		str = obj.value.(*core.String)
		str.AppendStr(v)
	}

	return codec.NewInteger(str.Len())
}

/**  string commands end  **/

func CommandNotFound(s *Session, args []*core.String) codec.Message {
	msg := codec.NewErrorString("Command not supported")

	return msg
}

func CommandParseFailed(s *Session, args []*core.String) codec.Message {
	msg := codec.NewErrorString("Command parse failed")

	return msg
}

type Commands struct {
	cmMap *core.Hash
}

func (t *Commands) Get(name *core.String) Command {
	i, find := t.cmMap.Get(name)

	if !find {
		return CommandNotFound
	}

	cmd, ok := i.(Command)
	if !ok {
		return CommandParseFailed
	}

	return cmd
}

func (t *Commands) add(name string, cmd Command) {
	t.cmMap.Put(core.NewStringStr(name), cmd)
}

func NewCommands() *Commands {
	cm := &Commands{}
	cm.cmMap = core.NewHashSize(100)

	return cm
}
