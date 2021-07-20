package session

import "hedis/codec"

const MinimalBufferSize = 1 << 4
const MaximumBufferSize = 1 << 32

type Buffer struct {
	buf   []byte
	start int
	end   int
	last  int
}

func NewBuffer(size int) *Buffer {
	buffer := &Buffer{}

	s := buffer.allocSize(size)
	buffer.buf = make([]byte, s, s)

	return buffer
}

func (t *Buffer) allocSize(require int) int {
	if require <= MinimalBufferSize {
		return MinimalBufferSize
	}

	if require <= (1 << 20) {
		size := MinimalBufferSize
		for i := 1; i <= 16; i++ {
			size <<= 1
			if size >= require {
				return size
			}
		}
	}

	size := require + (1 << 20)
	if size < 0 || size > MaximumBufferSize {
		panic("Buffer.allocSize caused out of memory")
	}

	return size
}

func (t *Buffer) grow(require int) {
	size := t.allocSize(require)
	buf := make([]byte, size, size)

	ind := 0
	for i := t.start; i < t.end; i++ {
		buf[ind] = t.buf[i]
		ind++
	}

	t.buf = buf
	t.setStart(0)
	t.end = ind
}

func (t *Buffer) copy(data []byte) {
	ind := t.end
	for i := 0; i < len(data); i++ {
		t.buf[ind] = data[i]
		ind++
	}

	t.end = ind
}

func (t *Buffer) setStart(start int) {
	t.start = start
	t.last = start
}

func (t *Buffer) Compact() {
	if t.start == 0 {
		return
	}

	ind := 0
	for i := t.start; i < t.end; i++ {
		t.buf[ind] = t.buf[i]
		ind++
	}

	t.setStart(0)
	t.end = ind
}

func (t *Buffer) Len() int {
	return t.end - t.start
}

func (t *Buffer) RealFree() int {
	return cap(t.buf) - t.Len()
}

func (t *Buffer) Free() int {
	return cap(t.buf) - t.end
}

func (t *Buffer) Append(data []byte) {
	if t.Free() < len(data) {
		if t.RealFree() >= len(data) {
			t.Compact()
		} else {
			require := t.Len() + len(data)
			t.grow(require)
		}
	}

	t.copy(data)
}

func (t *Buffer) ReadType() (mt codec.MessageType, b bool) {
	mt = codec.MessageTypeUnknown
	b = false

	if t.Len() < 1 {
		return
	}

	data := t.buf[t.start]
	t.start++

	b = true
	switch data {
	case '+':
		mt = codec.MessageTypeString
		break
	case '-':
		mt = codec.MessageTypeError
		break
	case ':':
		mt = codec.MessageTypeInteger
		break
	case '$':
		mt = codec.MessageTypeBulk
		break
	case '*':
		mt = codec.MessageTypeArray
		break
	}

	return
}

func (t *Buffer) ReadCRLF() (data []byte, b bool) {
	ind := t.last
	for ind < t.end-1 {
		if t.buf[ind] == '\r' && t.buf[ind+1] == '\n' {
			data = t.buf[t.start:ind]
			t.setStart(ind + 2)
			return data, true
		}
		ind++
	}

	return nil, false
}
