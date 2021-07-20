package codec

import "hedis/core"

type Bulk struct {
	str *core.String
}

func (t *Bulk) String() string {
	return ""
}

func (t *Bulk) Read(data []byte) (int, bool, error) {
	return 0, true, nil
}

func NewBulk() *Bulk {
	b := &Bulk{}

	return b
}
