package godec

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

type decodeReader struct {
	buf []byte
	pos int
}

func (self *decodeReader) readByte() (result byte, err error) {
	if self.pos >= len(self.buf) {
		err = io.EOF
		return
	}
	result = self.buf[self.pos]
	self.pos++
	return
}

func (self *decodeReader) readBytes(n int) (result []byte, err error) {
	if self.pos+n > len(self.buf) {
		err = io.EOF
		return
	}
	result = self.buf[self.pos : self.pos+n]
	self.pos += n
	return
}

func (self *decodeReader) readUint64() (result uint64, err error) {
	result, read := binary.Uvarint(self.buf[self.pos:])
	if read <= 0 {
		err = fmt.Errorf("Unable to read uint: %v")
		return
	}
	self.pos += read
	return
}

type Decoder struct {
	reader io.Reader
}

func Unmarshal(b []byte, i interface{}) (err error) {
	r := &decodeReader{
		buf: b,
	}
	if err = decode(r, i); err != nil {
		return
	}
	return
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader: r,
	}
}

type readerByteReader struct {
	io.Reader
}

func (self readerByteReader) ReadByte() (result byte, err error) {
	b := []byte{0}
	read, err := self.Reader.Read(b)
	if err != nil {
		return
	}
	if read != 1 {
		err = fmt.Errorf("Tried to read from %v, got nothing", self.Reader)
		return
	}
	result = b[0]
	return
}

func (self *Decoder) Decode(i interface{}) (err error) {
	l, err := binary.ReadUvarint(readerByteReader{self.reader})
	if err != nil {
		return
	}
	buf := make([]byte, int(l))
	if _, err = io.ReadAtLeast(self.reader, buf, len(buf)); err != nil {
		return
	}
	if err = Unmarshal(buf, i); err != nil {
		return
	}
	return
}

func decodeKind(r *decodeReader) (result reflect.Kind, err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	result = reflect.Kind(u)
	return
}

func rawdecodestring(r *decodeReader, s *string) (err error) {
	var size uint
	if err = rawdecodeuint(r, &size); err != nil {
		return
	}
	b, err := r.readBytes(int(size))
	if err != nil {
		return
	}
	*s = string(b)
	return
}

// The special case for byte slices is here, and we treat byte slices exactly like strings.
func decodeSliceOfuint8(r *decodeReader, v *[]uint8) (err error) {
	kind, err := decodeKind(r)
	if err != nil {
		return
	}
	if kind != stringKind {
		err = fmt.Errorf("Unable to decode %v into *[]uint8", kind)
		return
	}
	var l uint
	if err = rawdecodeuint(r, &l); err != nil {
		return
	}
	*v, err = r.readBytes(int(l))
	return
}

func decodeinterface__(r *decodeReader, i *interface{}) (err error) {
	return rawdecodeinterface__(r, i)
}

func rawdecodeinterface__(r *decodeReader, i *interface{}) (err error) {
	kind, err := decodeKind(r)
	if err != nil {
		return
	}
	switch kind {
	case interface__Kind:
		err = fmt.Errorf("Unable to decode raw interface to raw interface - and this should never become an issue anyway. This should never happen.")
	case stringKind:
		proxy := ""
		if err = rawdecodestring(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case intKind:
		proxy := int64(0)
		if err = rawdecodeint64(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case uintKind:
		proxy := uint64(0)
		if err = rawdecodeuint64(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case float64Kind:
		proxy := float64(0)
		if err = rawdecodefloat64(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case complex128Kind:
		proxy := complex(float64(0), float64(0))
		if err = rawdecodecomplex128(r, &proxy); err != nil {
			return
		}
		*i = proxy
	default:
		err = fmt.Errorf("Unknown kind to decode to interface: %v", kind)
	}
	return
}

func rawdecodereflect_Value(r *decodeReader, v *reflect.Value) (err error) {
	return
}

func rawdecodebool(r *decodeReader, b *bool) (err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	*b = u != 0
	return
}

func rawdecodecomplex128(r *decodeReader, c *complex128) (err error) {
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

func rawdecodecomplex64(r *decodeReader, c *complex64) (err error) {
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

func rawdecodefloat32(r *decodeReader, f *float32) (err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	*f = math.Float32frombits(uint32(u))
	return
}

func rawdecodefloat64(r *decodeReader, f *float64) (err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	*f = math.Float64frombits(u)
	return
}

func rawdecodeint(r *decodeReader, i *int) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int(i64)
	return
}

func rawdecodeint8(r *decodeReader, i *int8) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int8(i64)
	return
}

func rawdecodeint16(r *decodeReader, i *int16) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int16(i64)
	return
}

func rawdecodeint32(r *decodeReader, i *int32) (err error) {
	var i64 int64
	if err = rawdecodeint64(r, &i64); err != nil {
		return
	}
	*i = int32(i64)
	return
}

func rawdecodeint64(r *decodeReader, x *int64) (err error) {
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

func rawdecodeuint(r *decodeReader, u *uint) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint(u64)
	return
}

func rawdecodeuint8(r *decodeReader, u *uint8) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint8(u64)
	return
}

func rawdecodeuint16(r *decodeReader, u *uint16) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint16(u64)
	return
}

func rawdecodeuint32(r *decodeReader, u *uint32) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uint32(u64)
	return
}

func rawdecodeuintptr(r *decodeReader, u *uintptr) (err error) {
	var u64 uint64
	if err = rawdecodeuint64(r, &u64); err != nil {
		return
	}
	*u = uintptr(u64)
	return
}

func rawdecodeuint64(r *decodeReader, x *uint64) (err error) {
	*x, err = r.readUint64()
	return
}
