package codec

type Array struct {
	messages []Message
}

func (t *Array) String() string {
	return ""
}

func (t *Array) Read(data []byte) (int, bool, error) {
	return 0, true, nil
}

func NewArray() *Array {
	arr := &Array{}

	return arr
}
