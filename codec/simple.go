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

func (t *Simple) String() string {
	return t.str.String()
}

func (t *Simple) Read(reader *bufio.Reader) (bool, error) {
	line, prefix, err := reader.ReadLine()
	if err != nil {
		return false, err
	}

	t.str.Append(line)
	return !prefix, nil
}
