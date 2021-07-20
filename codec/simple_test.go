package codec

import "testing"

func TestSimple_Read(t *testing.T) {
	simple := NewSimple()

	reads, b, err := simple.Read([]byte("hello\r\n"))
	if err != nil {
		t.Failed()
	}
	if !b {
		t.Failed()
	}
	if reads != 7 {
		t.Failed()
	}
}

func TestSimple_HarfString(t *testing.T) {
	simple := NewSimple()

	reads, b, _ := simple.Read([]byte("hello"))
	if b {
		t.Failed()
	}
	if reads != 5 {
		t.Failed()
	}

	reads, b, _ = simple.Read([]byte("world\r\n"))
	if !b {
		t.Failed()
	}
	if reads != 7 {
		t.Failed()
	}

	if simple.String() != "helloworld" {
		t.Failed()
	}
}
