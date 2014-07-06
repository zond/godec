package godec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"sync"
	"unsafe"
)

var codecByType = map[reflect.Type]*codec{}
var codecLock = &sync.RWMutex{}

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

type codec struct {
	kind            reflect.Kind
	generatedEncode func(w io.Writer, v reflect.Value) (err error)
	generatedDecode func(r io.Reader, v reflect.Value) (err error)
}

func (self *codec) encode(w io.Writer, v reflect.Value) (err error) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if err = encodeKind(w, self.kind); err != nil {
		return
	}
	return self.generatedEncode(w, v)
}

func encodeKind(w io.Writer, k reflect.Kind) (err error) {
	return encodeUint(w, uint64(k))
}

func decodeKind(r io.Reader) (result reflect.Kind, err error) {
	u, err := decodeUint(r)
	if err != nil {
		return
	}
	result = reflect.Kind(u)
	return
}

func encodeFloat(w io.Writer, f float64) (err error) {
	return encodeUint(w, *(*uint64)(unsafe.Pointer(&f)))
}

func decodeFloat(r io.Reader) (f float64, err error) {
	u, err := decodeUint(r)
	if err != nil {
		return
	}
	f = *(*float64)(unsafe.Pointer(&u))
	return
}

func encodeFloatValue(w io.Writer, v reflect.Value) (err error) {
	return encodeFloat(w, v.Float())
}

func decodeFloatValue(r io.Reader, v reflect.Value) (err error) {
	f, err := decodeFloat(r)
	if err != nil {
		return
	}
	v.SetFloat(f)
	return
}

func decodeUint(r io.Reader) (result uint64, err error) {
	return binary.ReadUvarint(readerByteReader{r})
}

func decodeInt(r io.Reader) (result int64, err error) {
	return binary.ReadVarint(readerByteReader{r})
}

func encodeUint(w io.Writer, u uint64) (err error) {
	buf := make([]byte, binary.MaxVarintLen64)
	bw := binary.PutUvarint(buf, u)
	_, err = w.Write(buf[:bw])
	return
}

func encodeInt(w io.Writer, i int64) (err error) {
	buf := make([]byte, 8)
	bw := binary.PutVarint(buf, i)
	_, err = w.Write(buf[:bw])
	return
}

func encodeBoolValue(w io.Writer, v reflect.Value) (err error) {
	if v.Bool() {
		return encodeUint(w, 1)
	} else {
		return encodeUint(w, 0)
	}
}

func decodeBoolValue(r io.Reader, v reflect.Value) (err error) {
	i, err := decodeUint(r)
	if err != nil {
		return
	}
	if i == 0 {
		v.SetBool(false)
	} else {
		v.SetBool(true)
	}
	return
}

func decodeIntValue(r io.Reader, v reflect.Value) (err error) {
	i, err := decodeInt(r)
	if err != nil {
		return
	}
	v.SetInt(i)
	return
}

func decodeUintValue(r io.Reader, v reflect.Value) (err error) {
	u, err := decodeUint(r)
	if err != nil {
		return
	}
	v.SetUint(u)
	return
}

func encodeIntValue(w io.Writer, v reflect.Value) (err error) {
	return encodeInt(w, v.Int())
}

func encodeUintValue(w io.Writer, v reflect.Value) (err error) {
	return encodeUint(w, v.Uint())
}

func encodeComplexValue(w io.Writer, v reflect.Value) (err error) {
	cpy := v.Complex()
	if err = encodeFloat(w, real(cpy)); err != nil {
		return
	}
	return encodeFloat(w, imag(cpy))
}

func decodeComplexValue(r io.Reader, v reflect.Value) (err error) {
	re, err := decodeFloat(r)
	if err != nil {
		return
	}
	im, err := decodeFloat(r)
	if err != nil {
		return
	}
	v.SetComplex(complex(re, im))
	return
}

func (self *codec) decode(r io.Reader, v reflect.Value) (err error) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	kind, err := decodeKind(r)
	if err != nil {
		return
	}
	if self.kind != kind {
		err = fmt.Errorf("Wrong kind in stream, wanted %v but got %v", self.kind, kind)
		return
	}
	return self.generatedDecode(r, v)
}

