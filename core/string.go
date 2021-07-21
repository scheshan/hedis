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

func (t *String) Cap() int {
	return cap(t.buf)
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

func (t *String) Bytes() []byte {
	return t.buf
}

func NewEmptyString() *String {
	str := &String{}
	str.buf = make([]byte, 0, 0)

	return str
}

func NewMinimalString() *String {
	return NewString(16)
}

func NewString(size int) *String {
	str := &String{}
	str.buf = make([]byte, 0, size)

	return str
}

func NewStringStr(str string) *String {
	s := &String{}
	s.buf = []byte(str)

	return s
}
