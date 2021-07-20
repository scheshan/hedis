package codec

import (
	"hedis/core"
	"strconv"
)

type Bulk struct {
	lenStr *core.String
	length int
	str    *core.String
}

func (t *Bulk) String() string {
	return ""
}

func (t *Bulk) Read(data []byte) (int, bool, error) {
	reads := ReadCRLF(data)
	if reads == 0 {
		t.lenStr.Append(data)
		return reads, false, nil
	}

	t.lenStr.Append(data[0:reads])

	var err error
	t.length, err = strconv.Atoi(t.lenStr.String())
	if err != nil {
		return reads + 2, false, err
	}

	t.str = core.NewString(t.length)

	ind := reads + 2
	if len(data)-ind < t.length {
		t.str.Append(data[ind:])

		return len(data), false, nil
	}

	t.str.Append(data[ind : ind+t.length])
	return ind + t.length, true, nil
}

func NewBulk() *Bulk {
	b := &Bulk{}
	b.lenStr = core.NewEmptyString()

	return b
}
