package codec

import (
	"bufio"
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

func (t *Integer) Value() int {
	return t.num
}

func (t *Integer) Read(reader *bufio.Reader) (bool, error) {
	line, prefix, err := reader.ReadLine()
	if err != nil {
		return false, err
	}

	t.str.Append(line)

	if prefix {
		return false, nil
	}

	t.num, err = strconv.Atoi(t.str.String())
	return err != nil, err
}

func NewInteger() *Integer {
	i := &Integer{}
	i.str = core.NewMinimalString()

	return i
}
