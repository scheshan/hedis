package codec

type Error struct {
	*Simple
}

func NewError() *Error {
	err := &Error{}
	err.Simple = NewSimple()

	return err
}
