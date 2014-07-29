package godec

import (
	"fmt"
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
