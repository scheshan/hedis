package hedis

import (
	"bufio"
)

var codecStringPrefix = []byte("+")
var codecErrorPrefix = []byte("-")
var codecCRLF = []byte("\r\n")

type Message interface {
	Write(writer *bufio.Writer) (int, error)
}

type StringMessage struct {
	content *String
}

func (t *StringMessage) Write(writer *bufio.Writer) (int, error) {
	return writerWrite(writer, codecStringPrefix, t.content.buf, codecCRLF)
}

func NewStringMessage(content *String) *StringMessage {
	sm := new(StringMessage)
	sm.content = content
	return sm
}

type ErrorMessage struct {
	err error
}

func (t *ErrorMessage) Write(writer *bufio.Writer) (int, error) {
	return writerWrite(writer, codecErrorPrefix, []byte(t.err.Error()), codecCRLF)
}

func NewErrorMessage(err error) *ErrorMessage {
	em := &ErrorMessage{
		err: err,
	}
	return em
}

func writerWrite(writer *bufio.Writer, data ...[]byte) (int, error) {
	if data == nil || len(data) == 0 {
		return 0, nil
	}

	n := 0
	tn := 0
	var err error
	for i := range data {
		tn = 0
		if tn, err = writer.Write(data[i]); err != nil {
			return 0, err
		}
		n += tn
	}

	if err = writer.Flush(); err != nil {
		return 0, err
	}

	return n, nil
}
