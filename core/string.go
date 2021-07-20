package core

type String struct {
	buf []byte
}

func (t *String) Clear() {
	t.buf = t.buf[0:0]
}

func (t *String) Len() int {
	return len(t.buf)
}

func (t *String) Append(data []byte) {
	t.buf = append(t.buf, data...)
}

func (t *String) AppendStr(str *String) {
	t.Append(str.buf)
}

func (t *String) String() string {
	return string(t.buf)
}

func NewEmptyString() *String {
	str := &String{}
	str.buf = make([]byte, 16, 16)

	return str
}

func NewString(size int) *String {
	str := &String{}
	str.buf = make([]byte, size, size)

	return str
}
