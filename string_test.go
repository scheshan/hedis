package hedis

import (
	"fmt"
	"strings"
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

func TestString_Split(t *testing.T) {
	c := "a b c d e"
	s := NewStringS(c)
	arr1 := strings.Split(c, " ")
	arr2 := s.Split(" ")

	if len(arr1) != len(arr2) {
		t.Error("invalid")
	}
	for i, str := range arr1 {
		if arr2[i].String() != str {
			t.Error("invalid")
		}
	}
}
