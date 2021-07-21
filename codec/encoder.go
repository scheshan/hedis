package codec

import (
	"bufio"
	"errors"
	"hedis/core"
	"strconv"
)

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
	} else if _, ok := message.(*Array); ok {
		_, err = writer.WriteString("*")
	} else {
		err = errors.New("message not supported")
	}

	return err
}

func (t *Encoder) encodeLine(writer *bufio.Writer, prefix string, str *core.String) error {
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

	if err = t.encodeNumber(writer, msg.str.Len()); err != nil {
		return err
	}

	if _, err := writer.Write(msg.str.Bytes()); err != nil {
		return err
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
