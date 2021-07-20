package codec

import "testing"

func TestBulk_Read(t *testing.T) {
	b := NewBulk()

	b.Read([]byte("5"))
	b.Read([]byte("\r\n"))
	b.Read([]byte("hello"))
}
