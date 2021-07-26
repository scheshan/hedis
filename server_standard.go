package hedis

import (
	"errors"
	"time"
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

	go t.cleanTimeoutKeys()
}

func (t *StandardServer) PostStart() {
	t.baseServer.PostStart()
}

func (t *StandardServer) PreStop() {
	t.baseServer.PreStop()
}

func (t *StandardServer) SessionCreated(s *Session) {
	t.changeSessionDb(s, 0)

	t.baseServer.SessionCreated(s)
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
	t.addCommand("dbsize", t.CommandDbSize)

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

	t.addCommand("del", t.CommandDel)
	t.addCommand("exists", t.CommandExists)
	t.addCommand("expire", t.CommandExpire)

	t.addCommand("hset", t.CommandHSet)
	t.addCommand("hexists", t.CommandHExists)
	t.addCommand("hget", t.CommandHGet)
	t.addCommand("hgetall", t.CommandHGetAll)
	t.addCommand("hkeys", t.CommandHKeys)
	t.addCommand("hlen", t.CommandHLen)
	t.addCommand("hmget", t.CommandHMGet)
	t.addCommand("hmset", t.CommandHMSet)
	t.addCommand("hincrby", t.CommandHIncrBy)
	t.addCommand("hdel", t.CommandHDel)
	t.addCommand("hstrlen", t.CommandHStrLen)

	t.addCommand("sadd", t.CommandSAdd)
	t.addCommand("scard", t.CommandSCard)
	t.addCommand("sismember", t.CommandSIsMember)
	t.addCommand("smembers", t.CommandSMembers)
	t.addCommand("smismember", t.CommandSMIsMember)
	t.addCommand("srem", t.CommandSRem)
	t.addCommand("srandmember", t.CommandSRandMember)

	t.addCommand("llen", t.CommandLLen)
	t.addCommand("lpush", t.CommandLPush)
	t.addCommand("lpushx", t.CommandLPushX)
	t.addCommand("lpop", t.CommandLPop)
	t.addCommand("rpop", t.CommandRPop)
	t.addCommand("rpush", t.CommandRPush)
	t.addCommand("rpushx", t.CommandRPushX)
}

func (t *StandardServer) changeSessionDb(s *Session, db int) error {
	if db < 0 || db >= len(t.db) {
		return errors.New("invalid db")
	}

	s.db = t.db[db]
	return nil
}

func (t *StandardServer) cleanTimeoutKeys() {
	num := 10

	for t.running {
		<-time.After(time.Second * 10)

		for _, db := range t.db {
			if db.expires.Empty() {
				continue
			}

			for i := 0; i < num; i++ {
				key, value, find := db.expires.Random()
				if find {
					obj := value.(*Object)
					if obj.Expired() {
						db.Remove(key)
					}
				}
			}
		}
	}
}

//#region server commands

func (t *StandardServer) CommandDbSize(s *Session, args []*String) Message {
	num := s.db.ht.Size()
	return NewInteger(num)
}

//endregion

//#region connection commands

func (t *StandardServer) CommandSelect(s *Session, args []*String) Message {
	if len(args) != 1 {
		return ErrorInvalidArgNum
	}

	dbNum, err := args[0].ToInt()
	if err != nil {
		return NewErrorErr(err)
	}

	err = t.changeSessionDb(s, dbNum)
	if err != nil {
		return NewErrorErr(err)
	}

	return SimpleOK
}

//#endregion

//#region string commands

func (t *StandardServer) CommandSet(s *Session, args []*String) Message {
	if len(args) < 2 {
		return ErrorInvalidArgNum
	}

	k := args[0]
	v := args[1]

	_, obj, find, err := s.db.GetStringOrCreate(k, v)
	if err != nil {
		return NewErrorErr(err)
	}

	if find {
		obj.value = v
	}

	return NewInteger(1)
}

func (t *StandardServer) CommandGet(s *Session, args []*String) Message {
	if len(args) != 1 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	str, _, _, err := s.db.GetString(key)
	if err != nil {
		return NewErrorErr(err)
	}

	return NewBulkStr(str)
}

func (t *StandardServer) CommandGetSet(s *Session, args []*String) Message {
	if len(args) != 2 {
		return ErrorInvalidArgNum
	}

	k := args[0]
	v := args[1]

	str, obj, find, err := s.db.GetString(k)
	if err != nil {
		return NewErrorErr(err)
	}

	if !find {
		return NewBulkStr(nil)
	}

	obj.value = v
	return NewBulkStr(str)
}

func (t *StandardServer) CommandGetDel(s *Session, args []*String) Message {
	if len(args) != 1 {
		return ErrorInvalidArgNum
	}

	k := args[0]

	str, _, _, err := s.db.GetString(k)
	if err != nil {
		return NewErrorErr(err)
	}

	s.db.Remove(k)
	return NewBulkStr(str)
}

