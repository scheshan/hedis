package core

type String struct {
	hash int
	buf  []byte
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

func (t *String) AppendByte(b byte) {
	t.buf = append(t.buf, b)
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

func (t *String) HashCode() int {
	h := t.hash
	if h == 0 && t.Len() > 0 {
		for _, b := range t.buf {
			h = 31*h + int(b)
		}
	}

	return h
}

func (t *String) Equal(o *String) bool {
	if t.Len() == o.Len() {
		for i := 0; i < t.Len(); i++ {
			if t.buf[i] != o.buf[i] {
				return false
			}
		}

		return true
	}

	return false
}

func NewStringEmpty() *String {
	str := &String{}
	str.buf = make([]byte, 0, 0)

	return str
}

func NewStringMinSize() *String {
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
