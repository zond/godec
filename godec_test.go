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

func convertPtrMap(v reflect.Value) (result reflect.Value) {
	trueValue := reflect.ValueOf(true)
	boolType := trueValue.Type()
	valueType := v.Type().Elem()
	for valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}
	keyType := v.Type().Key()
	for keyType.Kind() == reflect.Ptr {
		keyType = keyType.Elem()
	}
	containerType := reflect.MapOf(valueType, boolType)
	result = reflect.MakeMap(reflect.MapOf(keyType, containerType))
	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)
		for key.Kind() == reflect.Ptr {
			key = key.Elem()
		}
		for value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		valueContainer := result.MapIndex(key)
		if !valueContainer.IsValid() {
			valueContainer = reflect.MakeMap(containerType)
			result.SetMapIndex(key, valueContainer)
		}
		valueContainer.SetMapIndex(value, trueValue)
	}
	return
}

func ptrMapDeepEqual(v1, v2 reflect.Value) bool {
	regMap1 := convertPtrMap(v1)
	regMap2 := convertPtrMap(v2)
	if regMap1.Len() != regMap2.Len() {
		fmt.Println("not same len")
		return false
	}
	return reflect.DeepEqual(regMap1.Interface(), regMap2.Interface())
}

func DeepEqual(i1, i2 interface{}) bool {
	v1 := reflect.ValueOf(i1)
	if v1.Kind() == reflect.Map && (v1.Type().Key().Kind() == reflect.Ptr || v1.Type().Elem().Kind() == reflect.Ptr) {
		return ptrMapDeepEqual(v1, reflect.ValueOf(i2))
	}
	return reflect.DeepEqual(i1, i2)
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
	dstElem := reflect.ValueOf(dst).Elem().Interface()
	toCmp := cmp
	if cmpVal := reflect.ValueOf(toCmp); cmpVal.Kind() == reflect.Ptr {
		toCmp = cmpVal.Elem().Interface()
	}
	if !DeepEqual(toCmp, dstElem) {
		t.Fatalf("Encoding/decoding %#v produced %#v", toCmp, dstElem)
	}
	b, err := Marshal(src)
	if err != nil {
		t.Fatalf("Unable to marshal %v: %v", src, err)
	}
	if err = Unmarshal(b, dst); err != nil {
		t.Fatalf("Unable to unmarshal to %v: %v", dstElem, err)
	}
	dstElem = reflect.ValueOf(dst).Elem().Interface()
	if !DeepEqual(toCmp, dstElem) {
		t.Fatalf("Marshalling/unmarshalling %v produced %v", toCmp, dstElem)
	}
}

type nestedThing1 map[string][]int
type nestedThing2 map[interface{}]map[int]interface{}

func TestEncodeDecodeNestedTypedThingsToInterfaces(t *testing.T) {
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

func TestEncodeDecodeNestedInterfaces(t *testing.T) {
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
