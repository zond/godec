package godec

import (
	"io"
	"reflect"

	"github.com/zond/godec/primitives"
)

type ioDecodeReader struct {
	Reader io.Reader
}

func (self ioDecodeReader) ReadByte() (result byte, err error) {
	buf := []byte{0}
	if _, err = self.Reader.Read(buf); err != nil {
		return
	}
	result = buf[0]
	return
}

func (self ioDecodeReader) ReadBytes(n int) (result []byte, err error) {
	result = make([]byte, n)
	_, err = io.ReadAtLeast(self.Reader, result, n)
	return
}

func (self ioDecodeReader) ReadUint64() (result uint64, err error) {
	b, err := self.ReadBytes(8)
	if err != nil {
		return
	}
	result = uint64(b[0])<<56 + uint64(b[1])<<48 + uint64(b[2])<<40 + uint64(b[3])<<32 + uint64(b[4])<<24 + uint64(b[5])<<16 + uint64(b[6])<<8 + uint64(b[7])
	return
}

type Decoder struct {
	primitives.DecodeReader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		DecodeReader: ioDecodeReader{r},
	}
}

func (self *Decoder) DecodeReflectValue(v reflect.Value) (err error) {
	return
}

func (self *Decoder) Decode(i interface{}) (err error) {
	return primitives.Decodeinterface__(self, i)
}
