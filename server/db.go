package server

import (
	"hedis/core"
)

type Db struct {
	ht           *core.Hash
	listBlocking *core.Hash
}

func (t *Db) Get(key *core.String) (*Object, bool) {
	i, ok := t.ht.Get(key)
	if !ok {
		return nil, ok
	}

	o, ok := i.(*Object)
	if !ok {
		return nil, ok
	}

	return o, true
}

func (t *Db) GetStringOrCreate(key *core.String, def *core.String) (*core.String, *Object, bool, error) {
	var res *core.String

	obj, find := t.Get(key)
	if !find {
		res = def
		obj = NewObject(ObjectTypeString, res)
		t.ht.Put(key, obj)
	} else {
		if obj.objType != ObjectTypeString {
			return nil, nil, find, ErrorInvalidObjectType
		}
		res = obj.value.(*core.String)
	}

	return res, obj, find, nil
}

func (t *Db) GetString(key *core.String) (*core.String, *Object, bool, error) {
	obj, find := t.Get(key)
	if !find {
		return nil, nil, find, nil
	}

	if obj.objType != ObjectTypeString {
		return nil, nil, find, ErrorInvalidObjectType
	}

	res := obj.value.(*core.String)
	return res, obj, find, nil
}

func (t *Db) GetHashOrCreate(key *core.String) (*core.Hash, *Object, bool, error) {
	var res *core.Hash

	obj, find := t.Get(key)
	if !find {
		res = core.NewHashSize(16)
		obj = NewObject(ObjectTypeHash, res)
		t.ht.Put(key, obj)
	} else {
		if obj.objType != ObjectTypeHash {
			return nil, nil, find, ErrorInvalidObjectType
		}
		res = obj.value.(*core.Hash)
	}

	return res, obj, find, nil
}

func (t *Db) GetHash(key *core.String) (*core.Hash, *Object, bool, error) {
	obj, find := t.Get(key)
	if !find {
		return nil, nil, find, nil
	}

	if obj.objType != ObjectTypeHash {
		return nil, nil, find, ErrorInvalidObjectType
	}

	res := obj.value.(*core.Hash)
	return res, obj, find, nil
}

func (t *Db) GetListOrCreate(key *core.String) (*core.List, *Object, bool, error) {
	var res *core.List

	obj, find := t.Get(key)
	if !find {
		res = core.NewList()
		obj := NewObject(ObjectTypeList, res)
		t.Put(key, obj)
	} else {
		if obj.objType != ObjectTypeList {
			return nil, nil, find, ErrorInvalidObjectType
		}
		res = obj.value.(*core.List)
	}

	return res, obj, find, nil
}

func (t *Db) GetList(key *core.String) (*core.List, *Object, bool, error) {
	obj, find := t.Get(key)
	if !find {
		return nil, nil, find, nil
	}

	if obj.objType != ObjectTypeList {
		return nil, nil, find, ErrorInvalidObjectType
	}

	res := obj.value.(*core.List)
	return res, obj, find, nil
}

func (t *Db) Put(key *core.String, obj *Object) {
	t.ht.Put(key, obj)
}

func (t *Db) Remove(key *core.String) bool {
	return t.ht.Remove(key)
}

func (t *Db) Exists(key *core.String) bool {
	return t.ht.Contains(key)
}

func (t *Db) AddListBlocking(s *Session, keys ...*core.String) {
	for _, key := range keys {
		v, find := t.listBlocking.Get(key)
		if !find {
			list := core.NewList()
			v = list
			t.listBlocking.Put(key, v)
		}

		v.(*core.List).AddHead(s)
	}
}

func (t *Db) RemoveListBlocking(s *Session, keys ...*core.String) {
	for _, key := range keys {
		v, _ := t.listBlocking.Get(key)
		list := v.(*core.List)
		list.RemoveFilter(func(v interface{}) bool {
			return v == s
		})
	}
}

func NewDb() *Db {
	db := &Db{}
	db.ht = core.NewHash()
	db.listBlocking = core.NewHash()

	return db
}
