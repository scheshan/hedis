package server

import (
	"errors"
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

var ErrorInvalidObjectType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
var MessageErrorInvalidArgNum = codec.NewErrorString("ERR wrong number of arguments for this command")
var MessageErrorInvalidObjectType = codec.NewErrorErr(ErrorInvalidObjectType)
var MessageSimpleOK = codec.NewSimpleString("ok")
var MessageSimpleNil = codec.NewSimpleStr(nil)
var MessageBulkNil = codec.NewBulkStr(nil)

//region connection commands

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

//endregion

//region string commands

func CommandSet(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]
	v := args[1]

	_, obj, find, err := s.Db().GetStringOrCreate(k, v)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	if find {
		obj.value = v
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

	str, _, _, err := s.Db().GetString(k)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	s.Db().Remove(k)
	return codec.NewBulkStr(str)
}

func CommandStrLen(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]

	str, _, _, err := s.Db().GetString(k)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	if str == nil {
		return codec.NewInteger(0)
	}
	return codec.NewInteger(str.Len())
}

func CommandAppend(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	k := args[0]
	v := args[1]

	str, _, find, err := s.Db().GetStringOrCreate(k, v)
	if err != nil {
		return codec.NewErrorErr(err)
	}
	if find {
		str.AppendStr(v)
	}

	return codec.NewInteger(str.Len())
}

func commandIncrBy(s *Session, key *core.String, num int) codec.Message {
	str, _, _, err := s.Db().GetStringOrCreate(key, core.NewStringStr("0"))
	if err != nil {
		return codec.NewErrorErr(err)
	}

	res, err := str.Incr(num)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	return codec.NewInteger(res)
}

func CommandIncr(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]

	return commandIncrBy(s, key, 1)
}

func CommandDecr(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]

	return commandIncrBy(s, key, -1)
}

func CommandIncrBy(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	arg := args[1]

	num, err := arg.ToInt()
	if err != nil {
		return codec.NewErrorErr(err)
	}

	return commandIncrBy(s, key, num)
}

func CommandDecrBy(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	arg := args[1]

	num, err := arg.ToInt()
	if err != nil {
		return codec.NewErrorErr(err)
	}

	return commandIncrBy(s, key, -num)
}

//endregion

// region keys commands

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

// endregion

//region hash commands

func CommandHSet(s *Session, args []*core.String) codec.Message {
	if len(args) < 3 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]

	ht, _, _, err := s.Db().GetHashOrCreate(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	for i := 1; i < len(args)-1; i += 2 {
		field := args[i]
		value := args[i+1]

		ht.Put(field, value)
	}

	return codec.NewInteger(1)
}

