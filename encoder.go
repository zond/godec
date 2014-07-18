package godec

import (
	"io"
	"reflect"

	"github.com/zond/godec/primitives"
)

type ioEncodeWriter struct {
	io.Writer
}

func (self ioEncodeWriter) WriteBytes(b []byte) (err error) {
	_, err = self.Writer.Write(b)
	return
}

func (self ioEncodeWriter) WriteByte(b byte) (err error) {
	_, err = self.Writer.Write([]byte{b})
	return
}

func (self ioEncodeWriter) WriteString(s string) (err error) {
	_, err = io.WriteString(self.Writer, s)
	return
}

func (self ioEncodeWriter) WriteUint64(u uint64) (err error) {
	return self.WriteBytes([]byte{
		byte(u >> 56),
		byte(u >> 48),
		byte(u >> 40),
		byte(u >> 32),
		byte(u >> 24),
		byte(u >> 16),
		byte(u >> 8),
		byte(u),
	})
}

type Encoder struct {
	primitives.EncodeWriter
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		EncodeWriter: ioEncodeWriter{w},
	}
}

func (self *Encoder) EncodeReflectValue(v reflect.Value) (err error) {
	return
}

func (self *Encoder) Encode(i interface{}) (err error) {
	return primitives.Encodeinterface__(self, i)
}
