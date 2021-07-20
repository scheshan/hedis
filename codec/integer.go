package codec

import "hedis/core"

type Integer struct {
	str *core.String
}

func (t *Integer) String() string {
	return ""
}

func (t *Integer) Read(data []byte) (int, bool, error) {
	return 0, true, nil
}

func NewInteger() *Integer {
	i := &Integer{}
	i.str = core.NewEmptyString()

	return i
}
