package godec

import (
	"io"
	"math"
	"reflect"
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
	return encodeinterface__(self.writer, i)
}

func encodeKind(w io.Writer, k reflect.Kind) (err error) {
	return rawencodeuint64(w, uint64(k))
}

func encodereflect_Value(w io.Writer, v reflect.Value) (err error) {
	return
}

func rawencodeinterface__(w io.Writer, v interface{}) (err error) {
	return encodeinterface__(w, v)
}

func rawencodereflect_Value(w io.Writer, v reflect.Value) (err error) {
	return encodereflect_Value(w, v)
}

// The special case for byte slices is here.
func encodeSliceOfuint8(w io.Writer, v []uint8) (err error) {
	if err = encodeKind(w, reflect.Slice); err != nil {
		return
	}
	if err = encodeKind(w, uint8Kind); err != nil {
		return
	}
	if err = rawencodeuint(w, uint(len(v))); err != nil {
		return
	}
	_, err = w.Write(v)
	return
}

func rawencodebool(w io.Writer, b bool) (err error) {
	if b {
		return rawencodeuint64(w, 1)
	} else {
		return rawencodeuint64(w, 0)
	}
}

func rawencodefloat32(w io.Writer, f float32) (err error) {
	return rawencodeuint64(w, uint64(math.Float32bits(f)))
}

func rawencodefloat64(w io.Writer, f float64) (err error) {
	return rawencodeuint64(w, math.Float64bits(f))
}

func rawencodeint(w io.Writer, u int) (err error) {
	return rawencodeint64(w, int64(u))
}

func rawencodeint8(w io.Writer, u int8) (err error) {
	return rawencodeint64(w, int64(u))
}

func rawencodeint16(w io.Writer, u int16) (err error) {
	return rawencodeint64(w, int64(u))
}

func rawencodeint32(w io.Writer, u int32) (err error) {
	return rawencodeint64(w, int64(u))
}

func rawencodeuintptr(w io.Writer, u uintptr) (err error) {
	return rawencodeuint64(w, uint64(u))
}

func rawencodeuint(w io.Writer, u uint) (err error) {
	return rawencodeuint64(w, uint64(u))
}

func rawencodeuint8(w io.Writer, u uint8) (err error) {
	return rawencodeuint64(w, uint64(u))
}

func rawencodeuint16(w io.Writer, u uint16) (err error) {
	return rawencodeuint64(w, uint64(u))
}

func rawencodeuint32(w io.Writer, u uint32) (err error) {
	return rawencodeuint64(w, uint64(u))
}

func rawencodeint64(w io.Writer, i int64) (err error) {
	ux := uint64(i) << 1
	if i < 0 {
		ux = ^ux
	}
	return rawencodeuint64(w, ux)
}

func rawencodeuint64(w io.Writer, u uint64) (err error) {
	for u >= 0x80 {
		if _, err = w.Write([]byte{byte(u) | 0x80}); err != nil {
			return
		}
		u >>= 7
	}
	if _, err = w.Write([]byte{byte(u)}); err != nil {
		return
	}
	return
}

func rawencodecomplex64(w io.Writer, c complex64) (err error) {
	if err = rawencodefloat64(w, float64(real(c))); err != nil {
		return
	}
	return rawencodefloat64(w, float64(imag(c)))
}

func rawencodecomplex128(w io.Writer, c complex128) (err error) {
	if err = rawencodefloat64(w, real(c)); err != nil {
		return
	}
	return rawencodefloat64(w, imag(c))
}

func rawencodestring(w io.Writer, s string) (err error) {
	if err = rawencodeuint(w, uint(len(s))); err != nil {
		return
	}
	_, err = io.WriteString(w, s)
	return
}
