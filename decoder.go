package godec

import (
	"encoding"
	"encoding/binary"
	"io"
	"math"
	"reflect"
	"time"
)

type decodeReader struct {
	buf []byte
	pos int
}

func (self *decodeReader) readByte() (result byte, err error) {
	if self.pos >= len(self.buf) {
		err = errorf("EOF")
		return
	}
	result = self.buf[self.pos]
	self.pos++
	return
}

func (self *decodeReader) readBytes(n int) (result []byte, err error) {
	if self.pos+n > len(self.buf) {
		err = errorf("EOF")
		return
	}
	result = self.buf[self.pos : self.pos+n]
	self.pos += n
	return
}

func (self *decodeReader) readUint64() (result uint64, err error) {
	result, read := binary.Uvarint(self.buf[self.pos:])
	if read <= 0 {
		err = errorf("Unable to read uint: %v", read)
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
	if err = decode(r, true, i); err != nil {
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
		err = errorf("Tried to read from %v, got nothing", self.Reader)
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

func decodeType(r *decodeReader) (result *Type, err error) {
	var u uint64
	if err = rawdecodeuint64(r, &u); err != nil {
		return
	}
	result = &Type{
		Base: Kind(u),
	}
	switch result.Base {
	case sliceKind:
		if result.Value, err = decodeType(r); err != nil {
			return
		}
	case mapKind:
		if result.Key, err = decodeType(r); err != nil {
			return
		}
		if result.Value, err = decodeType(r); err != nil {
			return
		}
	}
	return
}

func decodereflect_Value(r *decodeReader, decType bool, v reflect.Value) (err error) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Interface:
		var i interface{}
		if err = decodeinterface__(r, decType, &i); err != nil {
			return
		}
		v.Set(reflect.ValueOf(i))
	case reflect.Bool:
		var b bool
		if err = decodebool(r, decType, &b); err != nil {
			return
		}
		v.SetBool(b)
	case reflect.Int:
		var i int
		if err = decodeint(r, decType, &i); err != nil {
			return
		}
		v.SetInt(int64(i))
	case reflect.Int8:
		var i int8
		if err = decodeint8(r, decType, &i); err != nil {
			return
		}
		v.SetInt(int64(i))
	case reflect.Int16:
		var i int16
		if err = decodeint16(r, decType, &i); err != nil {
			return
		}
		v.SetInt(int64(i))
	case reflect.Int32:
		var i int32
		if err = decodeint32(r, decType, &i); err != nil {
			return
		}
		v.SetInt(int64(i))
	case reflect.Int64:
		var i int64
		if err = decodeint64(r, decType, &i); err != nil {
			return
		}
		v.SetInt(i)
	case reflect.Uint:
		var i uint
		if err = decodeuint(r, decType, &i); err != nil {
			return
		}
		v.SetUint(uint64(i))
	case reflect.Uint8:
		var i uint8
		if err = decodeuint8(r, decType, &i); err != nil {
			return
		}
		v.SetUint(uint64(i))
	case reflect.Uint16:
		var i uint16
		if err = decodeuint16(r, decType, &i); err != nil {
			return
		}
		v.SetUint(uint64(i))
	case reflect.Uint32:
		var i uint32
		if err = decodeuint32(r, decType, &i); err != nil {
			return
		}
		v.SetUint(uint64(i))
	case reflect.Uint64:
		var i uint64
		if err = decodeuint64(r, decType, &i); err != nil {
			return
		}
		v.SetUint(i)
	case reflect.Uintptr:
		var i uintptr
		if err = decodeuintptr(r, decType, &i); err != nil {
			return
		}
		v.SetUint(uint64(i))
	case reflect.Float32:
		var i float32
		if err = decodefloat32(r, decType, &i); err != nil {
			return
		}
		v.SetFloat(float64(i))
	case reflect.Float64:
		var i float64
		if err = decodefloat64(r, decType, &i); err != nil {
			return
		}
		v.SetFloat(i)
	case reflect.Complex64:
		var i complex64
		if err = decodecomplex64(r, decType, &i); err != nil {
			return
		}
		v.SetComplex(complex128(i))
	case reflect.Complex128:
		var i complex128
		if err = decodecomplex128(r, decType, &i); err != nil {
			return
		}
		v.SetComplex(i)
	case reflect.String:
		var i string
		if err = decodestring(r, decType, &i); err != nil {
			return
		}
		v.SetString(i)
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		refType := v.Type()
		if decType {
			var encodedType *Type
			if encodedType, err = decodeType(r); err != nil {
				return
			}
			var valType *Type
			if valType, err = getTypeOf(refType); err != nil {
				return
			}
			if !encodedType.Equal(valType) {
				err = errorf("Can't decode %v into %v", encodedType, valType)
				return
			}
		}
		var l uint
		if err = rawdecodeuint(r, &l); err != nil {
			return
		}
		v.Set(reflect.MakeSlice(refType, int(l), int(l)))
		for i := 0; i < int(l); i++ {
			elType := refType.Elem()
			origEl := reflect.New(elType)
			el := origEl.Elem()
			for ; elType.Kind() == reflect.Ptr; elType = elType.Elem() {
				nextEl := reflect.New(elType.Elem())
				el.Set(nextEl)
				el = el.Elem()
			}
			if err = decode(r, false, el.Addr().Interface()); err != nil {
				return
			}
			v.Index(i).Set(origEl.Elem())
		}
	case reflect.Map:
		refType := v.Type()
		if decType {
			var encodedType *Type
			if encodedType, err = decodeType(r); err != nil {
				return
			}
			var valType *Type
			if valType, err = getTypeOf(refType); err != nil {
				return
			}
			if !encodedType.Equal(valType) {
				err = errorf("Can't decode %v into %v", encodedType, valType)
				return
			}
		}
		var l uint
		if err = rawdecodeuint(r, &l); err != nil {
			return
		}
		v.Set(reflect.MakeMap(refType))
		for i := 0; i < int(l); i++ {
			keyType := refType.Key()
			origKey := reflect.New(keyType)
			key := origKey.Elem()
			for ; keyType.Kind() == reflect.Ptr; keyType = keyType.Elem() {
				nextKey := reflect.New(keyType.Elem())
				key.Set(nextKey)
				key = key.Elem()
			}
			if err = decode(r, false, key.Addr().Interface()); err != nil {
				return
			}
			elType := refType.Elem()
			origEl := reflect.New(elType)
			el := origEl.Elem()
			for ; elType.Kind() == reflect.Ptr; elType = elType.Elem() {
				nextEl := reflect.New(elType.Elem())
				el.Set(nextEl)
				el = el.Elem()
			}
			if err = decode(r, false, el.Addr().Interface()); err != nil {
				return
			}
			v.SetMapIndex(origKey.Elem(), origEl.Elem())
		}
	case reflect.Struct:
		if decType {
			var encodedType *Type
			if encodedType, err = decodeType(r); err != nil {
				return
			}
			if encodedType.Base != mapKind || encodedType.Key.Base != stringKind || encodedType.Value.Base != interface__Kind {
				err = errorf("Can't decode %v into %#v", encodedType, v)
				return
			}
		}
		var l uint
		if err = rawdecodeuint(r, &l); err != nil {
			return
		}
		fieldName := ""
		for i := 0; i < int(l); i++ {
			if err = rawdecodestring(r, &fieldName); err != nil {
				return
			}
			if field := v.FieldByName(fieldName); field.IsValid() {
				for fieldType := field.Type(); field.Kind() == reflect.Ptr; fieldType = field.Type() {
					nextField := reflect.New(fieldType.Elem())
					field.Set(nextField)
					field = field.Elem()
				}
				if err = decode(r, true, field.Addr().Interface()); err != nil {
					return
				}
			}
		}
	}
	return
}

func decodebinary_Unmarshaler(r *decodeReader, decType bool, bu encoding.BinaryUnmarshaler) (err error) {
	if decType {
		var t *Type
		t, err = decodeType(r)
		if err != nil {
			return
		}
		if t.Base != binaryUnMarshalerKind {
			err = errorf("Unable to decode %v into %v", t, bu)
			return
		}
	}
	var size uint
	if err = rawdecodeuint(r, &size); err != nil {
		return
	}
	b, err := r.readBytes(int(size))
	if err != nil {
		return
	}
	return bu.UnmarshalBinary(b)
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

func rawdecodetime_Time(r *decodeReader, t *time.Time) (err error) {
	var i int64
	if err = rawdecodeint64(r, &i); err != nil {
		return
	}
	*t = time.Unix(0, i)
	return
}

// The special case for byte slices is here, and we treat byte slices exactly like strings.
func decodeSliceOfuint8(r *decodeReader, decType bool, v *[]uint8) (err error) {
	if decType {
		var t *Type
		t, err = decodeType(r)
		if err != nil {
			return
		}
		if t.Base != stringKind {
			err = errorf("Unable to decode %v into *[]uint8", t)
			return
		}
	}
	var l uint
	if err = rawdecodeuint(r, &l); err != nil {
		return
	}
	*v, err = r.readBytes(int(l))
	return
}

func decodeInterfaceWithType(r *decodeReader, t *Type, i *interface{}) (err error) {
	switch t.Base {
	case interface__Kind:
		if err = rawdecodeinterface__(r, i); err != nil {
			return
		}
	case boolKind:
		proxy := true
		if err = rawdecodebool(r, &proxy); err != nil {
			return
		}
	case time_TimeKind:
		var proxy time.Time
		if err = rawdecodetime_Time(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case stringKind:
		proxy := ""
		if err = rawdecodestring(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case intKind:
		var proxy int
		if err = rawdecodeint(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case int8Kind:
		proxy := int8(0)
		if err = rawdecodeint8(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case int16Kind:
		proxy := int16(0)
		if err = rawdecodeint16(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case int32Kind:
		proxy := int32(0)
		if err = rawdecodeint32(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case int64Kind:
		proxy := int64(0)
		if err = rawdecodeint64(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case uintptrKind:
		proxy := uintptr(0)
		if err = rawdecodeuintptr(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case uintKind:
		proxy := uint(0)
		if err = rawdecodeuint(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case uint8Kind:
		proxy := uint8(0)
		if err = rawdecodeuint8(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case uint16Kind:
		proxy := uint16(0)
		if err = rawdecodeuint16(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case uint32Kind:
		proxy := uint32(0)
		if err = rawdecodeuint32(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case uint64Kind:
		proxy := uint64(0)
		if err = rawdecodeuint64(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case float32Kind:
		proxy := float32(0)
		if err = rawdecodefloat32(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case float64Kind:
		proxy := float64(0)
		if err = rawdecodefloat64(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case complex64Kind:
		proxy := complex(float32(0), float32(0))
		if err = rawdecodecomplex64(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case complex128Kind:
		proxy := complex(float64(0), float64(0))
		if err = rawdecodecomplex128(r, &proxy); err != nil {
			return
		}
		*i = proxy
	case sliceKind:
		var l uint
		if err = rawdecodeuint(r, &l); err != nil {
			return
		}
		sl := make([]interface{}, int(l))
		for i := 0; i < int(l); i++ {
			var el interface{}
			if err = decodeInterfaceWithType(r, t.Value, &el); err != nil {
				return
			}
			sl[i] = el
		}
		*i = sl
	case mapKind:
		var l uint
		if err = rawdecodeuint(r, &l); err != nil {
			return
		}
		m := make(map[interface{}]interface{}, int(l))
		for i := 0; i < int(l); i++ {
			var k interface{}
			if err = decodeInterfaceWithType(r, t.Key, &k); err != nil {
				return
			}
			var v interface{}
			if err = decodeInterfaceWithType(r, t.Value, &v); err != nil {
				return
			}
			m[k] = v
		}
		*i = m
	default:
		err = errorf("Unknown kind to decode to interface: %v", t)
	}
	return
}

func decodeinterface__(r *decodeReader, decType bool, i *interface{}) (err error) {
	return rawdecodeinterface__(r, i)
}

func rawdecodeinterface__(r *decodeReader, i *interface{}) (err error) {
	t, err := decodeType(r)
	if err != nil {
		return
	}
	return decodeInterfaceWithType(r, t, i)
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
