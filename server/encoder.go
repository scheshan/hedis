package server

import (
	"bufio"
	"errors"
	"strconv"
)

var encoder = &Encoder{}

type Encoder struct {
}

func (t *Encoder) Encode(writer *bufio.Writer, message Message) error {
	var err error
	if msg, ok := message.(*Simple); ok {
		return t.encodeLine(writer, "+", msg.str)
	} else if msg, ok := message.(*Error); ok {
		return t.encodeLine(writer, "-", msg.str)
	} else if msg, ok := message.(*Integer); ok {
		return t.encodeInteger(writer, msg)
	} else if msg, ok := message.(*Bulk); ok {
		return t.encodeBulk(writer, msg)
	} else if msg, ok := message.(*Array); ok {
		return t.encodeArray(writer, msg)
	} else if msg, ok := message.(*Inline); ok {
		return t.encodeInline(writer, msg)
	} else {
		return errors.New("message not supported")
	}

	return err
}

func (t *Encoder) encodeLine(writer *bufio.Writer, prefix string, str *String) error {
	_, err := writer.WriteString(prefix)
	if err != nil {
		return err
	}

	if _, err := writer.Write(str.Bytes()); err != nil {
		return err
	}

	return t.encodeCRLF(writer)
}

func (t *Encoder) encodeInteger(writer *bufio.Writer, msg *Integer) error {
	_, err := writer.WriteString(":")
	if err != nil {
		return err
	}

	return t.encodeNumber(writer, msg.num)
}

func (t *Encoder) encodeBulk(writer *bufio.Writer, msg *Bulk) error {
	_, err := writer.WriteString("$")
	if err != nil {
		return err
	}

	len := -1
	if msg.str != nil {
		len = msg.str.Len()
	}

	if err = t.encodeNumber(writer, len); err != nil {
		return err
	}

	if len > -1 {
		if _, err := writer.Write(msg.str.Bytes()); err != nil {
			return err
		}
		return t.encodeCRLF(writer)
	}

	return nil
}

func (t *Encoder) encodeArray(writer *bufio.Writer, msg *Array) error {
	_, err := writer.WriteString("*")
	if err != nil {
		return err
	}

	if err = t.encodeNumber(writer, len(msg.messages)); err != nil {
		return err
	}

	for _, m := range msg.messages {
		if err = t.Encode(writer, m); err != nil {
			return err
		}
	}

	return nil
}

func (t *Encoder) encodeInline(writer *bufio.Writer, msg *Inline) error {
	for i, m := range msg.args {
		if i > 0 {
			if _, err := writer.WriteString(" "); err != nil {
				return err
			}
		}

		if _, err := writer.Write(m.Bytes()); err != nil {
			return err
		}
	}

	return t.encodeCRLF(writer)
}

func (t *Encoder) encodeNumber(writer *bufio.Writer, num int) error {
	_, err := writer.WriteString(strconv.Itoa(num))
	if err != nil {
		return err
	}

	return t.encodeCRLF(writer)
}

func (t *Encoder) encodeCRLF(writer *bufio.Writer) error {
	_, err := writer.WriteString("\r\n")
	return err
}
