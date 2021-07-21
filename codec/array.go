package codec

import "bufio"

type Array struct {
	messages []Message
}

func (t *Array) String() string {
	return ""
}

func (t *Array) Read(reader *bufio.Reader) error {
	return nil
}

func NewArray() *Array {
	arr := &Array{}

	return arr
}
