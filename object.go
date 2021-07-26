package hedis

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

func NewObject(t ObjectType, v interface{}) *Object {
	obj := &Object{}
	obj.objType = t
	obj.value = v

	return obj
}
