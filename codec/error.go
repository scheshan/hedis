package codec

import (
	"bufio"
)

type Error struct {
	*Simple
}

func NewError() *Error {
	err := &Error{}
	err.Simple = NewSimple()

	return err
}

func (t *Error) Write(writer *bufio.Writer) (err error) {
	if _, err = writer.WriteString("-"); err != nil {
		return err
	}
	if _, err = writer.Write(t.str.Bytes()); err != nil {
		return err
	}

	return writeCRLF(writer)
}
