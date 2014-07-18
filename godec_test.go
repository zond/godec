package godec

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
)

func encodeDecode(t *testing.T, src, dst interface{}) {
	buf := &bytes.Buffer{}
	err := NewEncoder(buf).Encode(src)
	if err != nil {
		t.Fatalf("Unable to encode %v: %v", src, err)
	}
	if err = NewDecoder(buf).Decode(dst); err != nil {
		t.Fatalf("Unable to decode to %v: %v", reflect.ValueOf(dst).Elem().Interface(), err)
	}
	dstElem := reflect.ValueOf(dst).Elem().Interface()
	if !reflect.DeepEqual(src, dstElem) {
		t.Fatalf("Encoding/decoding %v produced %v", src, dstElem)
	}
	b, err := Marshal(src)
	if err != nil {
		t.Fatalf("Unable to marshal %v: %v", src, err)
	}
	if err = Unmarshal(b, dst); err != nil {
		t.Fatalf("Unable to unmarshal to %v: %v", dstElem, err)
	}
	dstElem = reflect.ValueOf(dst).Elem().Interface()
	if !reflect.DeepEqual(src, dstElem) {
		t.Fatalf("Marshalling/unmarshalling %v produced %v", src, dstElem)
	}
}

func TestEncodeDecodePrimitiveTypes(t *testing.T) {
	var i8 int8
	encodeDecode(t, int8(5), &i8)
	encodeDecode(t, int8(84), &i8)
	var i16 int16
	encodeDecode(t, int16(5), &i16)
	encodeDecode(t, int16(84), &i16)
	var i32 int32
	encodeDecode(t, int32(5), &i32)
	encodeDecode(t, int32(84), &i32)
	var i64 int64
	encodeDecode(t, int64(5), &i64)
	encodeDecode(t, int64(84), &i64)
	var f32 float32
	encodeDecode(t, float32(5.56), &f32)
	encodeDecode(t, float32(8484.0004), &f32)
	var f64 float64
	encodeDecode(t, float64(5.56), &f64)
	encodeDecode(t, float64(8484.0004), &f64)
	var c64 complex64
	encodeDecode(t, complex(float32(5.56), float32(5.1)), &c64)
	encodeDecode(t, complex(float32(4.11), float32(63.11)), &c64)
	var c128 complex128
	encodeDecode(t, complex(float64(5.56), float64(5.1)), &c128)
	encodeDecode(t, complex(float64(4.11), float64(63.11)), &c128)
	var b bool
	encodeDecode(t, true, &b)
	encodeDecode(t, false, &b)
}

func TestEncodeDecodeMaps(t *testing.T) {
	var mSI map[string]int
	encodeDecode(t, map[string]int{"hej": 1, "haha": 2}, &mSI)
	encodeDecode(t, map[string]int{"hejsan": 10, "hahaha": 20}, &mSI)
	var mIB map[int]bool
	encodeDecode(t, map[int]bool{44: true, 52: false}, &mIB)
	encodeDecode(t, map[int]bool{442: true, 523: true}, &mIB)
}

func TestEncodeDecodeSlices(t *testing.T) {
	var sS []string
	encodeDecode(t, []string{"hehu", "hepp"}, &sS)
	encodeDecode(t, []string{"kala", "bapa"}, &sS)
	var sI32 []int32
	encodeDecode(t, []int32{33, 44, 56, 1}, &sI32)
	encodeDecode(t, []int32{1, 2, 3}, &sI32)
}

func TestEncodeDecodeUint64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		i1 := uint64(rand.Int63())
		b, err := Marshal(i1)
		if err != nil {
			t.Fatalf("%v", err)
		}
		var i2 uint64
		if err := Unmarshal(b, &i2); err != nil {
			t.Fatalf("%v", err)
		}
		if i1 != i2 {
			t.Fatalf("Encoded %v to %v, and decoded that to %v", i1, b, i2)
		}
	}
}

func TestEncodeDecodeInt64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		i1 := rand.Int63()
		b, err := Marshal(i1)
		if err != nil {
			t.Fatalf("%v", err)
		}
		var i2 int64
		if err := Unmarshal(b, &i2); err != nil {
			t.Fatalf("%v", err)
		}
		if i1 != i2 {
			t.Fatalf("Encoded %v to %v, and decoded that to %v", i1, b, i2)
		}
	}
}
