package hedis

import "time"

type ObjectType int

const (
	ObjectTypeString ObjectType = iota
	ObjectTypeList
	ObjectTypeHash
)

type Object struct {
	objType ObjectType
	value   interface{}
	ttl     int64
}

func (t *Object) Expired() bool {
	return t.ttl > 0 && time.Now().UnixNano() > (t.ttl*1e6)
}

func (t *Object) ExpireAt(sec int64) {
	t.PExpireAt(sec * 1e3)
}

func (t *Object) PExpireAt(ms int64) {
	t.ttl = ms * 1e6
}

func NewObject(t ObjectType, v interface{}) *Object {
	obj := &Object{}
	obj.objType = t
	obj.value = v

	return obj
}
