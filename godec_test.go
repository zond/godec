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

func DeepEqual(i1, i2 interface{}) bool {
	return deepEqual(reflect.ValueOf(i1), reflect.ValueOf(i2))
}

func deepEqual(val1, val2 reflect.Value) bool {
	for k := val1.Kind(); k == reflect.Ptr || k == reflect.Interface; k = val1.Kind() {
		val1 = val1.Elem()
	}
	for k := val2.Kind(); k == reflect.Ptr || k == reflect.Interface; k = val2.Kind() {
		val2 = val2.Elem()
	}
	if val1.Kind() != val2.Kind() {
		return false
	}
	switch val1.Kind() {
	case reflect.Bool:
		return val1.Bool() == val2.Bool()
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return val1.Int() == val2.Int()
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
		return val1.Uint() == val2.Uint()
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return val1.Float() == val2.Float()
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return val1.Complex() == val2.Complex()
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if val1.Len() != val2.Len() {
			return false
		}
		for i := 0; i < val1.Len(); i++ {
			if !deepEqual(val1.Index(i), val2.Index(i)) {
				return false
			}
		}
	case reflect.Map:
		if val1.Len() != val2.Len() {
			return false
		}
		mapType := val1.Type()
		keyType := mapType.Key()
		for keyType.Kind() == reflect.Ptr {
			keyType = keyType.Elem()
		}
		fake1 := reflect.MakeMap(reflect.MapOf(keyType, reflect.SliceOf(mapType.Elem())))
		for _, key := range val1.MapKeys() {
			origKey := key
			for key.Kind() == reflect.Ptr {
				key = key.Elem()
			}
			var fakeVals reflect.Value
			if fakeVals = fake1.MapIndex(key); fakeVals.IsValid() {
				fakeVals = reflect.Append(fakeVals, val1.MapIndex(origKey))
			} else {
				fakeVals = reflect.MakeSlice(reflect.SliceOf(mapType.Elem()), 0, 1)
				fakeVals = reflect.Append(fakeVals, val1.MapIndex(origKey))
			}
			fake1.SetMapIndex(key, fakeVals)
		}
		for _, key := range val2.MapKeys() {
			value := val2.MapIndex(key)
			for key.Kind() == reflect.Ptr {
				key = key.Elem()
			}
			fake1Vals := fake1.MapIndex(key)
			if !fake1Vals.IsValid() {
				return false
			}
			found := false
			for i := 0; i < fake1Vals.Len(); i++ {
				if deepEqual(value, fake1Vals.Index(i)) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	case reflect.String:
		return val1.String() == val2.String()
	case reflect.Struct:
		if val1.Type() != val2.Type() {
			return false
		}
		switch i := val1.Interface().(type) {
		case time.Time:
			i2 := val2.Interface().(time.Time)
			return i.Equal(i2)
		default:
			typ := val1.Type()
			for i := 0; i < typ.NumField(); i++ {
				if !deepEqual(val1.Field(i), val2.Field(i)) {
					return false
				}
			}
		}
	default:
		return false
	}
	return true
}

func encodeDecode(t *testing.T, src, dst interface{}) {
	encodeDecodeWithCMP(t, src, src, dst)
}

func encodeDecodeWithCMP(t *testing.T, src, cmp, dst interface{}) {
	buf := &bytes.Buffer{}
	err := NewEncoder(buf).Encode(src)
	if err != nil {
		t.Fatalf("Unable to encode %v: %v", src, err)
	}
	if err = NewDecoder(buf).Decode(dst); err != nil {
		t.Fatalf("Unable to decode to %v: %v", reflect.ValueOf(dst).Elem().Interface(), err)
	}
	if !DeepEqual(cmp, dst) {
		t.Fatalf("Encoding/decoding %v produced %v", prettify(cmp), prettify(dst))
	}
	b, err := Marshal(src)
	if err != nil {
		t.Fatalf("Unable to marshal %v: %v", src, err)
	}
	if err = Unmarshal(b, dst); err != nil {
		t.Fatalf("Unable to unmarshal to %v: %v", dst, err)
	}
	if !DeepEqual(cmp, dst) {
		t.Fatalf("Marshalling/unmarshalling %v produced %v", prettify(cmp), prettify(dst))
	}
}

func TestManualEncodeDecodeMapOfInterface__PtrToInterface__(t *testing.T) {
	var dst map[*interface{}]interface{}
	var a interface{}
	a = 10
	encodeDecode(t, map[*interface{}]interface{}{
		&a: "a",
	}, &dst)
}

type nestedThing1 map[string][]int
type nestedThing2 map[interface{}]map[int]interface{}

func TestManualEncodeDecodeNestedTypedThingsToInterfaces(t *testing.T) {
	var dst interface{}
	encodeDecodeWithCMP(t, nestedThing1{
		"a": []int{1, 2, 3},
		"b": []int{4, 5, 6},
	}, map[interface{}]interface{}{
		"a": []interface{}{
			int64(1), int64(2), int64(3),
		},
		"b": []interface{}{
			int64(4), int64(5), int64(6),
		},
	}, &dst)
	encodeDecodeWithCMP(t, nestedThing2{
		"gna": {
			4: "hehu",
		},
		int64(42): {
			5: int64(44),
		},
	}, map[interface{}]interface{}{
		"gna": map[interface{}]interface{}{
			int64(4): "hehu",
		},
		int64(42): map[interface{}]interface{}{
			int64(5): int64(44),
		},
	}, &dst)
}

type structThing1 struct {
	A int
	B string
}
type structThing2 struct {
	A int
	B string
	C *structThing1
}

func TestManualEncodeDecodeStructTypesToInterfaces(t *testing.T) {
	var dst interface{}
	encodeDecodeWithCMP(t, structThing1{
		A: 33,
		B: "hehu",
	}, map[interface{}]interface{}{
		"A": 33,
		"B": "hehu",
	}, &dst)
	var dst2 map[string]interface{}
	encodeDecodeWithCMP(t, structThing1{
		A: 33,
		B: "hehu",
	}, map[string]interface{}{
		"A": 33,
		"B": "hehu",
	}, &dst2)
	encodeDecodeWithCMP(t, &structThing2{
		A: 41,
		B: "blapp",
		C: &structThing1{
			A: 11,
			B: "JAJA",
		},
	}, map[interface{}]interface{}{
		"A": 41,
		"B": "blapp",
		"C": map[interface{}]interface{}{
			"A": 11,
			"B": "JAJA",
		},
	}, &dst)
}

func TestManualEncodeDecodeNestedInterfaces(t *testing.T) {
	var dst interface{}
	encodeDecode(t, map[interface{}]interface{}{
		"hej": 33,
		41: []interface{}{
			"apa",
			"gnu",
		},
	}, &dst)
	encodeDecode(t, map[interface{}]interface{}{
		"hehu": map[interface{}]interface{}{
			33: []interface{}{
				31,
				"hehu",
				map[interface{}]interface{}{
					11: map[interface{}]interface{}{
						"gngn": "hepp",
					},
				},
			},
		},
	}, &dst)
	encodeDecode(t, []interface{}{
		"333",
		333,
		map[interface{}]interface{}{
			"apa": 11,
		},
		[]interface{}{
			[]interface{}{
				22,
				44,
				"hepp",
			},
			"hepp2",
		},
	}, &dst)
}

type sliceType1 []string
type sliceType2 [][]int
type sliceType3 []*string

func TestManualEncodeDecodeSliceTypes(t *testing.T) {
	var dst sliceType1
	encodeDecode(t, sliceType1{
		"a",
		"b",
	}, &dst)
	var dst2 sliceType2
	encodeDecode(t, sliceType2{
		[]int{
			3, 4, 5,
		},
		[]int{
			6, 7, 8,
		},
	}, &dst2)
	s1 := "a"
	s2 := "b"
	var dst3 sliceType3
	encodeDecode(t, sliceType3{
		&s1,
		&s2,
	}, &dst3)
}

type mapType1 map[string]int
type mapType2 map[int]map[string]bool
type mapType3 map[int]map[*string]*string

func TestManualEncodeDecodeMapTypes(t *testing.T) {
	var dst1 mapType1
	encodeDecode(t, mapType1{
		"a": 1,
		"b": 2,
	}, &dst1)
	var dst2 mapType2
	encodeDecode(t, mapType2{
		3: {
			"a": true,
			"b": false,
		},
		5: {
			"c": true,
			"d": true,
		},
	}, &dst2)
	a := "a"
	b := "b"
	c := "c"
	d := "d"
	var dst3 mapType3
	encodeDecode(t, mapType3{
		1: {
			&a: &b,
			&c: &d,
		},
		2: {
			&b: &a,
			&d: &c,
		},
	}, &dst3)
}
