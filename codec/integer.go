package codec

import (
	"bufio"
	"strconv"
)

type Integer struct {
	num int
}

func (t *Integer) String() string {
	return strconv.Itoa(t.num)
}

func (t *Integer) Value() int {
	return t.num
}

func (t *Integer) Read(reader *bufio.Reader) error {
	num, err := ReadInteger(reader)
	if err != nil {
		return err
	}

	t.num = num
	return nil
}

func NewInteger() *Integer {
	i := &Integer{}

	return i
}
