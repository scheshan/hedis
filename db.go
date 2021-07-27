package hedis

type Db struct {
	ht           *Hash
	expires      *Hash
	listBlocking *Hash
}

func (t *Db) Get(key *String) (*Object, bool) {
	i, ok := t.ht.Get(key)
	if !ok {
		return nil, ok
	}

	o, ok := i.(*Object)
	if !ok {
		return nil, ok
	}

	if o.Expired() {
		t.Remove(key)
		return nil, false
	}

	return o, true
}

func (t *Db) GetStringOrCreate(key *String, def *String) (*String, *Object, bool, error) {
	var res *String

	obj, find := t.Get(key)
	if !find {
		res = def
		obj = NewObject(ObjectTypeString, res)
		t.ht.Put(key, obj)
	} else {
		if obj.objType != ObjectTypeString {
			return nil, nil, find, ErrInvalidObjectType
		}
		res = obj.value.(*String)
	}

	return res, obj, find, nil
}

func (t *Db) GetString(key *String) (*String, *Object, bool, error) {
	obj, find := t.Get(key)
	if !find {
		return nil, nil, find, nil
	}

	if obj.objType != ObjectTypeString {
		return nil, nil, find, ErrInvalidObjectType
	}

	res := obj.value.(*String)
	return res, obj, find, nil
}

func (t *Db) GetHashOrCreate(key *String) (*Hash, *Object, bool, error) {
	var res *Hash

	obj, find := t.Get(key)
	if !find {
		res = NewHashSize(16)
		obj = NewObject(ObjectTypeHash, res)
		t.ht.Put(key, obj)
	} else {
		if obj.objType != ObjectTypeHash {
			return nil, nil, find, ErrInvalidObjectType
		}
		res = obj.value.(*Hash)
	}

	return res, obj, find, nil
}

func (t *Db) GetHash(key *String) (*Hash, *Object, bool, error) {
	obj, find := t.Get(key)
	if !find {
		return nil, nil, find, nil
	}

	if obj.objType != ObjectTypeHash {
		return nil, nil, find, ErrInvalidObjectType
	}

	res := obj.value.(*Hash)
	return res, obj, find, nil
}

func (t *Db) GetListOrCreate(key *String) (*List, *Object, bool, error) {
	var res *List

	obj, find := t.Get(key)
	if !find {
		res = NewList()
		obj := NewObject(ObjectTypeList, res)
		t.Put(key, obj)
	} else {
		if obj.objType != ObjectTypeList {
			return nil, nil, find, ErrInvalidObjectType
		}
		res = obj.value.(*List)
	}

	return res, obj, find, nil
}

func (t *Db) GetList(key *String) (*List, *Object, bool, error) {
	obj, find := t.Get(key)
	if !find {
		return nil, nil, find, nil
	}

	if obj.objType != ObjectTypeList {
		return nil, nil, find, ErrInvalidObjectType
	}

	res := obj.value.(*List)
	return res, obj, find, nil
}

func (t *Db) Put(key *String, obj *Object) {
	t.ht.Put(key, obj)
}

func (t *Db) Remove(key *String) bool {
	b := t.ht.Remove(key)

	t.expires.Remove(key)

	return b
}

func (t *Db) Exists(key *String) bool {
	return t.ht.Contains(key)
}

func (t *Db) AddListBlocking(s *Session, keys ...*String) {
	for _, key := range keys {
		v, find := t.listBlocking.Get(key)
		if !find {
			list := NewList()
			v = list
			t.listBlocking.Put(key, v)
		}

		v.(*List).AddHead(s)
	}
}

func (t *Db) RemoveListBlocking(s *Session, keys ...*String) {
	for _, key := range keys {
		v, _ := t.listBlocking.Get(key)
		list := v.(*List)
		list.RemoveFilter(func(v interface{}) bool {
			return v == s
		})
	}
}

func NewDb() *Db {
	db := &Db{}
	db.ht = NewHash()
	db.expires = NewHash()
	db.listBlocking = NewHash()

	return db
}
