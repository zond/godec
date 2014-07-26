package godec

import (
	"bytes"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var Runs = 1

func init() {
	envRuns, err := strconv.Atoi(os.Getenv("RUNS"))
	if err == nil {
		Runs = envRuns
	}
}

func randombool() bool {
	if rand.Int31n(2) == 0 {
		return true
	}
	return false
}

func randominterface__() interface{} {
	switch randomuint8() % 16 {
	case 0:
		return randomstring()
	case 1:
		return randomfloat64()
	case 2:
		return randomfloat32()
	case 3:
		return randomcomplex64()
	case 4:
		return randomcomplex128()
	case 5:
		return randomuintptr()
	case 6:
		return randomuint()
	case 7:
		return randomint()
	case 8:
		return randomint8()
	case 9:
		return randomint16()
	case 10:
		return randomint32()
	case 11:
		return randomint64()
	case 12:
		return randomuint8()
	case 13:
		return randomuint16()
	case 14:
		return randomuint32()
	case 15:
		return randomuint64()
	}
	panic("wtf")
}

func randomstring() string {
	l := rand.Int31() % 4096
	b := &bytes.Buffer{}
	for i := int32(0); i < l; i++ {
		b.Write([]byte{byte(rand.Int31())})
	}
	return b.String()
}

func randomfloat64() float64 {
	return float64(rand.Int63()) + rand.NormFloat64()
}

func randomfloat32() float32 {
	return float32(randomfloat64())
}

func randomcomplex64() complex64 {
	return complex(randomfloat32(), randomfloat32())
}

func randomcomplex128() complex128 {
	return complex(randomfloat64(), randomfloat64())
}

func randomuintptr() uintptr {
	return uintptr(rand.Int63())
}

func randomuint() uint {
	return uint(rand.Int63())
}

func randomint() int {
	return int(rand.Int63())
}

func randomint8() int8 {
	return int8(rand.Int63())
}

func randomint16() int16 {
	return int16(rand.Int63())
}

func randomint32() int32 {
	return int32(rand.Int63())
}

func randomint64() int64 {
	return int64(rand.Int63())
}

func randomuint8() uint8 {
	return uint8(rand.Int63())
}

func randomuint16() uint16 {
	return uint16(rand.Int63())
}

func randomuint32() uint32 {
	return uint32(rand.Int63())
}

func randomuint64() uint64 {
	return uint64(rand.Int63())
}

func randomtime_Time() time.Time {
	return time.Unix(0, randomint64())
}

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
	toCmp := src
	if srcVal := reflect.ValueOf(src); srcVal.Kind() == reflect.Ptr {
		toCmp = srcVal.Elem().Interface()
	}
	if !reflect.DeepEqual(toCmp, dstElem) {
		t.Fatalf("Encoding/decoding %v produced %v", toCmp, dstElem)
	}
	b, err := Marshal(src)
	if err != nil {
		t.Fatalf("Unable to marshal %v: %v", src, err)
	}
	if err = Unmarshal(b, dst); err != nil {
		t.Fatalf("Unable to unmarshal to %v: %v", dstElem, err)
	}
	dstElem = reflect.ValueOf(dst).Elem().Interface()
	if !reflect.DeepEqual(toCmp, dstElem) {
		t.Fatalf("Marshalling/unmarshalling %v produced %v", toCmp, dstElem)
	}
}

func TestEncodeDecodeMaps(t *testing.T) {
	var mSI map[string]int
	encodeDecode(t, map[string]int{"hej": 1, "haha": 2}, &mSI)
	encodeDecode(t, map[string]int{"hejsan": 10, "hahaha": 20}, &mSI)
	var mIB map[int]bool
	encodeDecode(t, map[int]bool{44: true, 52: false}, &mIB)
	encodeDecode(t, map[int]bool{442: true, 523: true}, &mIB)
}
