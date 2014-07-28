package godec

import (
	"encoding"
	"encoding/binary"
	"io"
	"math"
	"time"
)

type encodeWriter struct {
	buf []byte
	pos int
}

func newEncodeWriter() *encodeWriter {
	return &encodeWriter{
		buf: make([]byte, 1<<6),
	}
}

func (self *encodeWriter) writeBytes(b []byte) (err error) {
	l := len(b)
	self.grow(l)
	copy(self.buf[self.pos:], b)
	self.pos += l
	return
}

func (self *encodeWriter) writeUint64(u uint64) (err error) {
	self.grow(binary.MaxVarintLen64)
	wrote := binary.PutUvarint(self.buf[self.pos:], u)
	self.pos += wrote
	return
}

func (self *encodeWriter) grow(n int) {
	if n+self.pos > len(self.buf) {
		newBuf := make([]byte, len(self.buf)*2+n)
		copy(newBuf, self.buf)
		self.buf = newBuf
	}
}

func (self *encodeWriter) writeString(s string) (err error) {
	l := len(s)
	self.grow(l)
	copy(self.buf[self.pos:], s)
	self.pos += l
	return
}

func Marshal(i interface{}) (result []byte, err error) {
	w := &encodeWriter{
		buf: make([]byte, 1<<6),
	}
	if err = encodeinterface__(w, i); err != nil {
		return
	}
	result = w.buf[:w.pos]
	return
}

type Encoder struct {
	writer io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer: w,
	}
}

func (self *Encoder) Encode(i interface{}) (err error) {
	b, err := Marshal(i)
	if err != nil {
		return
	}
	lb := make([]byte, binary.MaxVarintLen64)
	wrote := binary.PutUvarint(lb, uint64(len(b)))
	if _, err = self.writer.Write(lb[:wrote]); err != nil {
		return
	}
	if _, err = self.writer.Write(b); err != nil {
		return
	}
	return
}

func encodeType(w *encodeWriter, t *Type) (err error) {
	if err = rawencodeuint64(w, uint64(t.Base)); err != nil {
		return
	}
	if t.Key != nil {
		if err = encodeType(w, t.Key); err != nil {
			return
		}
	}
	if t.Value != nil {
		if err = encodeType(w, t.Value); err != nil {
			return
		}
	}
	return
}

func rawencodetime_Time(w *encodeWriter, t time.Time) (err error) {
	return rawencodeint64(w, t.UnixNano())
}

func encodebinary_Marshaler(w *encodeWriter, bm encoding.BinaryMarshaler) (err error) {
	if err = encodeType(w, &Type{Base: binaryUnMarshalerKind}); err != nil {
		return
	}
	b, err := bm.MarshalBinary()
	if err != nil {
		return
	}
	if err = rawencodeuint(w, uint(len(b))); err != nil {
		return
	}
	if err = w.writeBytes(b); err != nil {
		return
	}
	return
}

func rawencodeinterface__(w *encodeWriter, v interface{}) (err error) {
	return encodeinterface__(w, v)
}

func rawencodestring(w *encodeWriter, s string) (err error) {
	if err = rawencodeuint64(w, uint64(len(s))); err != nil {
		return
	}
	err = w.writeString(s)
	return
}

// The special case for byte slices is here, and we treat byte slices exactly like strings.
func encodeSliceOfuint8(w *encodeWriter, v []uint8) (err error) {
	if err = encodeType(w, &Type{Base: stringKind}); err != nil {
		return
	}
	if err = rawencodeuint(w, uint(len(v))); err != nil {
		return
	}
	err = w.writeBytes(v)
	return
}

func rawencodebool(w *encodeWriter, b bool) (err error) {
	if b {
		return rawencodeuint64(w, 1)
	} else {
		return rawencodeuint64(w, 0)
	}
}

func rawencodefloat32(w *encodeWriter, f float32) (err error) {
	return rawencodeuint64(w, uint64(math.Float32bits(f)))
}

func rawencodefloat64(w *encodeWriter, f float64) (err error) {
	return rawencodeuint64(w, math.Float64bits(f))
}

func rawencodeint(w *encodeWriter, u int) (err error) {
	return rawencodeint64(w, int64(u))
}

func rawencodeint8(w *encodeWriter, u int8) (err error) {
	return rawencodeint64(w, int64(u))
}

func rawencodeint16(w *encodeWriter, u int16) (err error) {
	return rawencodeint64(w, int64(u))
}

func rawencodeint32(w *encodeWriter, u int32) (err error) {
	return rawencodeint64(w, int64(u))
}

func rawencodeuintptr(w *encodeWriter, u uintptr) (err error) {
	return rawencodeuint64(w, uint64(u))
}

func rawencodeuint(w *encodeWriter, u uint) (err error) {
	return rawencodeuint64(w, uint64(u))
}

func rawencodeuint8(w *encodeWriter, u uint8) (err error) {
	return rawencodeuint64(w, uint64(u))
}

func rawencodeuint16(w *encodeWriter, u uint16) (err error) {
	return rawencodeuint64(w, uint64(u))
}

func rawencodeuint32(w *encodeWriter, u uint32) (err error) {
	return rawencodeuint64(w, uint64(u))
}

func rawencodeint64(w *encodeWriter, i int64) (err error) {
	ux := uint64(i) << 1
	if i < 0 {
		ux = ^ux
	}
	return rawencodeuint64(w, ux)
}

func rawencodeuint64(w *encodeWriter, u uint64) (err error) {
	return w.writeUint64(u)
}

func rawencodecomplex64(w *encodeWriter, c complex64) (err error) {
	if err = rawencodefloat64(w, float64(real(c))); err != nil {
		return
	}
	return rawencodefloat64(w, float64(imag(c)))
}

func rawencodecomplex128(w *encodeWriter, c complex128) (err error) {
	if err = rawencodefloat64(w, real(c)); err != nil {
		return
	}
	return rawencodefloat64(w, imag(c))
}
