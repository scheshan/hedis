package codec

import "hedis/core"

type Simple struct {
	str *core.String
}

func NewSimple() *Simple {
	s := &Simple{}
	s.str = core.NewString()
	return s
}

func (t *Simple) String() string {
	return ""
}

func (t *Simple) Read(data []byte) (int, bool, error) {
	return 0, true, nil
}
