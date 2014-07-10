package godec

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"reflect"
)

type readerByteReader struct {
	io.Reader
}

func (self readerByteReader) ReadByte() (result byte, err error) {
	buf := make([]byte, 1)
	if _, err = self.Reader.Read(buf); err != nil {
		return
	}
	result = buf[0]
	return
}

func encodeKind(w io.Writer, k reflect.Kind) (err error) {
	return encodeuint64(w, uint64(k))
}

func encodereflect_Value(w io.Writer, v reflect.Value) (err error) {
	return
}

func decodeKind(r io.Reader) (result reflect.Kind, err error) {
	u, err := decodeuint(r)
	if err != nil {
		return
	}
	result = reflect.Kind(u)
	return
}

func encodebool(w io.Writer, b bool) (err error) {
	if b {
		return encodeuint64(w, 1)
	} else {
		return encodeuint64(w, 0)
	}
}

func encodefloat32(w io.Writer, f float32) (err error) {
	return encodeuint32(w, math.Float32bits(f))
}

func encodefloat64(w io.Writer, f float64) (err error) {
	return encodeuint64(w, math.Float64bits(f))
}

func decodefloat32(r io.Reader) (f float32, err error) {
	u, err := decodeuint32(r)
	f = math.Float32frombits(u)
	return
}

func decodefloat64(r io.Reader) (f float64, err error) {
	u, err := decodeuint64(r)
	f = math.Float64frombits(u)
	return
}

func decodeint(r io.Reader) (result int, err error) {
	u, err := decodeint64(r)
	result = int(u)
	return
}

func decodeint8(r io.Reader) (result int8, err error) {
	u, err := decodeint64(r)
	result = int8(u)
	return
}

func decodeint16(r io.Reader) (result int16, err error) {
	u, err := decodeint64(r)
	result = int16(u)
	return
}

func decodeint32(r io.Reader) (result int32, err error) {
	u, err := decodeint64(r)
	result = int32(u)
	return
}

func decodeint64(r io.Reader) (result int64, err error) {
	return binary.ReadVarint(readerByteReader{r})
}

func decodeuint(r io.Reader) (result uint, err error) {
	u, err := decodeuint64(r)
	result = uint(u)
	return
}

func decodeuint8(r io.Reader) (result uint8, err error) {
	u, err := decodeuint64(r)
	result = uint8(u)
	return
}

func decodeuint16(r io.Reader) (result uint16, err error) {
	u, err := decodeuint64(r)
	result = uint16(u)
	return
}

func decodeuint32(r io.Reader) (result uint32, err error) {
	u, err := decodeuint64(r)
	result = uint32(u)
	return
}

func decodeuint64(r io.Reader) (result uint64, err error) {
	return binary.ReadUvarint(readerByteReader{r})
}

func decodeInt(r io.Reader) (result int64, err error) {
	return binary.ReadVarint(readerByteReader{r})
}

func encodeint(w io.Writer, u int) (err error) {
	return encodeint64(w, int64(u))
}

func encodeint8(w io.Writer, u int8) (err error) {
	return encodeint64(w, int64(u))
}

func encodeint16(w io.Writer, u int16) (err error) {
	return encodeint64(w, int64(u))
}

func encodeint32(w io.Writer, u int32) (err error) {
	return encodeint64(w, int64(u))
}

func encodeint64(w io.Writer, u int64) (err error) {
	buf := make([]byte, binary.MaxVarintLen64)
	bw := binary.PutVarint(buf, u)
	_, err = w.Write(buf[:bw])
	return
}

func encodeuintptr(w io.Writer, u uintptr) (err error) {
	return encodeuint64(w, uint64(u))
}

func encodeuint(w io.Writer, u uint) (err error) {
	return encodeuint64(w, uint64(u))
}

func encodeuint8(w io.Writer, u uint8) (err error) {
	return encodeuint64(w, uint64(u))
}

func encodeuint16(w io.Writer, u uint16) (err error) {
	return encodeuint64(w, uint64(u))
}

func encodeuint32(w io.Writer, u uint32) (err error) {
	return encodeuint64(w, uint64(u))
}

func encodeuint64(w io.Writer, u uint64) (err error) {
	buf := make([]byte, binary.MaxVarintLen64)
	bw := binary.PutUvarint(buf, u)
	_, err = w.Write(buf[:bw])
	return
}

func encodecomplex64(w io.Writer, c complex64) (err error) {
	if err = encodefloat32(w, real(c)); err != nil {
		return
	}
	return encodefloat32(w, imag(c))
}

func encodecomplex128(w io.Writer, c complex128) (err error) {
	if err = encodefloat64(w, real(c)); err != nil {
		return
	}
	return encodefloat64(w, imag(c))
}

var compatibleKinds = [][]reflect.Kind{
	{reflect.Invalid},
	{reflect.Bool},
	{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64},
	{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64},
	{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64},
	{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64},
	{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64},
	{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr},
	{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr},
	{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr},
	{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr},
	{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr},
	{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr},
	{reflect.Float32, reflect.Float64},
	{reflect.Float32, reflect.Float64},
	{reflect.Complex64, reflect.Complex128},
	{reflect.Complex64, reflect.Complex128},
	{reflect.Array},
	{reflect.Chan},
	{reflect.Func},
	{reflect.Interface},
	{reflect.Map},
	{reflect.Ptr},
	{reflect.Slice},
	{reflect.String},
	{reflect.Struct},
	{reflect.UnsafePointer},
}

func kindsCompatible(k1, k2 reflect.Kind) bool {
	for _, comp := range compatibleKinds[int(k1)] {
		if k2 == comp {
			return true
		}
	}
	return false
}

func encodestring(w io.Writer, s string) (err error) {
	if err = encodeint(w, len(s)); err != nil {
		return
	}
	_, err = io.WriteString(w, s)
	return
}

func decodestring(r io.Reader) (result string, err error) {
	var size int
	if size, err = decodeint(r); err != nil {
		return
	}
	b := make([]byte, int(size))
	if _, err = io.ReadAtLeast(r, b, len(b)); err != nil {
		return
	}
	result = string(b)
	return
}

func Marshal(i interface{}) (result []byte, err error) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)
	if err = enc.Encode(i); err != nil {
		return
	}
	result = buf.Bytes()
	return
}

func Unmarshal(b []byte, i interface{}) (err error) {
	buf := bytes.NewBuffer(b)
	dec := NewDecoder(buf)
	if err = dec.Decode(i); err != nil {
		return
	}
	return
}
