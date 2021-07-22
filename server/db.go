package server

import "hedis/core"

type Db struct {
	ht *core.Hash
}

func NewDb() *Db {
	db := &Db{}
	db.ht = core.NewHashSize(16)

	return db
}
