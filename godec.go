package godec

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
)

const (
	stringKind Kind = iota
	boolKind
	intKind
	uintKind
	time_TimeKind
	interface__Kind
	int8Kind
	int16Kind
	int32Kind
	int64Kind
	uint8Kind
	uint16Kind
	uint32Kind
	uint64Kind
	uintptrKind
	float32Kind
	float64Kind
	complex64Kind
	complex128Kind
	sliceKind
	mapKind
	structKind
	binaryUnMarshalerKind
)

type stackErr struct {
	Err   error
	Stack string
}

func (self stackErr) String() string {
	return self.Err.Error() + "\n" + self.Stack
}

func (self stackErr) Error() string {
	return self.String()
}

func prettify(obj interface{}) string {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", obj)
	}
	return string(b)
}

func errorf(f string, args ...interface{}) error {
	stack := make([]byte, 4096)
	s := runtime.Stack(stack, false)
	if s < len(stack) {
		stack = stack[:s]
	}
	return stackErr{
		Err:   fmt.Errorf(f, args...),
		Stack: string(stack),
	}
}

type Type struct {
	Base  Kind
	Key   *Type
	Value *Type
}

func (self Type) String() string {
	switch self.Base {
	case sliceKind:
		return "[]" + self.Value.String()
	case mapKind:
		return "map[" + self.Key.String() + "]" + self.Value.String()
	default:
		return self.Base.String()
	}
}

func (self *Type) Equal(o *Type) bool {
	if self == nil && o == nil {
		return true
	}
	return self.Base == o.Base && self.Key.Equal(o.Key) && self.Value.Equal(o.Value)
}

type Kind uint64

func (self Kind) String() string {
	switch self {
	case interface__Kind:
		return "interface{}"
	case stringKind:
		return "string"
	case boolKind:
		return "bool"
	case intKind:
		return "int"
	case uintKind:
		return "uint"
	case int8Kind:
		return "int8"
	case int16Kind:
		return "int16"
	case int32Kind:
		return "int32"
	case int64Kind:
		return "int64"
	case uint8Kind:
		return "uint8"
	case uint16Kind:
		return "uint16"
	case uint32Kind:
		return "uint32"
	case uint64Kind:
		return "uint64"
	case uintptrKind:
		return "uintptr"
	case float32Kind:
		return "float32"
	case float64Kind:
		return "float64"
	case complex64Kind:
		return "complex64"
	case complex128Kind:
		return "complex128"
	case time_TimeKind:
		return "time.Time"
	case sliceKind:
		return "[]"
	case mapKind:
		return "map[]"
	case binaryUnMarshalerKind:
		return "encoding.BinaryUnMarshaler"
	}
	return fmt.Sprintf("Unrecognized Kind %#v", self)
}

func getTypeOf(t reflect.Type) (result *Type, err error) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Bool:
		result = &Type{Base: boolKind}
	case reflect.Int:
		result = &Type{Base: int64Kind}
	case reflect.Int8:
		result = &Type{Base: int8Kind}
	case reflect.Int16:
		result = &Type{Base: int16Kind}
	case reflect.Int32:
		result = &Type{Base: int64Kind}
	case reflect.Int64:
		result = &Type{Base: int64Kind}
	case reflect.Uint:
		result = &Type{Base: uintKind}
	case reflect.Uint8:
		result = &Type{Base: uint8Kind}
	case reflect.Uint16:
		result = &Type{Base: uint16Kind}
	case reflect.Uint32:
		result = &Type{Base: uint32Kind}
	case reflect.Uint64:
		result = &Type{Base: uint64Kind}
	case reflect.Uintptr:
		result = &Type{Base: uintptrKind}
	case reflect.Float32:
		result = &Type{Base: float32Kind}
	case reflect.Float64:
		result = &Type{Base: float64Kind}
	case reflect.Complex64:
		result = &Type{Base: complex64Kind}
	case reflect.Complex128:
		result = &Type{Base: complex128Kind}
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		var valueType *Type
		if valueType, err = getTypeOf(t.Elem()); err != nil {
			return
		}
		result = &Type{Base: sliceKind, Value: valueType}
	case reflect.Map:
		var keyType *Type
		if keyType, err = getTypeOf(t.Key()); err != nil {
			return
		}
		var valueType *Type
		if valueType, err = getTypeOf(t.Elem()); err != nil {
			return
		}
		result = &Type{Base: mapKind, Key: keyType, Value: valueType}
	case reflect.String:
		result = &Type{Base: stringKind}
	case reflect.Struct:
		result = &Type{Base: mapKind, Key: &Type{Base: stringKind}, Value: &Type{Base: interface__Kind}}
	case reflect.Interface:
		result = &Type{Base: interface__Kind}
	default:
		err = errorf("Unable to encode %v", t)
	}
	return
}
