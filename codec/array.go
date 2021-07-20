package codec

import "bufio"

type Array struct {
	messages []Message
}

func (t *Array) String() string {
	return ""
}

func (t *Array) Read(reader *bufio.Reader) (bool, error) {
	return true, nil
}

func NewArray() *Array {
	arr := &Array{}

	return arr
}
