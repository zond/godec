package godec

import (
	"io"
	"reflect"

	"github.com/zond/godec/primitives"
)

type Encoder struct {
	io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		Writer: w,
	}
}

func (self *Encoder) EncodeReflectValue(v reflect.Value) (err error) {
	return
}

func (self *Encoder) Encode(i interface{}) (err error) {
	return primitives.Encodeinterface__(self, i)
}
