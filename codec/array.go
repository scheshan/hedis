package codec

import "bufio"

type Array struct {
	length   int
	messages []Message
}

func (t *Array) String() string {
	return ""
}

func (t *Array) Read(reader *bufio.Reader) error {
	num, err := ReadInteger(reader)
	if err != nil {
		return err
	}

	t.length = num
	t.messages = make([]Message, t.length, t.length)

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

func NewArray() *Array {
	arr := &Array{}

	return arr
}
