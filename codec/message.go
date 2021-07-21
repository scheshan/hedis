package codec

import (
	"bufio"
	"bytes"
	"errors"
)

var InvalidMessage = errors.New("invalid message")

type Message interface {
	String() string
	Read(reader *bufio.Reader) error
	Write(writer *bufio.Writer) error
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

func readSymbol(reader *bufio.Reader) (num int, negative bool, err error) {
	num = 0
	negative = false
	var b byte

	b, err = reader.ReadByte()
	if err != nil {
		return
	}

	if b == '-' {
		negative = true
	} else if b >= '0' && b <= '9' {
		num = int(b - '0')
	} else {
		err = InvalidMessage
	}

	return
}

func readCRLF(reader *bufio.Reader) error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if b != '\r' {
		return err
	}

	return readLF(reader)
}

func readLF(reader *bufio.Reader) error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if b != '\n' {
		return InvalidMessage
	}

	return nil
}

func readInteger(reader *bufio.Reader) (res int, err error) {
	num, negative, err := readSymbol(reader)
	if err != nil {
		return 0, err
	}

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return num, err
		}

		if b >= '0' && b <= '9' {
			num = num*10 + int(b-'0')
		} else if b == '\r' {
			if err = readLF(reader); err != nil {
				return 0, err
			}

			res = num
			if negative {
				res = -res
			}
			return res, nil
		} else {
			return 0, InvalidMessage
		}
	}
}

func toString(message Message) string {
	buf := bytes.NewBuffer(make([]byte, 1024))
	writer := bufio.NewWriter(buf)

	err := message.Write(writer)
	if err != nil {
		return err.Error()
	}

	if err = writer.Flush(); err != nil {
		return err.Error()
	}

	return buf.String()
}

func writeCRLF(writer *bufio.Writer) error {
	_, err := writer.WriteString("\r\n")

	return err
}
