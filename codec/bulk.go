package codec

import (
	"bufio"
	"hedis/core"
)

type Bulk struct {
	length int
	str    *core.String
}

func (t *Bulk) String() string {
	return t.str.String()
}

func (t *Bulk) Read(reader *bufio.Reader) error {
	num, err := ReadInteger(reader)
	if err != nil {
		return err
	}

	t.length = num

	if t.length < 0 {
		t.str = core.NewEmptyString()
		return nil
	} else if t.length == 0 {
		return readCRLF(reader)
	} else {
		t.str = core.NewString(t.length)
	}

	for t.str.Len() < t.length {
		require := t.length - t.str.Len()
		if require > reader.Size() {
			require = reader.Size()
		}

		peek, err := reader.Peek(require)
		if err != nil {
			return err
		}

		t.str.Append(peek)
		_, _ = reader.Discard(len(peek))
	}

	err = readCRLF(reader)
	return err
}

func NewBulk() *Bulk {
	b := &Bulk{}

	return b
}
