package godec

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

type Decoder struct {
	reader io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader: r,
	}
}

func (self *Decoder) Decode(i interface{}) (err error) {
	return decodeinterface__(self.reader, i)
}

func decodeKind(r io.Reader) (result reflect.Kind, err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	result = reflect.Kind(u)
	return
}

func decodereflect_Value(r io.Reader, v *reflect.Value) (err error) {
	return
}

// The special case for byte slices is here.
func decodeSliceOfuint8(r io.Reader, v *[]uint8) (err error) {
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
	var l int
	if err = rawdecodeint(r, &l); err != nil {
		return
	}
	*v = make([]uint8, l)
	_, err = io.ReadAtLeast(r, *v, l)
	return
}

func rawdecodeinterface__(r io.Reader, i *interface{}) (err error) {
	return
}

func rawdecodereflect_Value(r io.Reader, v *reflect.Value) (err error) {
	return
}

func rawdecodebool(r io.Reader, b *bool) (err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	*b = u != 0
	return
}

func rawdecodecomplex128(r io.Reader, c *complex128) (err error) {
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

func rawdecodecomplex64(r io.Reader, c *complex64) (err error) {
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

func rawdecodefloat32(r io.Reader, f *float32) (err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	*f = math.Float32frombits(uint32(u))
	return
}

func rawdecodefloat64(r io.Reader, f *float64) (err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	*f = math.Float64frombits(u)
	return
}

func rawdecodeint(r io.Reader, i *int) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int(i64)
	return
}

func rawdecodeint8(r io.Reader, i *int8) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int8(i64)
	return
}

func rawdecodeint16(r io.Reader, i *int16) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int16(i64)
	return
}

func rawdecodeint32(r io.Reader, i *int32) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int32(i64)
	return
}

func rawdecodeint64(r io.Reader, i *int64) (err error) {
	*i, err = binary.ReadVarint(readerByteReader{r})
	return
}

func rawdecodeuint(r io.Reader, u *uint) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint(u64)
	return
}

func rawdecodeuint8(r io.Reader, u *uint8) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint8(u64)
	return
}

func rawdecodeuint16(r io.Reader, u *uint16) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint16(u64)
	return
}

func rawdecodeuint32(r io.Reader, u *uint32) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint32(u64)
	return
}

func rawdecodeuintptr(r io.Reader, u *uintptr) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uintptr(u64)
	return
}

func rawdecodeuint64(r io.Reader, u *uint64) (err error) {
	*u, err = binary.ReadUvarint(readerByteReader{r})
	return
}

func rawdecodestring(r io.Reader, s *string) (err error) {
	var size int
	if err = rawdecodeint(r, &size); err != nil {
		return
	}
	b := make([]byte, int(size))
	if _, err = io.ReadAtLeast(r, b, len(b)); err != nil {
		return
	}
	*s = string(b)
	return
}
