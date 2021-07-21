package codec

import (
	"bufio"
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

func NewSimpleStr(str string) *Simple {
	s := &Simple{}
	s.str = core.NewStringStr(str)

	return s
}

func (t *Simple) String() string {
	return toString(t)
}

func (t *Simple) Read(reader *bufio.Reader) error {
	finish := false
	for !finish {
		line, prefix, err := reader.ReadLine()
		if err != nil {
			return err
		}

		t.str.Append(line)
		finish = !prefix
	}

	return nil
}

func (t *Simple) Write(writer *bufio.Writer) (err error) {
	if _, err = writer.WriteString("+"); err != nil {
		return err
	}
	if _, err = writer.Write(t.str.Bytes()); err != nil {
		return err
	}

	return writeCRLF(writer)
}
