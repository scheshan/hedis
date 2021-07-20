package codec

import (
	"hedis/core"
)

type Simple struct {
	str *core.String
}

func NewSimple() *Simple {
	s := &Simple{}
	s.str = core.NewMinimalString()

	return s
}

func (t *Simple) String() string {
	return t.str.String()
}

func (t *Simple) Read(data []byte) (int, bool, error) {
	reads := ReadCRLF(data)

	if reads == 0 {
		t.str.Append(data)
		return len(data), false, nil
	}

	t.str.Append(data[0:reads])
	return reads + 2, true, nil
}
