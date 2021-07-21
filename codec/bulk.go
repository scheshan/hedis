package codec

import (
	"bufio"
	"hedis/core"
	"strconv"
)

type Bulk struct {
	length int
	str    *core.String
}

func (t *Bulk) String() string {
	return toString(t)
}

func (t *Bulk) Read(reader *bufio.Reader) error {
	num, err := readInteger(reader)
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

func (t *Bulk) Write(writer *bufio.Writer) (err error) {
	if _, err = writer.WriteString("$"); err != nil {
		return err
	}
	if _, err = writer.WriteString(strconv.Itoa(t.length)); err != nil {
		return err
	}
	if err = writeCRLF(writer); err != nil {
		return err
	}

	if t.length >= 0 {
		if _, err = writer.Write(t.str.Bytes()); err != nil {
			return err
		}

		return writeCRLF(writer)
	}

	return nil
}

func NewBulk() *Bulk {
	b := &Bulk{}

	return b
}