func (t *StandardServer) CommandStrLen(s *Session, args []*String) Message {
	if len(args) != 1 {
		return ErrorInvalidArgNum
	}

	k := args[0]

	str, _, _, err := s.db.GetString(k)
	if err != nil {
		return NewErrorErr(err)
	}

	if str == nil {
		return NewInteger(0)
	}
	return NewInteger(str.Len())
}

func (t *StandardServer) CommandAppend(s *Session, args []*String) Message {
	if len(args) != 2 {
		return ErrorInvalidArgNum
	}

	k := args[0]
	v := args[1]

	str, _, find, err := s.db.GetStringOrCreate(k, v)
	if err != nil {
		return NewErrorErr(err)
	}
	if find {
		str.AppendStr(v)
	}

	return NewInteger(str.Len())
}

func (t *StandardServer) commandIncrBy(s *Session, key *String, num int) Message {
	str, _, _, err := s.db.GetStringOrCreate(key, NewStringStr("0"))
	if err != nil {
		return NewErrorErr(err)
	}

	res, err := str.Incr(num)
	if err != nil {
		return NewErrorErr(err)
	}

	return NewInteger(res)
}

func (t *StandardServer) CommandIncr(s *Session, args []*String) Message {
	if len(args) != 1 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	return t.commandIncrBy(s, key, 1)
}

func (t *StandardServer) CommandDecr(s *Session, args []*String) Message {
	if len(args) != 1 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	return t.commandIncrBy(s, key, -1)
}

func (t *StandardServer) CommandIncrBy(s *Session, args []*String) Message {
	if len(args) != 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	arg := args[1]

	num, err := arg.ToInt()
	if err != nil {
		return NewErrorErr(err)
	}

	return t.commandIncrBy(s, key, num)
}

func (t *StandardServer) CommandDecrBy(s *Session, args []*String) Message {
	if len(args) != 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	arg := args[1]

	num, err := arg.ToInt()
	if err != nil {
		return NewErrorErr(err)
	}

	return t.commandIncrBy(s, key, -num)
}

//endregion

// region keys commands

func (t *StandardServer) CommandDel(s *Session, args []*String) Message {
	if len(args) == 0 {
		return ErrorInvalidArgNum
	}

	res := 0
	for _, key := range args {
		if s.db.Remove(key) {
			res++
		}
	}

	return NewInteger(res)
}

func (t *StandardServer) CommandExists(s *Session, args []*String) Message {
	if len(args) == 0 {
		return ErrorInvalidArgNum
	}

	res := 0
	for _, key := range args {
		if s.db.Exists(key) {
			res++
		}
	}

	return NewInteger(res)
}

func (t *StandardServer) CommandExpire(s *Session, args []*String) Message {
	if len(args) < 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	sec, err := args[1].ToInt()
	if err != nil {
		return NewErrorErr(err)
	}

	obj, find := s.db.Get(key)
	if !find {
		return IntegerZero
	}

	dur := time.Second * time.Duration(sec)
	expire := time.Now().Add(dur)
	obj.ttl = expire.Unix()
	s.db.expires.Put(key, obj)

	return IntegerOne
}

// endregion

//region hash commands

func (t *StandardServer) CommandHSet(s *Session, args []*String) Message {
	if len(args) < 3 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	ht, _, _, err := s.db.GetHashOrCreate(key)
	if err != nil {
		return NewErrorErr(err)
	}

	for i := 1; i < len(args)-1; i += 2 {
		field := args[i]
		value := args[i+1]

		ht.Put(field, value)
	}

	return NewInteger(1)
}

func (t *StandardServer) CommandHExists(s *Session, args []*String) Message {
	if len(args) < 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	ht, _, find, err := s.db.GetHash(key)
	if err != nil {
		return NewErrorErr(err)
	}
	if !find {
		return NewInteger(0)
	}

	res := 0
	if ht.Contains(args[1]) {
		res = 1
	}

	return NewInteger(res)
}

func (t *StandardServer) CommandHGet(s *Session, args []*String) Message {
	if len(args) < 2 {
		return ErrorInvalidArgNum
	}

	var str *String

	key := args[0]

	ht, _, find, err := s.db.GetHash(key)
	if err != nil {
		return NewErrorErr(err)
	}
	if find {
		i, find := ht.Get(args[1])
		if find {
			str = i.(*String)
		}
	}

	return NewBulkStr(str)
}

