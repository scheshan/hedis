package hedis

import (
	"bufio"
	"errors"
)

var decoder = &Decoder{}

type Decoder struct {
}

func (t *Decoder) Decode(reader *bufio.Reader) (Message, error) {
	msg, err := t.readMessage(reader)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (t *Decoder) readMessage(reader *bufio.Reader) (Message, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch b {
	case '+':
		msg := &Simple{}
		str, err := t.readLine(reader)
		if err != nil {
			return nil, err
		}

		msg.str = str
		return msg, nil
	case '-':
		msg := &Error{}
		str, err := t.readLine(reader)
		if err != nil {
			return nil, err
		}

		msg.str = str
		return msg, nil
	case ':':
		msg := &Integer{}
		num, err := t.readInteger(reader)
		if err != nil {
			return nil, err
		}

		msg.num = num
		return msg, nil
	case '$':
		msg := &Bulk{}
		str, err := t.readBulk(reader)
		if err != nil {
			return nil, err
		}

		msg.str = str
		return msg, nil
	case '*':
		return t.readArray(reader)
	default:
		return t.readInline(reader, b)
	}
}

func (t *Decoder) readLine(reader *bufio.Reader) (*String, error) {
	str := NewStringEmpty()

	finish := false

	for !finish {
		line, prefix, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}

		str.Append(line)
		finish = !prefix
	}

	return str, nil
}

func (t *Decoder) readSymbol(reader *bufio.Reader) (num int, negative bool, err error) {
	num = 0
	negative = false
	var b byte

	b, err = reader.ReadByte()
	if err != nil {
		return
	}

	if b == '-' {
		negative = true
	} else if b == '+' {
		negative = false
	} else if b >= '0' && b <= '9' {
		num = int(b - '0')
	} else {
		err = ErrInvalidMessage
	}

	return
}

func (t *Decoder) readInteger(reader *bufio.Reader) (int, error) {
	num, negative, err := t.readSymbol(reader)
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

			if num < 0 {
				return 0, errors.New("number overflow")
			}
		} else if b == '\r' {
			if err = t.readLF(reader); err != nil {
				return 0, err
			}

			if negative {
				num = -num
			}
			return num, nil
		} else {
			return 0, ErrInvalidMessage
		}
	}
}

func (t *Decoder) readBulk(reader *bufio.Reader) (*String, error) {
	num, err := t.readInteger(reader)
	if err != nil {
		return nil, err
	}

	var str *String
	if num < 0 {
		str = NewStringEmpty()
	} else if num == 0 {
		if err = t.readCRLF(reader); err != nil {
			return nil, err
		}
		str = NewStringEmpty()
	} else {
		str = NewStringSize(num)

		for str.Len() < num {
			require := num - str.Len()
			if require > reader.Size() {
				require = reader.Size()
			}

			peek, err := reader.Peek(require)
			if err != nil {
				return nil, err
			}

			str.Append(peek)
			_, _ = reader.Discard(len(peek))
		}

		if err = t.readCRLF(reader); err != nil {
			return nil, err
		}
	}

	return str, nil
}

func (t *Decoder) readArray(reader *bufio.Reader) (*Array, error) {
	num, err := t.readInteger(reader)
	if err != nil {
		return nil, err
	}

	if num <= 0 {
		num = 0
	}
	messages := make([]Message, 0, num)
	arr := &Array{}
	arr.messages = messages

	if num == 0 {
		return arr, nil
	}

	for num > 0 {
		msg, err := t.readMessage(reader)
		if err != nil {
			return nil, err
		}

		arr.messages = append(arr.messages, msg)
		num--

	}
	return arr, nil
}

func (t *Decoder) readInline(reader *bufio.Reader, b byte) (*Inline, error) {
	inline := &Inline{}

	arg := NewStringMinSize()
	arg.AppendByte(b)
	inline.args = append(inline.args, arg)

	finish := false

	for !finish {
		line, prefix, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(line); i++ {
			if line[i] == ' ' {
				arg = NewStringEmpty()
				inline.args = append(inline.args, arg)
			} else {
				arg.AppendByte(line[i])
			}
		}
		finish = !prefix
	}

	return inline, nil
}

func (t *Decoder) readCRLF(reader *bufio.Reader) error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if b != '\r' {
		return err
	}

	return t.readLF(reader)
}

func (t *Decoder) readLF(reader *bufio.Reader) error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if b != '\n' {
		return ErrInvalidMessage
	}

	return nil
}
