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
	return t.ttl > 0 && time.Now().Unix() > t.ttl
}

func NewObject(t ObjectType, v interface{}) *Object {
	obj := &Object{}
	obj.objType = t
	obj.value = v

	return obj
}
