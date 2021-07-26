package hedis

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
)

var ErrInvalidMessage = errors.New("invalid message")
var ErrInvalidObjectType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")

var SimpleOK = NewSimpleString("ok")
var SimpleNil = NewSimpleStr(nil)

var BulkNil = NewBulkStr(nil)

var IntegerZero = NewInteger(0)
var IntegerOne = NewInteger(1)

var ErrorCommandNotFound = NewErrorString("Command not supported")
var ErrorInvalidArgNum = NewErrorString("ERR wrong number of arguments for this command")
var ErrorInvalidObjectType = NewErrorErr(ErrInvalidObjectType)

type Message interface {
	Command() *String
	Args() []*String
	ToString() *String
}

type Simple struct {
	str *String
}

func (t *Simple) Command() *String {
	return t.str
}

func (t *Simple) Args() []*String {
	return nil
}

func (t *Simple) ToString() *String {
	return t.str
}

type Integer struct {
	num int
}

func (t *Integer) Command() *String {
	return nil
}

func (t *Integer) Args() []*String {
	return nil
}

func (t *Integer) ToString() *String {
	return NewStringStr(strconv.Itoa(t.num))
}

type Error struct {
	str *String
}

func (t *Error) Command() *String {
	return nil
}

func (t *Error) Args() []*String {
	return nil
}

func (t *Error) ToString() *String {
	return t.str
}

type Bulk struct {
	str *String
}

func (t *Bulk) Command() *String {
	return t.str
}

func (t *Bulk) Args() []*String {
	return nil
}

func (t *Bulk) ToString() *String {
	return t.str
}

type Array struct {
	messages []Message
}

func (t *Array) Command() *String {
	if len(t.messages) == 0 {
		return nil
	}

	msg, ok := t.messages[0].(*Bulk)
	if !ok {
		return nil
	}

	return msg.str
}

func (t *Array) Args() []*String {
	args := make([]*String, len(t.messages)-1, len(t.messages)-1)

	for i := 1; i < len(t.messages); i++ {
		args[i-1] = t.messages[i].ToString()
	}

	return args
}

func (t *Array) ToString() *String {
	return nil
}

func (t *Array) AppendStr(str *String) {
	t.messages = append(t.messages, NewBulkStr(str))
}

func (t *Array) AppendMessage(msg Message) {
	t.messages = append(t.messages, msg)
}

type Inline struct {
	args []*String
}

func (t *Inline) Command() *String {
	if len(t.args) > 0 {
		return t.args[0]
	}

	return nil
}

func (t *Inline) Args() []*String {
	if len(t.args) > 1 {
		return t.args[1:]
	}

	return nil
}

func (t *Inline) ToString() *String {
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
	res.str = NewStringStr(text)

	return res
}

func NewSimpleStr(str *String) *Simple {
	res := &Simple{}
	res.str = str

	return res
}

func NewErrorString(text string) *Error {
	res := &Error{}
	res.str = NewStringStr(text)

	return res
}

func NewErrorErr(err error) *Error {
	res := &Error{}
	res.str = NewStringStr(err.Error())

	return res
}

func NewBulkString(text string) *Bulk {
	str := NewStringStr(text)
	return NewBulkStr(str)
}

func NewBulkStr(str *String) *Bulk {
	bulk := &Bulk{}
	bulk.str = str

	return bulk
}

func NewInteger(num int) *Integer {
	i := &Integer{}
	i.num = num

	return i
}

func NewArraySize(num int) *Array {
	arr := &Array{}
	arr.messages = make([]Message, 0, num)

	return arr
}

func NewArrayEmpty() *Array {
	return NewArraySize(0)
}
