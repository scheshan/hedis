package codec

import (
	"bufio"
	"hedis/core"
	"strconv"
)

const (
	bulkStateReadLength = 1
	bulkStateReadString = 2
)

type Bulk struct {
	state  int
	length *Integer
	str    *core.String
}

func (t *Bulk) String() string {
	if t.state == bulkStateReadLength {
		return ""
	}

	return t.str.String()
}

func (t *Bulk) Read(reader *bufio.Reader) (bool, error) {
	if t.state == bulkStateReadLength {
		finish, err := t.length.Read(reader)
		if err != nil {
			return false, err
		}

		if !finish {
			return false, nil
		}

		t.state = bulkStateReadString
		size := t.length.num
		if size < 0 {
			size = 0
		}
		t.str = core.NewString(size)

		return false, nil
	}

	reader.read
}

func (t *Bulk) Read(data []byte) (int, bool, error) {
	if t.state == bulkStateReadLength {
		//读取长度
		reads := ReadCRLF(data)
		if reads == 0 {
			t.str.Append(data)
			return reads, false, nil
		}

		t.str.Append(data[0:reads])
		var err error
		t.length, err = strconv.Atoi(t.str.String())
		if err != nil {
			return reads + 2, false, err
		}

		if t.str.Cap() >= t.length {
			t.str.Clear()
		} else {
			t.str = core.NewString(t.length)
		}

		t.state = bulkStateReadString
		return reads + 2, false, nil
	} else if t.length-t.str.Len() > 0 {
		//读取文本
		remain := t.length - t.str.Len()

		if len(data) < remain {
			t.str.Append(data[0:])

			return len(data), false, nil
		}

		t.str.Append(data[0:t.length])
		return t.length, false, nil
	} else {
		//读取CRLF
		reads := ReadCRLF(data)
		if reads == 0 {

		}
	}
}

func NewBulk() *Bulk {
	b := &Bulk{}
	b.state = bulkStateReadLength
	b.length = NewInteger()

	return b
}
