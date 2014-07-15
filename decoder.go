package godec

import (
	"io"
	"reflect"

	"github.com/zond/godec/primitives"
)

type IODecodeReader struct {
	Reader io.Reader
}

func (self IODecodeReader) ReadByte() (result byte, err error) {
	buf := []byte{0}
	if _, err = self.Reader.Read(buf); err != nil {
		return
	}
	result = buf[0]
	return
}

func (self IODecodeReader) ReadBytes(n int) (result []byte, err error) {
	result = make([]byte, n)
	_, err = io.ReadAtLeast(self.Reader, result, n)
	return
}

type BytesDecodeReader struct {
	Buf []byte
	Pos int
}

func (self *BytesDecodeReader) ReadByte() (result byte, err error) {
	if self.Pos >= len(self.Buf) {
		err = io.EOF
		return
	}
	result = self.Buf[self.Pos]
	self.Pos++
	return
}

func (self *BytesDecodeReader) ReadBytes(n int) (result []byte, err error) {
	if self.Pos+n > len(self.Buf) {
		err = io.EOF
		return
	}
	result = self.Buf[self.Pos : self.Pos+n]
	self.Pos += n
	return
}

type Decoder struct {
	primitives.DecodeReader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		DecodeReader: IODecodeReader{r},
	}
}

func (self *Decoder) DecodeReflectValue(v reflect.Value) (err error) {
	return
}

func (self *Decoder) Decode(i interface{}) (err error) {
	return primitives.Decodeinterface__(self, i)
}
