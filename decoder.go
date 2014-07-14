package godec

import "github.com/zond/godec/primitives"

type Decoder struct {
	reader primitives.DecodeReader
}

func NewDecoder(r primitives.DecodeReader) *Decoder {
	return &Decoder{
		reader: r,
	}
}

func (self *Decoder) Decode(i interface{}) (err error) {
	return primitives.Decodeinterface__(self.reader, i)
}
