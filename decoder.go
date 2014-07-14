package godec

import (
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

type DecodeReader interface {
	io.ByteReader
	ReadBytes(n int) ([]byte, error)
}

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
	reader DecodeReader
}

func NewDecoder(r DecodeReader) *Decoder {
	return &Decoder{
		reader: r,
	}
}

func (self *Decoder) Decode(i interface{}) (err error) {
	return decodeinterface__(self.reader, i)
}

func decodeKind(r DecodeReader) (result reflect.Kind, err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	result = reflect.Kind(u)
	return
}

func decodereflect_Value(r DecodeReader, v reflect.Value) (err error) {
	return
}

// The special case for byte slices is here.
func decodeSliceOfuint8(r DecodeReader, v *[]uint8) (err error) {
	kind, err := decodeKind(r)
	if err != nil {
		return
	}
	if kind != reflect.Slice {
		err = fmt.Errorf("Unable to decode %v into *[]uint8", kind)
		return
	}
	elkind, err := decodeKind(r)
	if err != nil {
		return
	}
	if elkind != uint8Kind {
		err = fmt.Errorf("Unable to decode %v into uint8", elkind)
		return
	}
	var l uint
	if err = rawdecodeuint(r, &l); err != nil {
		return
	}
	*v, err = r.ReadBytes(int(l))
	return
}

func rawdecodeinterface__(r DecodeReader, i *interface{}) (err error) {
	return
}

func rawdecodereflect_Value(r DecodeReader, v *reflect.Value) (err error) {
	return
}

func rawdecodebool(r DecodeReader, b *bool) (err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	*b = u != 0
	return
}

func rawdecodecomplex128(r DecodeReader, c *complex128) (err error) {
	var re float64
	if err = rawdecodefloat64(r, &re); err != nil {
		return
	}
	var im float64
	if err = rawdecodefloat64(r, &im); err != nil {
		return
	}
	*c = complex(re, im)
	return
}

func rawdecodecomplex64(r DecodeReader, c *complex64) (err error) {
	var re float64
	if err = rawdecodefloat64(r, &re); err != nil {
		return
	}
	var im float64
	if err = rawdecodefloat64(r, &im); err != nil {
		return
	}
	*c = complex(float32(re), float32(im))
	return
}

func rawdecodefloat32(r DecodeReader, f *float32) (err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	*f = math.Float32frombits(uint32(u))
	return
}

func rawdecodefloat64(r DecodeReader, f *float64) (err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	*f = math.Float64frombits(u)
	return
}

func rawdecodeint(r DecodeReader, i *int) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int(i64)
	return
}

func rawdecodeint8(r DecodeReader, i *int8) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int8(i64)
	return
}

func rawdecodeint16(r DecodeReader, i *int16) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int16(i64)
	return
}

func rawdecodeint32(r DecodeReader, i *int32) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int32(i64)
	return
}

func rawdecodeint64(r DecodeReader, x *int64) (err error) {
	var ux uint64
	if err = rawdecodeuint64(r, &ux); err != nil {
		return
	}
	*x = int64(ux >> 1)
	if ux&1 != 0 {
		*x = ^*x
	}
	return
}

func rawdecodeuint(r DecodeReader, u *uint) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint(u64)
	return
}

func rawdecodeuint8(r DecodeReader, u *uint8) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint8(u64)
	return
}

func rawdecodeuint16(r DecodeReader, u *uint16) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint16(u64)
	return
}

func rawdecodeuint32(r DecodeReader, u *uint32) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint32(u64)
	return
}

func rawdecodeuintptr(r DecodeReader, u *uintptr) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uintptr(u64)
	return
}

var overflow = errors.New("binary: varint overflows a 64-bit integer")

func rawdecodeuint64(r DecodeReader, x *uint64) (err error) {
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

func rawdecodestring(r DecodeReader, s *string) (err error) {
	var size uint
	if err = rawdecodeuint(r, &size); err != nil {
		return
	}
	b, err := r.ReadBytes(int(size))
	if err != nil {
		return
	}
	*s = string(b)
	return
}
