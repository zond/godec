package godec

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	binc "github.com/ugorji/go/codec"
	"github.com/zond/godec/primitives"
)

func encodeDecode(t *testing.T, src, dst interface{}) {
	b, err := Marshal(src)
	if err != nil {
		t.Fatalf("Unable to encode %v: %v", src, err)
		return
	}
	if err = Unmarshal(b, dst); err != nil {
		t.Fatalf("Unable to decode to %v: %v", reflect.ValueOf(dst).Elem(), err)
	}
	dstElem := reflect.ValueOf(dst).Elem().Interface()
	if !reflect.DeepEqual(src, dstElem) {
		t.Fatalf("Encoding/decoding %v produced %v", src, dstElem)
	}
}

type marshaller interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

var bigSlice = []string{}
var bigMap = map[string]string{}

func init() {
	for i := 0; i < 1000; i++ {
		s := fmt.Sprintf("String nr %v", i)
		bigSlice = append(bigSlice, s)
		bigMap[s] = s
	}
}

const (
	encode = 1 << iota
	decode
)

func runBenchMap(b *testing.B, m marshaller) {
	target := map[string]string{}
	for i := 0; i < b.N; i++ {
		res, err := m.Marshal(bigMap)
		if err != nil {
			b.Fatalf("%v", err)
		}
		if len(res) == 0 {
			b.Fatalf("Zero marshalling?")
		}
		if err = m.Unmarshal(res, &target); err != nil {
			b.Fatalf("%v", err)
		}
		b.StopTimer()
		if len(bigMap) != len(target) {
			b.Fatalf("Incorrect unmarshal")
		}
		for k, v := range bigMap {
			if v != target[k] {
				b.Fatalf("Incorrect unmarshal")
			}
		}
		b.StartTimer()
	}
}

func runBenchSlice(b *testing.B, m marshaller, ops int) {
	var encoded []byte
	var err error
	if ops&encode == 0 {
		b.StopTimer()
		if encoded, err = m.Marshal(bigSlice); err != nil {
			b.Fatalf("%v", err)
		}
		if len(encoded) == 0 {
			b.Fatalf("Zero marshalling?")
		}
		b.StartTimer()
	}
	target := []string{}
	for i := 0; i < b.N; i++ {
		if ops&encode == encode {
			if encoded, err = m.Marshal(bigSlice); err != nil {
				b.Fatalf("%v", err)
			}
			if len(encoded) == 0 {
				b.Fatalf("Zero marshalling?")
			}
		}
		if ops&decode == decode {
			if err = m.Unmarshal(encoded, &target); err != nil {
				b.Fatalf("%v", err)
			}
			b.StopTimer()
			if len(target) != len(bigSlice) {
				b.Fatalf("Incorrect unmarshal")
			}
			for i, v := range target {
				if bigSlice[i] != v {
					b.Fatalf("Incorrect unmarshal")
				}
			}
			b.StartTimer()
		}
	}
}

type jsonMarshaller struct{}

func (self jsonMarshaller) Marshal(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (self jsonMarshaller) Unmarshal(b []byte, i interface{}) error {
	return json.Unmarshal(b, i)
}

var bh binc.BincHandle

type bincMarshaller struct{}

func (self bincMarshaller) Marshal(i interface{}) (b []byte, err error) {
	if err = binc.NewEncoderBytes(&b, &bh).Encode(i); err != nil {
		return
	}
	return
}

func (self bincMarshaller) Unmarshal(b []byte, i interface{}) (err error) {
	if err = binc.NewDecoderBytes(b, &bh).Decode(i); err != nil {
		return
	}
	return
}

type godecMarshaller struct{}

func (self godecMarshaller) Marshal(i interface{}) (b []byte, err error) {
	return Marshal(i)
}

func (self godecMarshaller) Unmarshal(b []byte, i interface{}) (err error) {
	return Unmarshal(b, i)
}

type gobMarshaller struct{}

func (self gobMarshaller) Marshal(i interface{}) (b []byte, err error) {
	buf := &bytes.Buffer{}
	if err = gob.NewEncoder(buf).Encode(i); err != nil {
		return
	}
	b = buf.Bytes()
	return
}

func (self gobMarshaller) Unmarshal(b []byte, i interface{}) (err error) {
	if err = gob.NewDecoder(bytes.NewBuffer(b)).Decode(i); err != nil {
		return
	}
	return
}

func BenchmarkJSONStringMap(b *testing.B) {
	runBenchMap(b, jsonMarshaller{})
}

func BenchmarkBincStringMap(b *testing.B) {
	runBenchMap(b, bincMarshaller{})
}

func BenchmarkGobStringMap(b *testing.B) {
	runBenchMap(b, gobMarshaller{})
}

func BenchmarkGodecStringMap(b *testing.B) {
	runBenchMap(b, godecMarshaller{})
}

func BenchmarkJSONStringSlice(b *testing.B) {
	runBenchSlice(b, jsonMarshaller{}, encode|decode)
}

func BenchmarkBincStringSlice(b *testing.B) {
	runBenchSlice(b, bincMarshaller{}, encode|decode)
}

func BenchmarkGobStringSlice(b *testing.B) {
	runBenchSlice(b, gobMarshaller{}, encode|decode)
}

func BenchmarkGodecStringSlice(b *testing.B) {
	runBenchSlice(b, godecMarshaller{}, encode|decode)
}

func BenchmarkGodecStringSliceEncode(b *testing.B) {
	runBenchSlice(b, godecMarshaller{}, encode)
}

func BenchmarkBincStringSliceEncode(b *testing.B) {
	runBenchSlice(b, bincMarshaller{}, encode)
}

func BenchmarkGodecStringSliceDecode(b *testing.B) {
	runBenchSlice(b, godecMarshaller{}, decode)
}

func BenchmarkBincStringSliceDecode(b *testing.B) {
	runBenchSlice(b, bincMarshaller{}, decode)
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
		buf := &bytes.Buffer{}
		i1 := uint64(rand.Int63())
		if err := primitives.Rawencodeuint64(&Encoder{Writer: buf}, i1); err != nil {
			t.Fatalf("%v", err)
		}
		var i2 uint64
		by := buf.Bytes()
		if err := primitives.Rawdecodeuint64(&Decoder{DecodeReader: IODecodeReader{buf}}, &i2); err != nil {
			t.Fatalf("%v", err)
		}
		if i1 != i2 {
			t.Fatalf("Encoded %v to %v, and decoded that to %v", i1, by, i2)
		}
	}
}

func TestEncodeDecodeInt64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		buf := &bytes.Buffer{}
		i1 := rand.Int63()
		if err := primitives.Rawencodeint64(&Encoder{buf}, i1); err != nil {
			t.Fatalf("%v", err)
		}
		var i2 int64
		by := buf.Bytes()
		if err := primitives.Rawdecodeint64(&Decoder{DecodeReader: IODecodeReader{buf}}, &i2); err != nil {
			t.Fatalf("%v", err)
		}
		if i1 != i2 {
			t.Fatalf("Encoded %v to %v, and decoded that to %v", i1, by, i2)
		}
	}
}
