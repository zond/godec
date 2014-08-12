package godec

import (
	"bytes"
	"fmt"
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

func DeepEqual(i1, i2 interface{}) error {
	return deepEqual(reflect.ValueOf(i1), reflect.ValueOf(i2))
}

func deepEqual(val1, val2 reflect.Value) error {
	for k := val1.Kind(); k == reflect.Ptr || k == reflect.Interface; k = val1.Kind() {
		val1 = val1.Elem()
	}
	for k := val2.Kind(); k == reflect.Ptr || k == reflect.Interface; k = val2.Kind() {
		val2 = val2.Elem()
	}
	if !val1.IsValid() {
		if !val2.IsValid() {
			return nil
		} else {
			return fmt.Errorf("%v != %v", val1, val2)
		}
	} else if !val2.IsValid() {
		return fmt.Errorf("%v != %v", val1, val2)
	}
	if val1.Kind() != val2.Kind() {
		return fmt.Errorf("%+v != %+v (val1.Kind() = %v, val2.Kind() = %v)", val1.Interface(), val2.Interface(), val1.Kind(), val2.Kind())
	}
	switch val1.Kind() {
	case reflect.Bool:
		if val1.Bool() != val2.Bool() {
			return fmt.Errorf("%+v != %+v", val1.Bool(), val2.Bool())
		}
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		if val1.Int() != val2.Int() {
			return fmt.Errorf("%+v != %+v", val1.Int(), val2.Int())
		}
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
		if val1.Uint() != val2.Uint() {
			return fmt.Errorf("%+v != %+v", val1.Uint(), val2.Uint())
		}
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		if val1.Float() != val2.Float() {
			return fmt.Errorf("%+v != %+v", val1.Float(), val2.Float())
		}
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		if val1.Complex() != val2.Complex() {
			return fmt.Errorf("%+v != %+v", val1.Complex(), val2.Complex())
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if val1.Len() != val2.Len() {
			return fmt.Errorf("%+v != %+v (val1.Len() = %v, val2.Len() = %v)", val1.Interface(), val2.Interface(), val1.Len(), val2.Len())
		}
		for i := 0; i < val1.Len(); i++ {
			if err := deepEqual(val1.Index(i), val2.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		if val1.Len() != val2.Len() {
			return fmt.Errorf("%+v != %+v (val1.Len() = %v, val2.Len() = %v)", val1.Interface(), val2.Interface(), val1.Len(), val2.Len())
		}
		mapType := val1.Type()
		keyType := mapType.Key()
		for keyType.Kind() == reflect.Ptr {
			keyType = keyType.Elem()
		}
		fake1 := reflect.MakeMap(reflect.MapOf(keyType, reflect.SliceOf(mapType.Elem())))
		var nilValue reflect.Value
		var nilKey reflect.Value
		for _, key := range val1.MapKeys() {
			origKey := key
			for key.Kind() == reflect.Ptr {
				key = key.Elem()
			}
			if key.IsValid() {
				var fakeVals reflect.Value
				if fakeVals = fake1.MapIndex(key); fakeVals.IsValid() {
					fakeVals = reflect.Append(fakeVals, val1.MapIndex(origKey))
				} else {
					fakeVals = reflect.MakeSlice(reflect.SliceOf(mapType.Elem()), 0, 1)
					fakeVals = reflect.Append(fakeVals, val1.MapIndex(origKey))
				}
				fake1.SetMapIndex(key, fakeVals)
			} else {
				nilKey = origKey
				nilValue = val1.MapIndex(origKey)
			}
		}
		for _, key := range val2.MapKeys() {
			value := val2.MapIndex(key)
			for key.Kind() == reflect.Ptr {
				key = key.Elem()
			}
			if key.IsValid() {
				fake1Vals := fake1.MapIndex(key)
				if !fake1Vals.IsValid() {
					return fmt.Errorf("%+v != %+v (%v not found in former)", val1.Interface(), val2.Interface(), key.Interface())
				}
				found := false
				errors := []error{}
				for i := 0; i < fake1Vals.Len(); i++ {
					if err := deepEqual(value, fake1Vals.Index(i)); err == nil {
						found = true
						break
					} else {
						errors = append(errors, err)
					}
				}
				if !found {
					return fmt.Errorf("%+v != %+v:\n%v", val1.Interface(), val2.Interface(), errors)
				}
			}
		}
		if nilKey.IsValid() {
			value := val2.MapIndex(nilKey)
			if value.IsValid() {
				if err := deepEqual(value, nilValue); err != nil {
					return fmt.Errorf("%+v != %+v (nil key in val1 != nil key in val1 != nil key in val2: %v)", val1.Interface(), val2.Interface(), err)
				}
			} else {
				return fmt.Errorf("%+v != %+v (no nil key in latter)", val1.Interface(), val2.Interface())
			}
		}
	case reflect.String:
		if val1.String() != val2.String() {
			return fmt.Errorf("%#v != %#v", val1.String(), val2.String())
		}
	case reflect.Struct:
		if val1.Type() != val2.Type() {
			return fmt.Errorf("%+v != %+v (val1.Type() = %v, val2.Type() = %v)", val1.Interface(), val2.Interface(), val1.Type(), val2.Type())
		}
		switch i := val1.Interface().(type) {
		case time.Time:
			i2 := val2.Interface().(time.Time)
			if !i.Equal(i2) {
				return fmt.Errorf("%+v != %+v", i, i2)
			}
		default:
			typ := val1.Type()
			for i := 0; i < typ.NumField(); i++ {
				if val1.Field(i).IsValid() && val2.Field(i).IsValid() {
					if err := deepEqual(val1.Field(i), val2.Field(i)); err != nil {
						return err
					}
				} else if val1.Field(i).IsValid() || val2.Field(i).IsValid() {
					return fmt.Errorf("%+v != %+v (val1.%v.IsValid() = %v, val2.%v.IsValid() = %v)", val1.Interface(), val2.Interface(), typ.Field(i).Name, val1.Field(i).IsValid(), typ.Field(i).Name, val2.Field(i).IsValid())
				}
			}
		}
	default:
		return fmt.Errorf("Unrecognized kind %v", val1.Kind())
	}
	return nil
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
	if err := DeepEqual(cmp, dst); err != nil {
		t.Fatalf("Encoding/decoding\n%v\nproduced\n%v\n%v", cmp, dst, err)
	}
	b, err := Marshal(src)
	if err != nil {
		t.Fatalf("Unable to marshal %v: %v", src, err)
	}
	if err = Unmarshal(b, dst); err != nil {
		t.Fatalf("Unable to unmarshal to %v: %v", dst, err)
	}
	if err := DeepEqual(cmp, dst); err != nil {
		t.Fatalf("Marshalling/unmarshalling\n%v\nproduced\n%v\n", cmp, dst, err)
	}
}

func TestManualEncodeDecodeEmptyInterface(t *testing.T) {
	var dst interface{}
	var src interface{}
	encodeDecode(t, src, &dst)
}

func TestManualEncodeDecodeSliceOfNilstringPtr(t *testing.T) {
	var dst []*string
	src := make([]*string, 5)
	encodeDecode(t, src, &dst)
}

type nilStruct struct {
	A *string
}

func TestManualEncodeDecodeStructWithNilFields(t *testing.T) {
	var dst nilStruct
	encodeDecode(t, nilStruct{}, &dst)
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

func TestManualEncodeDecodeFlatStructTypes(t *testing.T) {
	var dst structThing1
	encodeDecode(t, structThing1{
		A: 33,
		B: "hehu",
	}, &dst)
	var dst2 structThing2
	encodeDecode(t, &structThing2{
		A: 41,
		B: "blapp",
		C: &structThing1{
			A: 11,
			B: "JAJA",
		},
	}, &dst2)
}

type testStructT1 struct {
	A int
	B string
	C []byte
	D time.Time
	E []string
	F byte
}

type testStructT2 struct {
	testStructT1
	G testStructT1
	H *testStructT2
	I *[]int
	J map[string]*map[int]int
	K map[int]*testStructT2
}

var testStruct1 = testStructT1{
	A: 4412,
	B: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	C: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
	D: time.Unix(1407311220, 0),
	E: []string{"aaaaaaaaaaa", "bbbbbbbbbbbb", "ccccccccccccc", "dddddddddddddd", "eeeeeeeeeeeeee", "fffffffffffffff"},
	F: 43,
}

var testStruct3 = testStructT2{
	testStructT1: testStruct1,
	G:            testStruct1,
	I:            &smallIntSlice,
	J: map[string]*map[int]int{
		"a": &smallIntMap,
		"b": &smallIntMap,
	},
}

var testStruct2 = testStructT2{
	testStructT1: testStruct1,
	G:            testStruct1,
	H:            &testStruct3,
	I:            &smallIntSlice,
	J: map[string]*map[int]int{
		"a": &smallIntMap,
		"b": &smallIntMap,
	},
	K: map[int]*testStructT2{
		44: &testStruct3,
		91: &testStruct3,
	},
}

func TestManualEncodeDecodeNestedStructTypes(t *testing.T) {
	var dst testStructT2
	encodeDecode(t, testStruct2, &dst)
}

func TestManualEncodeDecodeNestedStructTypesToInterfaces(t *testing.T) {
	var dst interface{}
	t1 := map[interface{}]interface{}{
		"A": 4412,
		"B": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"C": string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}),
		"D": time.Unix(1407311220, 0),
		"E": []string{"aaaaaaaaaaa", "bbbbbbbbbbbb", "ccccccccccccc", "dddddddddddddd", "eeeeeeeeeeeeee", "fffffffffffffff"},
		"F": uint8(43),
	}
	sl1 := []interface{}{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	}
	m1 := map[interface{}]interface{}{
		int64(1):  int64(4),
		int64(4):  int64(5),
		int64(7):  int64(1),
		int64(8):  int64(11),
		int64(91): int64(32),
	}
	t3 := map[interface{}]interface{}{
		"A": 4412,
		"B": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"C": string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}),
		"D": time.Unix(1407311220, 0),
		"E": []string{"aaaaaaaaaaa", "bbbbbbbbbbbb", "ccccccccccccc", "dddddddddddddd", "eeeeeeeeeeeeee", "fffffffffffffff"},
		"F": uint8(43),
		"G": t1,
		"I": sl1,
		"J": map[interface{}]interface{}{
			"a": &m1,
			"b": &m1,
		},
		"K": map[interface{}]interface{}{},
	}
	t2 := map[interface{}]interface{}{
		"A": 4412,
		"B": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"C": string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}),
		"D": time.Unix(1407311220, 0),
		"E": []string{"aaaaaaaaaaa", "bbbbbbbbbbbb", "ccccccccccccc", "dddddddddddddd", "eeeeeeeeeeeeee", "fffffffffffffff"},
		"F": uint8(43),
		"G": t1,
		"H": t3,
		"I": sl1,
		"J": map[interface{}]interface{}{
			"a": m1,
			"b": m1,
		},
		"K": map[interface{}]interface{}{
			int64(44): t3,
			int64(91): t3,
		},
	}
	encodeDecodeWithCMP(t, testStruct2, t2, &dst)
}

func TestManualEncodeDecodeFlatStructTypesToInterfaces(t *testing.T) {
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
