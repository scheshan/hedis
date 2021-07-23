package server

import (
	"hedis/core"
)

type Db struct {
	ht *core.Hash
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

func (t *Db) Put(key *core.String, obj *Object) error {
	return t.ht.Put(key, obj)
}

func (t *Db) Remove(key *core.String) bool {
	return t.ht.Remove(key)
}

func (t *Db) Exists(key *core.String) bool {
	return t.ht.Contains(key)
}

func NewDb() *Db {
	db := &Db{}
	db.ht = core.NewHashSize(16)

	return db
}