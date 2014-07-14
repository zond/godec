package godec

import (
	"bytes"
	"reflect"
)

const (
	reflect_ValueKind = reflect.Struct
	interface__Kind   = reflect.Interface
	stringKind        = reflect.String
	boolKind          = reflect.Uint64
	intKind           = reflect.Int64
	int8Kind          = reflect.Int64
	int16Kind         = reflect.Int64
	int32Kind         = reflect.Int64
	int64Kind         = reflect.Int64
	uintKind          = reflect.Uint64
	uint8Kind         = reflect.Uint64
	uint16Kind        = reflect.Uint64
	uint32Kind        = reflect.Uint64
	uint64Kind        = reflect.Uint64
	uintptrKind       = reflect.Uint64
	float32Kind       = reflect.Float64
	float64Kind       = reflect.Float64
	complex64Kind     = reflect.Complex128
	complex128Kind    = reflect.Complex128
)

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
	dec := NewDecoder(&BytesDecodeReader{Buf: b})
	if err = dec.Decode(i); err != nil {
		return
	}
	return
}
