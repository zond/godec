package primitives

import (
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

const (
	interface__Kind = reflect.Interface
	stringKind      = reflect.String
	boolKind        = reflect.Uint64
	intKind         = reflect.Int64
	int8Kind        = reflect.Int64
	int16Kind       = reflect.Int64
	int32Kind       = reflect.Int64
	int64Kind       = reflect.Int64
	uintKind        = reflect.Uint64
	uint8Kind       = reflect.Uint64
	uint16Kind      = reflect.Uint64
	uint32Kind      = reflect.Uint64
	uint64Kind      = reflect.Uint64
	uintptrKind     = reflect.Uint64
	float32Kind     = reflect.Float64
	float64Kind     = reflect.Float64
	complex64Kind   = reflect.Complex128
	complex128Kind  = reflect.Complex128
)

type Encoder interface {
	io.Writer
	EncodeReflectValue(reflect.Value) error
}

type DecodeReader interface {
	ReadByte() (byte, error)
	ReadBytes(n int) ([]byte, error)
}

type Decoder interface {
	DecodeReader
	DecodeReflectValue(reflect.Value) error
}

func EncodeKind(w Encoder, k reflect.Kind) (err error) {
	return Rawencodeuint64(w, uint64(k))
}

func Rawencodeinterface__(w Encoder, v interface{}) (err error) {
	return Encodeinterface__(w, v)
}

func Rawencodebool(w Encoder, b bool) (err error) {
	if b {
		return Rawencodeuint64(w, 1)
	} else {
		return Rawencodeuint64(w, 0)
	}
}

func Rawencodefloat32(w Encoder, f float32) (err error) {
	return Rawencodeuint64(w, uint64(math.Float32bits(f)))
}

func Rawencodefloat64(w Encoder, f float64) (err error) {
	return Rawencodeuint64(w, math.Float64bits(f))
}

func Rawencodeint(w Encoder, u int) (err error) {
	return Rawencodeint64(w, int64(u))
}

func Rawencodeint8(w Encoder, u int8) (err error) {
	return Rawencodeint64(w, int64(u))
}

func Rawencodeint16(w Encoder, u int16) (err error) {
	return Rawencodeint64(w, int64(u))
}

func Rawencodeint32(w Encoder, u int32) (err error) {
	return Rawencodeint64(w, int64(u))
}

func Rawencodeuintptr(w Encoder, u uintptr) (err error) {
	return Rawencodeuint64(w, uint64(u))
}

func Rawencodeuint(w Encoder, u uint) (err error) {
	return Rawencodeuint64(w, uint64(u))
}

func Rawencodeuint8(w Encoder, u uint8) (err error) {
	return Rawencodeuint64(w, uint64(u))
}

func Rawencodeuint16(w Encoder, u uint16) (err error) {
	return Rawencodeuint64(w, uint64(u))
}

func Rawencodeuint32(w Encoder, u uint32) (err error) {
	return Rawencodeuint64(w, uint64(u))
}

func Rawencodeint64(w Encoder, i int64) (err error) {
	ux := uint64(i) << 1
	if i < 0 {
		ux = ^ux
	}
	return Rawencodeuint64(w, ux)
}

func Rawencodeuint64(w Encoder, u uint64) (err error) {
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

func Rawencodecomplex64(w Encoder, c complex64) (err error) {
	if err = Rawencodefloat64(w, float64(real(c))); err != nil {
		return
	}
	return Rawencodefloat64(w, float64(imag(c)))
}

func Rawencodecomplex128(w Encoder, c complex128) (err error) {
	if err = Rawencodefloat64(w, real(c)); err != nil {
		return
	}
	return Rawencodefloat64(w, imag(c))
}

func DecodeKind(r Decoder) (result reflect.Kind, err error) {
	var u uint64
	if err = Rawdecodeuint64(r, &u); err != nil {
		return
	}
	result = reflect.Kind(u)
	return
}

func Rawencodestring(w Encoder, s string) (err error) {
	if err = Rawencodeuint64(w, uint64(len(s))); err != nil {
		return
	}
	_, err = io.WriteString(w, s)
	return
}

// The special case for byte slices is here, and we treat byte slices exactly like strings.
func EncodeSliceOfuint8(w Encoder, v []uint8) (err error) {
	if err = EncodeKind(w, stringKind); err != nil {
		return
	}
	if err = Rawencodeuint(w, uint(len(v))); err != nil {
		return
	}
	_, err = w.Write(v)
	return
}

func Rawdecodestring(r Decoder, s *string) (err error) {
	var size uint
	if err = Rawdecodeuint(r, &size); err != nil {
		return
	}
	b, err := r.ReadBytes(int(size))
	if err != nil {
		return
	}
	*s = string(b)
	return
}

// The special case for byte slices is here, and we treat byte slices exactly like strings.
func DecodeSliceOfuint8(r Decoder, v *[]uint8) (err error) {
	kind, err := DecodeKind(r)
	if err != nil {
		return
	}
	if kind != stringKind {
		err = fmt.Errorf("Unable to decode %v into *[]uint8", kind)
		return
	}
	var l uint
	if err = Rawdecodeuint(r, &l); err != nil {
		return
	}
	*v, err = r.ReadBytes(int(l))
	return
}

func Rawdecodeinterface__(r Decoder, i *interface{}) (err error) {
	return
}

func Rawdecodereflect_Value(r Decoder, v *reflect.Value) (err error) {
	return
}

func Rawdecodebool(r Decoder, b *bool) (err error) {
	var u uint64
	if err = Rawdecodeuint64(r, &u); err != nil {
		return
	}
	*b = u != 0
	return
}

func Rawdecodecomplex128(r Decoder, c *complex128) (err error) {
	var re float64
	if err = Rawdecodefloat64(r, &re); err != nil {
		return
	}
	var im float64
	if err = Rawdecodefloat64(r, &im); err != nil {
		return
	}
	*c = complex(re, im)
	return
}

func Rawdecodecomplex64(r Decoder, c *complex64) (err error) {
	var re float64
	if err = Rawdecodefloat64(r, &re); err != nil {
		return
	}
	var im float64
	if err = Rawdecodefloat64(r, &im); err != nil {
		return
	}
	*c = complex(float32(re), float32(im))
	return
}

func Rawdecodefloat32(r Decoder, f *float32) (err error) {
	var u uint64
	if err = Rawdecodeuint64(r, &u); err != nil {
		return
	}
	*f = math.Float32frombits(uint32(u))
	return
}

func Rawdecodefloat64(r Decoder, f *float64) (err error) {
	var u uint64
	if err = Rawdecodeuint64(r, &u); err != nil {
		return
	}
	*f = math.Float64frombits(u)
	return
}

func Rawdecodeint(r Decoder, i *int) (err error) {
	var i64 int64
	if err = Rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int(i64)
	return
}

func Rawdecodeint8(r Decoder, i *int8) (err error) {
	var i64 int64
	if err = Rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int8(i64)
	return
}

func Rawdecodeint16(r Decoder, i *int16) (err error) {
	var i64 int64
	if err = Rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int16(i64)
	return
}

func Rawdecodeint32(r Decoder, i *int32) (err error) {
	var i64 int64
	if err = Rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int32(i64)
	return
}

func Rawdecodeint64(r Decoder, x *int64) (err error) {
	var ux uint64
	if err = Rawdecodeuint64(r, &ux); err != nil {
		return
	}
	*x = int64(ux >> 1)
	if ux&1 != 0 {
		*x = ^*x
	}
	return
}

func Rawdecodeuint(r Decoder, u *uint) (err error) {
	var u64 uint64
	if err = Rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint(u64)
	return
}

func Rawdecodeuint8(r Decoder, u *uint8) (err error) {
	var u64 uint64
	if err = Rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint8(u64)
	return
}

func Rawdecodeuint16(r Decoder, u *uint16) (err error) {
	var u64 uint64
	if err = Rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint16(u64)
	return
}

func Rawdecodeuint32(r Decoder, u *uint32) (err error) {
	var u64 uint64
	if err = Rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint32(u64)
	return
}

func Rawdecodeuintptr(r Decoder, u *uintptr) (err error) {
	var u64 uint64
	if err = Rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uintptr(u64)
	return
}

var overflow = errors.New("binary: varint overflows a 64-bit integer")

func Rawdecodeuint64(r Decoder, x *uint64) (err error) {
	*x = 0
	var s uint
	var b byte
	for i := 0; ; i++ {
		if b, err = r.ReadByte(); err != nil {
			return
		}
		if err != nil {
			return
		}
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				err = overflow
				return
			}
			*x = *x | uint64(b)<<s
			return
		}
		*x = *x | uint64(b&0x7f)<<s
		s += 7
	}
	return
}
