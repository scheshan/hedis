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

var MessageErrorInvalidArgNum = codec.NewErrorString("ERR wrong number of arguments for this command")
var MessageErrorInvalidObjectType = codec.NewErrorString("WRONGTYPE Operation against a key holding the wrong kind of value")
var MessageSimpleOK = codec.NewSimpleString("ok")
var MessageSimpleNil = codec.NewSimpleStr(nil)

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

/**  hash commands start  **/

func CommandHSet(s *Session, args []*core.String) codec.Message {
	if len(args) < 3 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	obj, find := s.Db().Get(key)

	if find && obj.objType != ObjectTypeHash {
		return MessageErrorInvalidObjectType
	}

	if obj == nil {
		obj = NewObject(ObjectTypeHash, core.NewHashSize(16))
		if err := s.Db().Put(key, obj); err != nil {
			return codec.NewErrorErr(err)
		}
	}

	ht := obj.value.(*core.Hash)
	for i := 1; i < len(args)-1; i += 2 {
		field := args[i]
		value := args[i+1]

		if err := ht.Put(field, value); err != nil {
			return codec.NewErrorErr(err)
		}
	}

	return codec.NewInteger(1)
}

func CommandHExists(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	obj, find := s.Db().Get(key)
	if !find {
		return codec.NewInteger(0)
	}

	if obj.objType != ObjectTypeHash {
		return MessageErrorInvalidObjectType
	}
	ht := obj.value.(*core.Hash)

	res := 0
	if ht.Contains(args[1]) {
		res = 1
	}

	return codec.NewInteger(res)
}

func CommandHGet(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	var str *core.String

	key := args[0]
	obj, find := s.Db().Get(key)
	if find {
		if obj.objType != ObjectTypeHash {
			return MessageErrorInvalidObjectType
		}

		ht := obj.value.(*core.Hash)
		i, find := ht.Get(args[1])
		if find {
			str = i.(*core.String)
		}
	}

	return codec.NewBulkStr(str)
}

func CommandHGetAll(s *Session, args []*core.String) codec.Message {
	if len(args) < 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	obj, find := s.Db().Get(key)
	if !find {
		return codec.NewArrayEmpty()
	}

	if obj.objType != ObjectTypeHash {
		return MessageErrorInvalidObjectType
	}

	ht := obj.value.(*core.Hash)
	msg := codec.NewArraySize(ht.Size() << 1)
	ht.Iterate(func(k *core.String, v interface{}) {
		msg.AppendStr(k)
		msg.AppendStr(v.(*core.String))
	})

	return msg
}

func CommandHKeys(s *Session, args []*core.String) codec.Message {
	if len(args) < 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	obj, find := s.Db().Get(key)
	if !find {
		return codec.NewArrayEmpty()
	}

	if obj.objType != ObjectTypeHash {
		return MessageErrorInvalidObjectType
	}

	ht := obj.value.(*core.Hash)
	msg := codec.NewArraySize(ht.Size())
	ht.Iterate(func(k *core.String, v interface{}) {
		msg.AppendStr(k)
	})

	return msg
}

func CommandHLen(s *Session, args []*core.String) codec.Message {
	if len(args) < 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	obj, find := s.Db().Get(key)
	if !find {
		return codec.NewInteger(0)
	}

	if obj.objType != ObjectTypeHash {
		return MessageErrorInvalidObjectType
	}

	ht := obj.value.(*core.Hash)
	return codec.NewInteger(ht.Size())
}

func CommandHMGet(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	var ht *core.Hash

	key := args[0]
	obj, find := s.Db().Get(key)
	if find {
		if obj.objType != ObjectTypeHash {
			return MessageErrorInvalidObjectType
		}

		ht = obj.value.(*core.Hash)
	}

	msg := codec.NewArraySize(len(args) - 1)
	for i := 1; i < len(args); i++ {
		if ht != nil {
			v, find := ht.Get(args[i])
			if find {
				msg.AppendStr(v.(*core.String))
				continue
			}
		}
		msg.AppendStr(nil)
	}

	return msg
}

func CommandHMSet(s *Session, args []*core.String) codec.Message {
	if len(args) < 3 {
		return MessageErrorInvalidArgNum
	}

	var ht *core.Hash

	key := args[0]
	obj, find := s.Db().Get(key)
	if find && obj.objType != ObjectTypeHash {
		return MessageErrorInvalidObjectType
	}

	if find {
		ht = obj.value.(*core.Hash)
	} else {
		ht = core.NewHashSize(16)
		obj = NewObject(ObjectTypeHash, ht)
		if err := s.Db().Put(key, obj); err != nil {
			return codec.NewErrorErr(err)
		}
	}

	for i := 1; i < len(args)-1; i++ {
		f := args[i]
		v := args[i+1]

		if err := ht.Put(f, v); err != nil {
			return codec.NewErrorErr(err)
		}
	}

	return MessageSimpleOK
}

/**  hash commands end  **/

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
