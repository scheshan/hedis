package codec

import (
	"bufio"
	"strconv"
)

type Array struct {
	length   int
	messages []Message
}

func (t *Array) String() string {
	return toString(t)
}

func (t *Array) Read(reader *bufio.Reader) error {
	num, err := readInteger(reader)
	if err != nil {
		return err
	}

	t.length = num
	t.messages = make([]Message, 0, t.length)

	if t.length == 0 {
		return nil
	}

	ind := 0
	for ind < t.length {
		msg, err := ReadMessage(reader)
		if err != nil {
			return err
		}

		if err = msg.Read(reader); err != nil {
			return err
		}

		t.messages = append(t.messages, msg)
		ind++
	}

	return nil
}

func (t *Array) Write(writer *bufio.Writer) (err error) {
	if _, err = writer.WriteString("*"); err != nil {
		return err
	}
	if _, err = writer.WriteString(strconv.Itoa(t.length)); err != nil {
		return err
	}
	if _, err = writer.WriteString("\r\n"); err != nil {
		return err
	}

	for _, message := range t.messages {
		if err = message.Write(writer); err != nil {
			return err
		}
	}

	return nil
}

func NewArray() *Array {
	arr := &Array{}

	return arr
}
