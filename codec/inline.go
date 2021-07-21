package codec

import (
	"bufio"
	"hedis/core"
)

type Inline struct {
	args []*core.String
}

func NewInline() *Inline {
	inline := &Inline{}

	return inline
}

func (t *Inline) String() string {

}

func (t *Inline) Read(reader *bufio.Reader) error {

}

func (t *Inline) Write(writer *bufio.Writer) error {

}
