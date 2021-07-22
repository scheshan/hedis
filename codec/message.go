package codec

import (
	"bufio"
	"bytes"
	"errors"
	"hedis/core"
	"strconv"
)

var InvalidMessage = errors.New("invalid message")

type Message interface {
	Command() *core.String
	Args() []*core.String
	ToString() *core.String
}

type Simple struct {
	str *core.String
}

func (t *Simple) Command() *core.String {
	return t.str
}

func (t *Simple) Args() []*core.String {
	return nil
}

func (t *Simple) ToString() *core.String {
	return t.str
}

type Integer struct {
	num int
}

func (t *Integer) Command() *core.String {
	return nil
}

func (t *Integer) Args() []*core.String {
	return nil
}

func (t *Integer) ToString() *core.String {
	return core.NewStringStr(strconv.Itoa(t.num))
}

type Error struct {
	str *core.String
}

func (t *Error) Command() *core.String {
	return nil
}

func (t *Error) Args() []*core.String {
	return nil
}

func (t *Error) ToString() *core.String {
	return t.str
}

type Bulk struct {
	str *core.String
}

func (t *Bulk) Command() *core.String {
	return t.str
}

func (t *Bulk) Args() []*core.String {
	return nil
}

func (t *Bulk) ToString() *core.String {
	return t.str
}

type Array struct {
	messages []Message
}

func (t *Array) Command() *core.String {
	if len(t.messages) == 0 {
		return nil
	}

	msg, ok := t.messages[0].(*Bulk)
	if !ok {
		return nil
	}

	return msg.str
}

func (t *Array) Args() []*core.String {
	args := make([]*core.String, len(t.messages)-1, len(t.messages)-1)

	for i := 1; i < len(t.messages); i++ {
		args[i-1] = t.messages[i].ToString()
	}

	return args
}

func (t *Array) ToString() *core.String {
	return nil
}

type Inline struct {
	args []*core.String
}

func (t *Inline) Command() *core.String {
	if len(t.args) > 0 {
		return t.args[0]
	}

	return nil
}

func (t *Inline) Args() []*core.String {
	return t.args[1:]
}

func (t *Inline) ToString() *core.String {
	return nil
}

func Decode(reader *bufio.Reader) (Message, error) {
	return decoder.Decode(reader)
}

func Encode(writer *bufio.Writer, message Message) error {
	return encoder.Encode(writer, message)
}

func EncodeString(message Message) (string, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	writer := bufio.NewWriter(buffer)

	if err := Encode(writer, message); err != nil {
		return "", err
	}

	if err := writer.Flush(); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func NewSimpleString(text string) *Simple {
	res := &Simple{}
	res.str = core.NewStringStr(text)

	return res
}

func NewSimpleStr(str *core.String) *Simple {
	res := &Simple{}
	res.str = str

	return res
}

func NewErrorString(text string) *Error {
	res := &Error{}
	res.str = core.NewStringStr(text)

	return res
}

func NewBulkString(text string) *Bulk {
	str := core.NewStringStr(text)
	return NewBulkStr(str)
}

func NewBulkStr(str *core.String) *Bulk {
	bulk := &Bulk{}
	bulk.str = str

	return bulk
}
