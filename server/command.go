package server

import (
	"hedis/codec"
	"hedis/core"
	"strconv"
)

type CommandContext struct {
	session *Session
	name    *core.String
	args    []*core.String
	command Command
}

type Command func(s *Session, args []*core.String) codec.Message

var MessageErrorInvalidArgNum = codec.NewErrorString("Invalid arg num")
var MessageErrorInvalidObjectType = codec.NewErrorString("Invalid object type")
var MessageSimpleOK = codec.NewSimpleString("ok")

/**  connection commands start  **/

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
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	msg := codec.NewBulkStr(args[0])

	return msg
}

func CommandSelect(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	dbNum, err := strconv.Atoi(args[0].String())
	if err != nil {
		return codec.NewErrorErr(err)
	}

	err = s.SelectDb(dbNum)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	return MessageSimpleOK
}

/**  connection commands end  **/

/**  string commands start  **/

func CommandSet(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]
	v := args[1]

	obj, find := s.Db().Get(k)
	if find {
		if obj.objType != ObjectTypeString {
			return MessageErrorInvalidObjectType
		}

		obj.value = v
		return codec.NewInteger(1)
	}

	obj = NewObject(ObjectTypeString, v)
	if err := s.Db().Put(k, obj); err != nil {
		return codec.NewErrorErr(err)
	}

	return codec.NewInteger(1)
}

func CommandGet(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]
	obj, find := s.Db().Get(k)
	if !find {
		return codec.NewBulkStr(nil)
	}
	if obj.objType != ObjectTypeString {
		return MessageErrorInvalidObjectType
	}

	str := obj.value.(*core.String)
	return codec.NewBulkStr(str)
}

func CommandGetSet(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]
	v := args[1]

	obj, find := s.Db().Get(k)
	if !find {
		return codec.NewBulkStr(nil)
	}
	if obj.objType != ObjectTypeString {
		return MessageErrorInvalidObjectType
	}

	ov := obj.value.(*core.String)
	obj.value = v
	return codec.NewBulkStr(ov)
}

func CommandGetDel(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]
	obj, find := s.Db().Get(k)
	if !find {
		return codec.NewBulkStr(nil)
	}
	if obj.objType != ObjectTypeString {
		return MessageErrorInvalidObjectType
	}

	s.Db().Remove(k)
	str := obj.value.(*core.String)

	return codec.NewBulkStr(str)
}

func CommandStrLen(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]
	obj, find := s.Db().Get(k)
	if !find {
		return codec.NewInteger(0)
	}
	if obj.objType != ObjectTypeString {
		return MessageErrorInvalidObjectType
	}

	str := obj.value.(*core.String)
	return codec.NewInteger(str.Len())
}

func CommandAppend(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]
	v := args[1]

	var str *core.String

	obj, find := s.Db().Get(k)
	if !find {
		str = v
		obj = NewObject(ObjectTypeString, str)
		if err := s.Db().Put(k, obj); err != nil {
			return codec.NewErrorErr(err)
		}
	} else {
		if obj.objType != ObjectTypeString {
			return MessageErrorInvalidObjectType
		}

		str = obj.value.(*core.String)
		str.AppendStr(v)
	}

	return codec.NewInteger(str.Len())
}

/**  string commands end  **/

/**  keys commands start  **/

func CommandDel(s *Session, args []*core.String) codec.Message {
	if len(args) == 0 {
		return MessageErrorInvalidArgNum
	}

	res := 0
	for _, key := range args {
		if s.Db().Remove(key) {
			res++
		}
	}

	return codec.NewInteger(res)
}

func CommandExists(s *Session, args []*core.String) codec.Message {
	if len(args) == 0 {
		return MessageErrorInvalidArgNum
	}

	res := 0
	for _, key := range args {
		if s.Db().Exists(key) {
			res++
		}
	}

	return codec.NewInteger(res)
}

/**  keys commands end  **/

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
