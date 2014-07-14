package godec

import (
	"io"

	"github.com/zond/godec/primitives"
)

type Encoder struct {
	writer io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer: w,
	}
}

func (self *Encoder) Encode(i interface{}) (err error) {
	return primitives.Encodeinterface__(self.writer, i)
}
