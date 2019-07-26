package hedis

import (
	"fmt"
	"testing"
)

func TestString_Slice(t *testing.T) {
	c := "abcdefghijklmnopqrstuvwxyz"

	s := NewStringS(c)

	if s.Slice(0).String() != c {
		t.Error("invalid")
	}
	if s.Slice(26).String() != "" {
		t.Error("invalid")
	}
	if s.Slice(-1).String() != "" {
		t.Error("invalid")
	}
}

func TestString_SliceLength(t *testing.T) {
	c := "abcdefghijklmnopqrstuvwxyz"

	s := NewStringS(c)

	if s.SliceLength(0, 26).String() != c {
		t.Error("invalid")
	}
	if s.SliceLength(5, 5).String() != "fghij" {
		t.Error("invalid")
	}
	if s.SliceLength(-1, 5).String() != "" {
		t.Error("invalid")
	}
	if s.SliceLength(5, -1).String() != "" {
		t.Error("invalid")
	}
	fmt.Println(s.SliceLength(5, 30).String())
	if s.SliceLength(5, 30).String() != "" {
		t.Error("invalid")
	}
}
