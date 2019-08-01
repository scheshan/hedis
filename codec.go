package hedis

import (
	"bufio"
	"fmt"
	"strconv"
)

var codecStringPrefix = []byte("+")
var codecErrorPrefix = []byte("-")
var codecNumberPrefix = []byte(":")
var codecBatchPrefix = []byte("$")
var codecArrayPrefix = []byte("*")
var codecCRLF = []byte("\r\n")

type Message interface {
	Write(writer *bufio.Writer) (int, error)
}

type StringMessage struct {
	content string
}

func (t *StringMessage) Write(writer *bufio.Writer) (int, error) {
	return writerWrite(writer, codecStringPrefix, []byte(t.content), codecCRLF)
}

func NewStringMessage(content string) *StringMessage {
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

type NumberMessage struct {
	value float64
}

func (t *NumberMessage) Write(writer *bufio.Writer) (int, error) {
	str := fmt.Sprintf("%v", t.value)

	return writerWrite(writer, codecNumberPrefix, []byte(str), codecCRLF)
}

func NewNumberMessage(value float64) *NumberMessage {
	nm := &NumberMessage{
		value: value,
	}
	return nm
}

type BatchMessage struct {
	str *String
}

func (t *BatchMessage) Write(writer *bufio.Writer) (int, error) {
	l := strconv.Itoa(t.str.Len())

	return writerWrite(writer, codecBatchPrefix, []byte(l), codecCRLF, t.str.buf, codecCRLF)
}

func NewBatchMessage(str *String) *BatchMessage {
	bm := &BatchMessage{
		str: str,
	}
	return bm
}

type ArrayMessage struct {
	children []*BatchMessage
}

func (t *ArrayMessage) Write(writer *bufio.Writer) (int, error) {
	l := strconv.Itoa(len(t.children))

	total := 0
	n := 0
	var err error

	total, err = writerWrite(writer, codecArrayPrefix, []byte(l), codecCRLF)
	if err != nil {
		return 0, err
	}

	for i := range t.children {
		n, err = t.children[i].Write(writer)
		if err != nil {
			return 0, err
		}

		total += n
	}

	return total, nil
}

func (t *ArrayMessage) Append(msg *BatchMessage) {
	t.children = append(t.children, msg)
}

func NewArrayMessage(children ...*BatchMessage) *ArrayMessage {
	am := &ArrayMessage{
		children: children,
	}
	return am
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

	return n, nil
}