func CommandHExists(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]

	ht, _, find, err := s.Db().GetHash(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}
	if !find {
		return codec.NewInteger(0)
	}

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

	ht, _, find, err := s.Db().GetHash(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}
	if find {
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

	ht, _, find, err := s.Db().GetHash(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	if !find {
		return codec.NewArrayEmpty()
	}

	msg := codec.NewArraySize(ht.Size() << 1)
	ht.Iterate(func(k *core.String, v interface{}) {
		msg.AppendStr(k)
		msg.AppendStr(v.(*core.String))
	})

	return msg
}

func CommandHKeys(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]

	ht, _, find, err := s.Db().GetHash(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	if !find {
		return codec.NewArrayEmpty()
	}

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

	ht, _, find, err := s.Db().GetHash(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}
	if !find {
		return codec.NewInteger(0)
	}

	return codec.NewInteger(ht.Size())
}

func CommandHMGet(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	ht, _, find, err := s.Db().GetHash(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	msg := codec.NewArraySize(len(args) - 1)
	for i := 1; i < len(args); i++ {
		if find {
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

	key := args[0]

	ht, _, _, err := s.Db().GetHashOrCreate(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	for i := 1; i < len(args)-1; i++ {
		f := args[i]
		v := args[i+1]

		ht.Put(f, v)
	}

	return MessageSimpleOK
}

func CommandHIncrBy(s *Session, args []*core.String) codec.Message {
	if len(args) != 3 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	field := args[1]
	arg := args[2]

	num, err := arg.ToInt()
	if err != nil {
		return codec.NewErrorErr(err)
	}

	ht, _, _, err := s.Db().GetHashOrCreate(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	var str *core.String
	v, find := ht.Get(field)
	if !find {
		str = core.NewStringStr("0")
		ht.Put(field, str)
	} else {
		str = v.(*core.String)
	}

	res, err := str.Incr(num)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	return codec.NewInteger(res)
}

func CommandHDel(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	fields := args[1:]

	ht, _, find, err := s.Db().GetHash(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	num := 0
	for _, field := range fields {
		if find && ht.Remove(field) {
			num++
		}
	}

	return codec.NewInteger(num)
}

//TODO hincrbyfloat, hrandfield, hcan, hsetnx, hstrlen, hvals

//endregion

//region set commands

func CommandSAdd(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	field := args[1]

	ht, _, _, err := s.Db().GetHashOrCreate(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	if ht.Contains(field) {
		return codec.NewInteger(0)
	}

	ht.Put(field, core.HashDefaultValue)
	return codec.NewInteger(1)
}

func CommandSCard(s *Session, args []*core.String) codec.Message {
	return CommandHLen(s, args)
}

func CommandSIsMember(s *Session, args []*core.String) codec.Message {
	if len(args) != 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	field := args[1]

	ht, _, find, err := s.Db().GetHash(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	if find && ht.Contains(field) {
		return codec.NewInteger(1)
	}

	return codec.NewInteger(0)
}

func CommandSMembers(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	return CommandHKeys(s, args)
}

func CommandSMIsMember(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	fields := args[1:]

	ht, _, find, err := s.Db().GetHash(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	res := codec.NewArraySize(len(fields))
	for i := 0; i < len(fields); i++ {
		if find && ht.Contains(fields[i]) {
			res.AppendMessage(codec.NewInteger(1))
		} else {
			res.AppendMessage(codec.NewInteger(0))
		}
	}

	return res
}

func CommandSRem(s *Session, args []*core.String) codec.Message {
	return CommandHDel(s, args)
}

func CommandSRandMember(s *Session, args []*core.String) codec.Message {
	if len(args) < 1 || len(args) > 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]

	ht, _, find, err := s.Db().GetHash(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	if len(args) == 1 {
		var str *core.String
		if find {
			str, _, _ = ht.Random()
		}
		return codec.NewBulkStr(str)
	} else {
		arg := args[1]
		num, err := arg.ToInt()
		if err != nil {
			return codec.NewErrorErr(err)
		}

		arr := codec.NewArraySize(num)
		for i := 0; i < num; i++ {
			var str *core.String
			if find {
				str, _, _ = ht.Random()
			}
			arr.AppendStr(str)
		}
		return arr
	}
}

//TODO sdiff, sdiffstore, sinter, sinterstore, smove, spop, sscan, sunion, sunionstore

//endregion

//region list commands

func CommandLLen(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]

	list, _, find, err := s.Db().GetList(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	num := 0
	if find {
		num = list.Len()
	}

	return codec.NewInteger(num)
}

func CommandLPush(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	values := args[1:]

	list, _, _, err := s.Db().GetListOrCreate(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	for _, v := range values {
		list.AddHead(v)
	}

	return codec.NewInteger(list.Len())
}

func CommandLPushX(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	values := args[1:]

	list, _, find, err := s.Db().GetList(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	num := 0
	if find {
		for _, v := range values {
			list.AddHead(v)
		}
		num = list.Len()
	}

	return codec.NewInteger(num)
}

func CommandLPop(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]

	list, _, find, err := s.Db().GetList(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	var str *core.String
	if find {
		v, li, find := list.GetHead()
		if find {
			str = v.(*core.String)
			list.Remove(li)
		}
	}

	return codec.NewBulkStr(str)
}

func CommandRPop(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]

	list, _, find, err := s.Db().GetList(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	var str *core.String
	if find {
		v, li, find := list.GetTail()
		if find {
			str = v.(*core.String)
			list.Remove(li)
		}
	}

	return codec.NewBulkStr(str)
}

func CommandRPush(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	values := args[1:]

	list, _, _, err := s.Db().GetListOrCreate(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	for _, v := range values {
		list.AddTail(v)
	}

	return codec.NewInteger(list.Len())
}

func CommandRPushX(s *Session, args []*core.String) codec.Message {
	if len(args) < 2 {
		return MessageErrorInvalidArgNum
	}

	key := args[0]
	values := args[1:]

	list, _, find, err := s.Db().GetList(key)
	if err != nil {
		return codec.NewErrorErr(err)
	}

	num := 0
	if find {
		for _, v := range values {
			list.AddTail(v)
		}
		num = list.Len()
	}

	return codec.NewInteger(num)
}

//TODO blpop, brpop, brpoplpush, blmove, lindex, linsert, lpos, lrange, lrem, lset, ltrim, rpoplpush, lmove

//endregion

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
