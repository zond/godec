package godec

import (
	"encoding"
	"encoding/binary"
	"encoding/gob"
	"io"
	"math"
	"reflect"
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
	if err = encodeinterface__(w, true, i, nil); err != nil {
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

func encodereflect_Value(w *encodeWriter, encType bool, v reflect.Value) (err error) {
	for k := v.Kind(); k == reflect.Ptr || k == reflect.Interface; k = v.Kind() {
		v = v.Elem()
	}
	var typ *Type
	refType := v.Type()
	if typ, err = getTypeOf(refType); err != nil {
		return
	}
	if encType {
		if err = encodeType(w, typ); err != nil {
			return
		}
	}
	switch v.Kind() {
	case reflect.Bool:
		err = rawencodebool(w, v.Bool())
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		err = rawencodeint64(w, v.Int())
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
		err = rawencodeuint64(w, v.Uint())
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		err = rawencodefloat64(w, v.Float())
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		err = rawencodecomplex128(w, v.Complex())
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		if err = rawencodeuint(w, uint(v.Len())); err != nil {
			return
		}
		for i := 0; i < v.Len(); i++ {
			el := v.Index(i)
			for el.Kind() == reflect.Ptr {
				el = el.Elem()
			}
			if err = encodeinterface__(w, typ.Value.Base == interface__Kind, el.Interface(), &el); err != nil {
				return
			}
		}
	case reflect.Map:
		if err = rawencodeuint(w, uint(v.Len())); err != nil {
			return
		}
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			for key.Kind() == reflect.Ptr {
				key = key.Elem()
			}
			for value.Kind() == reflect.Ptr {
				value = value.Elem()
			}
			if err = encodeinterface__(w, typ.Key.Base == interface__Kind, key.Interface(), &key); err != nil {
				return
			}
			if err = encodeinterface__(w, typ.Value.Base == interface__Kind, value.Interface(), &value); err != nil {
				return
			}
		}
	case reflect.String:
		err = rawencodestring(w, v.String())
	case reflect.Struct:
		names := []string{}
		values := []reflect.Value{}
		collectStructValues(v, refType, &names, &values)
		if err = rawencodeuint(w, uint(len(names))); err != nil {
			return
		}
		for index, name := range names {
			if err = rawencodestring(w, name); err != nil {
				return
			}
			if err = encodeinterface__(w, true, values[index].Interface(), &values[index]); err != nil {
				return
			}
		}
	default:
		err = errorf("Unable to encode %v", v.Interface())
	}
	return
}

func collectStructValues(v reflect.Value, refType reflect.Type, names *[]string, values *[]reflect.Value) {
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		if field.PkgPath == "" {
			if field.Anonymous {
				collectStructValues(v.Field(i), field.Type, names, values)
			} else {
				val := v.Field(i)
				for val.Kind() == reflect.Ptr {
					val = val.Elem()
				}
				if val.IsValid() {
					*names = append(*names, field.Name)
					*values = append(*values, val)
				}
			}
		}
	}
}

func rawencodetime_Time(w *encodeWriter, t time.Time) (err error) {
	b, err := t.MarshalBinary()
	if err != nil {
		return
	}
	if err = rawencodeuint(w, uint(len(b))); err != nil {
		return
	}
	return w.writeBytes(b)
}

func encodegob_GobEncoder(w *encodeWriter, encType bool, ge gob.GobEncoder) (err error) {
	if encType {
		if err = encodeType(w, &Type{Base: gobDEncoderKind}); err != nil {
			return
		}
	}
	b, err := ge.GobEncode()
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

func encodebinary_Marshaler(w *encodeWriter, encType bool, bm encoding.BinaryMarshaler) (err error) {
	if encType {
		if err = encodeType(w, &Type{Base: binaryUnMarshalerKind}); err != nil {
			return
		}
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

func rawencodestring(w *encodeWriter, s string) (err error) {
	if err = rawencodeuint64(w, uint64(len(s))); err != nil {
		return
	}
	err = w.writeString(s)
	return
}

// The special case for byte slices is here, and we treat byte slices exactly like strings.
func encodeSliceOfuint8(w *encodeWriter, encType bool, v []uint8, val *reflect.Value) (err error) {
	if encType {
		if err = encodeType(w, &Type{Base: stringKind}); err != nil {
			return
		}
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
