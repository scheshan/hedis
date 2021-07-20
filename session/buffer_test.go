package session

import (
	"testing"
)

func TestBuffer_allocSize(t *testing.T) {
	buf := &Buffer{}

	size := buf.allocSize(5)
	if size != MinimalBufferSize {
		t.Fail()
	}

	size = buf.allocSize(16)
	if size != 16 {
		t.Fail()
	}

	size = buf.allocSize(17)
	if size != 32 {
		t.Fail()
	}

	size = buf.allocSize(1048577)
	if size != 1048577+1048576 {
		t.Fail()
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Fail()
		}
	}()

	size = buf.allocSize(1 << 32)
}

func TestBuffer_grow(t *testing.T) {
	buf := NewBuffer(MinimalBufferSize)
	buf.Append(make([]byte, 5, 5))
	buf.Append(make([]byte, 17, 17))

	free := buf.RealFree()
	if free != 10 {
		t.Fail()
	}

	buf.ReadType()
	buf.Append(make([]byte, 11, 11))

	free = buf.RealFree()
	if free != 0 {
		t.Fail()
	}
}
