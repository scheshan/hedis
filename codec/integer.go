package codec

import (
	"hedis/core"
	"strconv"
)

type Integer struct {
	str *core.String
	num int
}

func (t *Integer) String() string {
	return ""
}

func (t *Integer) Read(data []byte) (int, bool, error) {
	reads := ReadCRLF(data)
	if reads == 0 {
		t.str.Append(data)
		return reads, false, nil
	}

	t.str.Append(data[0:reads])

	var err error
	t.num, err = strconv.Atoi(t.str.String())

	if err != nil {
		return reads + 2, false, err
	}

	return reads + 2, true, err
}

func NewInteger() *Integer {
	i := &Integer{}
	i.str = core.NewEmptyString()

	return i
}
