package codec

import (
	"errors"
)

var InvalidMessage = errors.New("invalid message")

type Message interface {
	String() string
	Read(data []byte) (int, bool, error)
}

func ReadMessage(data []byte) (Message, error) {
	if len(data) == 0 {
		return nil, InvalidMessage
	}

	b := data[0]

	switch b {
	case '+':
		return NewSimple(), nil
	case '-':
		return NewError(), nil
	case ':':
		return NewInteger(), nil
	case '$':
		return NewBulk(), nil
	case '*':
		return NewArray(), nil
	}

	return nil, InvalidMessage
}

func ReadCRLF(data []byte) int {
	for i := 0; i < len(data)-2; i++ {
		if data[i+1] == '\r' && data[i+2] == '\n' {
			return i + 1
		}
	}

	return -1
}
