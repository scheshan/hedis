package codec

type MessageType int8

const (
	MessageTypeUnknown MessageType = iota
	MessageTypeString
	MessageTypeError
	MessageTypeInteger
	MessageTypeBulk
	MessageTypeArray
)

type Message interface {
	String() string
}