func (t *StandardServer) CommandHGetAll(s *Session, args []*String) Message {
	if len(args) < 1 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	ht, _, find, err := s.db.GetHash(key)
	if err != nil {
		return NewErrorErr(err)
	}

	if !find {
		return NewArrayEmpty()
	}

	msg := NewArraySize(ht.Size() << 1)
	ht.Iterate(func(k *String, v interface{}) {
		msg.AppendStr(k)
		msg.AppendStr(v.(*String))
	})

	return msg
}

func (t *StandardServer) CommandHKeys(s *Session, args []*String) Message {
	if len(args) != 1 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	ht, _, find, err := s.db.GetHash(key)
	if err != nil {
		return NewErrorErr(err)
	}

	if !find {
		return NewArrayEmpty()
	}

	msg := NewArraySize(ht.Size())
	ht.Iterate(func(k *String, v interface{}) {
		msg.AppendStr(k)
	})

	return msg
}

func (t *StandardServer) CommandHLen(s *Session, args []*String) Message {
	if len(args) < 1 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	ht, _, find, err := s.db.GetHash(key)
	if err != nil {
		return NewErrorErr(err)
	}
	if !find {
		return NewInteger(0)
	}

	return NewInteger(ht.Size())
}

func (t *StandardServer) CommandHMGet(s *Session, args []*String) Message {
	if len(args) < 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	ht, _, find, err := s.db.GetHash(key)
	if err != nil {
		return NewErrorErr(err)
	}

	msg := NewArraySize(len(args) - 1)
	for i := 1; i < len(args); i++ {
		if find {
			v, find := ht.Get(args[i])
			if find {
				msg.AppendStr(v.(*String))
				continue
			}
		}
		msg.AppendStr(nil)
	}

	return msg
}

func (t *StandardServer) CommandHMSet(s *Session, args []*String) Message {
	if len(args) < 3 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	ht, _, _, err := s.db.GetHashOrCreate(key)
	if err != nil {
		return NewErrorErr(err)
	}

	for i := 1; i < len(args)-1; i++ {
		f := args[i]
		v := args[i+1]

		ht.Put(f, v)
	}

	return SimpleOK
}

func (t *StandardServer) CommandHIncrBy(s *Session, args []*String) Message {
	if len(args) != 3 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	field := args[1]
	arg := args[2]

	num, err := arg.ToInt()
	if err != nil {
		return NewErrorErr(err)
	}

	ht, _, _, err := s.db.GetHashOrCreate(key)
	if err != nil {
		return NewErrorErr(err)
	}

	var str *String
	v, find := ht.Get(field)
	if !find {
		str = NewStringStr("0")
		ht.Put(field, str)
	} else {
		str = v.(*String)
	}

	res, err := str.Incr(num)
	if err != nil {
		return NewErrorErr(err)
	}

	return NewInteger(res)
}

func (t *StandardServer) CommandHDel(s *Session, args []*String) Message {
	if len(args) != 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	fields := args[1:]

	ht, _, find, err := s.db.GetHash(key)
	if err != nil {
		return NewErrorErr(err)
	}

	num := 0
	for _, field := range fields {
		if find && ht.Remove(field) {
			num++
		}
	}

	return NewInteger(num)
}

func (t *StandardServer) CommandHStrLen(s *Session, args []*String) Message {
	if len(args) != 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	field := args[1]

	ht, _, find, err := s.db.GetHash(key)
	if err != nil {
		return NewErrorErr(err)
	}

	num := 0
	if find {
		v, f := ht.Get(field)
		if f {
			str := v.(*String)
			num = str.Len()
		}
	}

	return NewInteger(num)
}

//TODO hincrbyfloat, hrandfield, hcan, hsetnx, hvals

//endregion

//region set commands

func (t *StandardServer) CommandSAdd(s *Session, args []*String) Message {
	if len(args) != 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	field := args[1]

	ht, _, _, err := s.db.GetHashOrCreate(key)
	if err != nil {
		return NewErrorErr(err)
	}

	if ht.Contains(field) {
		return NewInteger(0)
	}

	ht.Put(field, HashDefaultValue)
	return NewInteger(1)
}

func (t *StandardServer) CommandSCard(s *Session, args []*String) Message {
	return t.CommandHLen(s, args)
}

func (t *StandardServer) CommandSIsMember(s *Session, args []*String) Message {
	if len(args) != 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	field := args[1]

	ht, _, find, err := s.db.GetHash(key)
	if err != nil {
		return NewErrorErr(err)
	}

	if find && ht.Contains(field) {
		return NewInteger(1)
	}

	return NewInteger(0)
}

func (t *StandardServer) CommandSMembers(s *Session, args []*String) Message {
	if len(args) != 1 {
		return ErrorInvalidArgNum
	}

	return t.CommandHKeys(s, args)
}