func createCodec(t reflect.Type) (result *codec, err error) {
	result = &codec{
		kind: t.Kind(),
	}
	switch result.kind {
	case reflect.Invalid:
		err = fmt.Errorf("Unable to create codec for invalid kind")
		return
	case reflect.Bool:
		result.generatedEncode = encodeBoolValue
		result.generatedDecode = decodeBoolValue
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		result.generatedEncode = encodeIntValue
		result.generatedDecode = decodeIntValue
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uintptr:
		result.generatedEncode = encodeUintValue
		result.generatedDecode = decodeUintValue
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		result.generatedEncode = encodeFloatValue
		result.generatedDecode = decodeFloatValue
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		result.generatedEncode = encodeComplexValue
		result.generatedDecode = decodeComplexValue
	case reflect.Array:
	case reflect.Chan:
		err = fmt.Errorf("Unable to create codec for %v", t)
		return
	case reflect.Func:
		err = fmt.Errorf("Unable to create codec for %v", t)
		return
	case reflect.Interface:
		err = fmt.Errorf("Unable to create codec for %v", t)
		return
	case reflect.Ptr:
		err = fmt.Errorf("Unable to create codec for %v", t)
		return
	case reflect.Map:
		var keyCodec *codec
		if keyCodec, err = getCodec(t.Key()); err != nil {
			return
		}
		var valueCodec *codec
		if valueCodec, err = getCodec(t.Elem()); err != nil {
			return
		}
		result.generatedEncode = func(w io.Writer, v reflect.Value) (err error) {
			if err = encodeUint(w, uint64(v.Len())); err != nil {
				return
			}
			for _, key := range v.MapKeys() {
				if err = keyCodec.encode(w, key); err != nil {
					return
				}
				value := v.MapIndex(key)
				if err = valueCodec.encode(w, value); err != nil {
					return
				}
			}
			return
		}
		result.generatedDecode = func(r io.Reader, v reflect.Value) (err error) {
			var size uint64
			if size, err = decodeUint(r); err != nil {
				return
			}
			typ := v.Type()
			v.Set(reflect.MakeMap(typ))
			for i := uint64(0); i < size; i++ {
				key := reflect.New(typ.Key())
				if err = keyCodec.decode(r, key); err != nil {
					return
				}
				value := reflect.New(typ.Elem())
				if err = valueCodec.decode(r, value); err != nil {
					return
				}
				v.SetMapIndex(key.Elem(), value.Elem())
			}
			return
		}
	case reflect.Slice:
	case reflect.String:
		result.generatedEncode = func(w io.Writer, v reflect.Value) (err error) {
			b := []byte(v.String())
			if err = encodeUint(w, uint64(len(b))); err != nil {
				return
			}
			_, err = w.Write(b)
			return
		}
		result.generatedDecode = func(r io.Reader, v reflect.Value) (err error) {
			var size uint64
			if size, err = decodeUint(r); err != nil {
				return
			}
			b := make([]byte, int(size))
			if _, err = io.ReadAtLeast(r, b, len(b)); err != nil {
				return
			}
			v.SetString(string(b))
			return
		}
	case reflect.Struct:
	case reflect.UnsafePointer:
	}
	return
}

func getCodec(typ reflect.Type) (result *codec, err error) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	codecLock.RLock()
	result, found := codecByType[typ]
	codecLock.RUnlock()
	if found {
		return
	}
	if result, err = createCodec(typ); err != nil {
		return
	}
	codecLock.Lock()
	codecByType[typ] = result
	codecLock.Unlock()
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
	val := reflect.ValueOf(i)
	c, err := getCodec(val.Type())
	if err != nil {
		return
	}
	return c.encode(self.writer, val)
}

type Decoder struct {
	reader io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader: r,
	}
}

func (self *Decoder) Decode(i interface{}) (err error) {
	val := reflect.ValueOf(i)
	c, err := getCodec(val.Type())
	if err != nil {
		return
	}
	return c.decode(self.reader, val)
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
