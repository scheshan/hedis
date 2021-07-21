package codec

import (
	"bufio"
	"strconv"
)

type Integer struct {
	num int
}

func (t *Integer) String() string {
	return toString(t)
}

func (t *Integer) Value() int {
	return t.num
}

func (t *Integer) Read(reader *bufio.Reader) error {
	num, err := readInteger(reader)
	if err != nil {
		return err
	}

	t.num = num
	return nil
}

func (t *Integer) Write(writer *bufio.Writer) (err error) {
	if _, err = writer.WriteString(":"); err != nil {
		return err
	}
	if _, err = writer.WriteString(strconv.Itoa(t.num)); err != nil {
		return err
	}
	return writeCRLF(writer)
}

func NewInteger() *Integer {
	i := &Integer{}

	return i
}
