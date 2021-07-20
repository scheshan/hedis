package codec

import (
	"bufio"
	"errors"
)

var InvalidMessage = errors.New("invalid message")

type Message interface {
	String() string
	Read(reader *bufio.Reader) (bool, error)
}

func ReadMessage(reader *bufio.Reader) (Message, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

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

	return 0
}
