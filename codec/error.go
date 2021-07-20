package codec

import "hedis/core"

type Error struct {
	str *core.String
}

func (t *Error) String() string {
	return ""
}

func (t *Error) Read(data []byte) (int, bool, error) {
	return 0, true, nil
}

func NewError() *Error {
	err := &Error{}
	err.str = core.NewString()

	return err
}