func (t *StandardServer) CommandSMIsMember(s *Session, args []*String) Message {
	if len(args) < 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	fields := args[1:]

	ht, _, find, err := s.db.GetHash(key)
	if err != nil {
		return NewErrorErr(err)
	}

	res := NewArraySize(len(fields))
	for i := 0; i < len(fields); i++ {
		if find && ht.Contains(fields[i]) {
			res.AppendMessage(NewInteger(1))
		} else {
			res.AppendMessage(NewInteger(0))
		}
	}

	return res
}

func (t *StandardServer) CommandSRem(s *Session, args []*String) Message {
	return t.CommandHDel(s, args)
}

func (t *StandardServer) CommandSRandMember(s *Session, args []*String) Message {
	if len(args) < 1 || len(args) > 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	ht, _, find, err := s.db.GetHash(key)
	if err != nil {
		return NewErrorErr(err)
	}

	if len(args) == 1 {
		var str *String
		if find {
			str, _, _ = ht.Random()
		}
		return NewBulkStr(str)
	} else {
		arg := args[1]
		num, err := arg.ToInt()
		if err != nil {
			return NewErrorErr(err)
		}

		arr := NewArraySize(num)
		for i := 0; i < num; i++ {
			var str *String
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

func (t *StandardServer) CommandLLen(s *Session, args []*String) Message {
	if len(args) != 1 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	list, _, find, err := s.db.GetList(key)
	if err != nil {
		return NewErrorErr(err)
	}

	num := 0
	if find {
		num = list.Len()
	}

	return NewInteger(num)
}

func (t *StandardServer) CommandLPush(s *Session, args []*String) Message {
	if len(args) < 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	values := args[1:]

	list, _, _, err := s.db.GetListOrCreate(key)
	if err != nil {
		return NewErrorErr(err)
	}

	for _, v := range values {
		list.AddHead(v)
	}

	return NewInteger(list.Len())
}

func (t *StandardServer) CommandLPushX(s *Session, args []*String) Message {
	if len(args) < 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	values := args[1:]

	list, _, find, err := s.db.GetList(key)
	if err != nil {
		return NewErrorErr(err)
	}

	num := 0
	if find {
		for _, v := range values {
			list.AddHead(v)
		}
		num = list.Len()
	}

	return NewInteger(num)
}

func (t *StandardServer) CommandLPop(s *Session, args []*String) Message {
	if len(args) != 1 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	list, _, find, err := s.db.GetList(key)
	if err != nil {
		return NewErrorErr(err)
	}

	var str *String
	if find {
		v, li, find := list.GetHead()
		if find {
			str = v.(*String)
			list.Remove(li)
		}
	}

	return NewBulkStr(str)
}

func (t *StandardServer) CommandRPop(s *Session, args []*String) Message {
	if len(args) != 1 {
		return ErrorInvalidArgNum
	}

	key := args[0]

	list, _, find, err := s.db.GetList(key)
	if err != nil {
		return NewErrorErr(err)
	}

	var str *String
	if find {
		v, li, find := list.GetTail()
		if find {
			str = v.(*String)
			list.Remove(li)
		}
	}

	return NewBulkStr(str)
}

func (t *StandardServer) CommandRPush(s *Session, args []*String) Message {
	if len(args) < 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	values := args[1:]

	list, _, _, err := s.db.GetListOrCreate(key)
	if err != nil {
		return NewErrorErr(err)
	}

	for _, v := range values {
		list.AddTail(v)
	}

	return NewInteger(list.Len())
}

func (t *StandardServer) CommandRPushX(s *Session, args []*String) Message {
	if len(args) < 2 {
		return ErrorInvalidArgNum
	}

	key := args[0]
	values := args[1:]

	list, _, find, err := s.db.GetList(key)
	if err != nil {
		return NewErrorErr(err)
	}

	num := 0
	if find {
		for _, v := range values {
			list.AddTail(v)
		}
		num = list.Len()
	}

	return NewInteger(num)
}

func (t *StandardServer) CommandBLPop(s *Session, args []*String) Message {
	//if len(args) != 1 {
	//	return ErrorInvalidArgNum
	//}
	//
	//keys := args[1:]
	//
	//for _, key := range keys {
	//	list, _, find, err := s.db.GetList(key)
	//	if err != nil {
	//		return NewErrorErr(err)
	//	}
	//
	//	if find && list.Len() > 0 {
	//		v, li, _ := list.GetHead()
	//		list.Remove(li)
	//
	//		return NewBulkStr(v.(*String))
	//	}
	//}
	//
	//s.flag |= SessionFlagBlocking
	//s.db.AddListBlocking(s, keys...)
	//s.AddListBlocking(keys...)

	return nil
}

//TODO blpop, brpop, brpoplpush, blmove, lindex, linsert, lpos, lrange, lrem, lset, ltrim, rpoplpush, lmove

//endregion
