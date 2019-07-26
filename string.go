package hedis

import (
	"bytes"
	"strings"
)

type String struct {
	buf []byte
}

func NewStringS(str string) *String {
	s := new(String)
	s.buf = []byte(str)

	return s
}

func NewStringB(buf []byte) *String {
	s := new(String)
	s.buf = buf

	return s
}

func NewEmptyString() *String {
	s := new(String)
	return s
}

func (t *String) String() string {
	return string(t.buf)
}

func (t *String) Append(buf []byte) {
	if t.buf != nil {
		t.buf = append(t.buf, buf...)
	} else {
		t.buf = make([]byte, len(buf), len(buf))
		copy(t.buf, buf)
	}
}

func (t *String) Split(sep string) []*String {
	content := t.String()
	arr := strings.Split(content, sep)

	r := make([]*String, len(arr), len(arr))
	for i, str := range arr {
		r[i] = NewStringS(str)
	}

	return r
}

func (t *String) SplitB(sep []byte) []*String {
	return t.Split(string(sep))
}

func (t *String) SliceLength(start, length int) *String {
	if start < 0 || start >= len(t.buf) {
		return NewEmptyString()
	}
	if length < 0 {
		return NewEmptyString()
	}

	end := start + length
	if end > len(t.buf) {
		return NewEmptyString()
	}

	buf := make([]byte, length, length)
	copy(buf, t.buf[start:end+1])

	return NewStringB(buf)
}

func (t *String) Slice(start int) *String {

	return t.SliceLength(start, len(t.buf)-start)
}

func (t *String) Index(str string) int {
	return strings.Index(t.String(), str)
}

func (t *String) Len() int {
	return len(t.buf)
}

func (t *String) Equals(str string) bool {
	r := bytes.Runes([]byte(str))

	if len(t.buf) != len(r) {
		return false
	}

	for i, b := range t.buf {
		if b != byte(r[i]) {
			return false
		}
	}

	return true
}
