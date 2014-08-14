package godec

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"labix.org/v2/mgo/bson"

	binc "github.com/ugorji/go/codec"
	"github.com/vmihailenco/msgpack"
)

type marshaller interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

var bigIntSlice = []int64{}
var bigStringSlice = []string{}
var bigMap = map[string]string{}
var randomNumbers []uint64

func init() {
	rand.Seed(0)
	for i := 0; i < 10000; i++ {
		s := fmt.Sprintf("String nr %v", i)
		bigStringSlice = append(bigStringSlice, s)
		bigIntSlice = append(bigIntSlice, rand.Int63())
		bigMap[s] = s
		randomNumbers = append(randomNumbers, uint64(rand.Int63()))
	}
}

const (
	doEncode = 1 << iota
	doDecode
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

type BenchStructT1 struct {
	A int
	B string
	C []byte
	D time.Time
	E []string
	F byte
}

var benchStruct1 = BenchStructT1{
	A: 4412,
	B: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	C: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
	D: time.Now(),
	E: []string{"aaaaaaaaaaa", "bbbbbbbbbbbb", "ccccccccccccc", "dddddddddddddd", "eeeeeeeeeeeeee", "fffffffffffffff"},
	F: 43,
}

type benchStructT2 struct {
	BenchStructT1
	G BenchStructT1
	H *benchStructT2
	I *[]int
	J map[string]*map[int]int
	K map[int]*benchStructT2
}

var smallIntSlice = []int{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
}

var smallIntMap = map[int]int{
	1:  4,
	4:  5,
	7:  1,
	8:  11,
	91: 32,
}

var benchStruct3 = benchStructT2{
	BenchStructT1: benchStruct1,
	G:             benchStruct1,
	I:             &smallIntSlice,
	J: map[string]*map[int]int{
		"a": &smallIntMap,
		"b": &smallIntMap,
	},
}

var benchStruct2 = benchStructT2{
	BenchStructT1: benchStruct1,
	G:             benchStruct1,
	H:             &benchStruct3,
	I:             &smallIntSlice,
	J: map[string]*map[int]int{
		"a": &smallIntMap,
		"b": &smallIntMap,
	},
	K: map[int]*benchStructT2{
		44: &benchStruct3,
		91: &benchStruct3,
	},
}

func runBenchNestedStruct(b *testing.B, m marshaller, ops int) {
	var encoded []byte
	var err error
	if ops&doEncode == 0 {
		b.StopTimer()
		if encoded, err = m.Marshal(benchStruct1); err != nil {
			b.Fatal("%v", err)
		}
		if len(encoded) == 0 {
			b.Fatalf("Zero marshalling")
		}
		b.StartTimer()
	}
	target := benchStructT2{}
	for i := 0; i < b.N; i++ {
		if ops&doEncode == doEncode {
			if encoded, err = m.Marshal(benchStruct2); err != nil {
				b.Fatalf("%v", err)
			}
			if len(encoded) == 0 {
				b.Fatalf("Zero marshalling?")
			}
		}
		if ops&doDecode == doDecode {
			if err = m.Unmarshal(encoded, &target); err != nil {
				b.Fatalf("%v", err)
			}
			b.StopTimer()
			if err := DeepEqual(target, benchStruct2); err != nil {
				b.Fatalf("Incorrect unmarshal: %+v != %v: %v", target, benchStruct2, err)
			}
			b.StartTimer()
		}
	}
}

func runBenchFlatStruct(b *testing.B, m marshaller, ops int) {
	var encoded []byte
	var err error
	if ops&doEncode == 0 {
		b.StopTimer()
		if encoded, err = m.Marshal(benchStruct1); err != nil {
			b.Fatal("%v", err)
		}
		if len(encoded) == 0 {
			b.Fatalf("Zero marshalling")
		}
		b.StartTimer()
	}
	target := BenchStructT1{}
	for i := 0; i < b.N; i++ {
		if ops&doEncode == doEncode {
			if encoded, err = m.Marshal(benchStruct1); err != nil {
				b.Fatalf("%v", err)
			}
			if len(encoded) == 0 {
				b.Fatalf("Zero marshalling?")
			}
		}
		if ops&doDecode == doDecode {
			if err = m.Unmarshal(encoded, &target); err != nil {
				b.Fatalf("%v", err)
			}
			b.StopTimer()
			if err := DeepEqual(target, benchStruct1); err != nil {
				b.Fatalf("Incorrect unmarshal: %v", err)
			}
			b.StartTimer()
		}
	}
}

func runBenchIntSlice(b *testing.B, m marshaller, ops int) {
	var encoded []byte
	var err error
	if ops&doEncode == 0 {
		b.StopTimer()
		if encoded, err = m.Marshal(bigIntSlice); err != nil {
			b.Fatalf("%v", err)
		}
		if len(encoded) == 0 {
			b.Fatalf("Zero marshalling?")
		}
		b.StartTimer()
	}
	target := []int64{}
	for i := 0; i < b.N; i++ {
		if ops&doEncode == doEncode {
			if encoded, err = m.Marshal(bigIntSlice); err != nil {
				b.Fatalf("%v", err)
			}
			if len(encoded) == 0 {
				b.Fatalf("Zero marshalling?")
			}
		}
		if ops&doDecode == doDecode {
			if err = m.Unmarshal(encoded, &target); err != nil {
				b.Fatalf("%v", err)
			}
			b.StopTimer()
			if len(target) != len(bigIntSlice) {
				b.Fatalf("Incorrect unmarshal")
			}
			for i, v := range target {
				if bigIntSlice[i] != v {
					b.Fatalf("Incorrect unmarshal")
				}
			}
			b.StartTimer()
		}
	}
}

func runBenchStringSlice(b *testing.B, m marshaller, ops int) {
	var encoded []byte
	var err error
	if ops&doEncode == 0 {
		b.StopTimer()
		if encoded, err = m.Marshal(bigStringSlice); err != nil {
			b.Fatalf("%v", err)
		}
		if len(encoded) == 0 {
			b.Fatalf("Zero marshalling?")
		}
		b.StartTimer()
	}
	target := []string{}
	for i := 0; i < b.N; i++ {
		if ops&doEncode == doEncode {
			if encoded, err = m.Marshal(bigStringSlice); err != nil {
				b.Fatalf("%v", err)
			}
			if len(encoded) == 0 {
				b.Fatalf("Zero marshalling?")
			}
		}
		if ops&doDecode == doDecode {
			if err = m.Unmarshal(encoded, &target); err != nil {
				b.Fatalf("%v", err)
			}
			b.StopTimer()
			if len(target) != len(bigStringSlice) {
				b.Fatalf("Incorrect unmarshal")
			}
			for i, v := range target {
				if bigStringSlice[i] != v {
					b.Fatalf("Incorrect unmarshal")
				}
			}
			b.StartTimer()
		}
	}
}

func runBenchUint8Encode(b *testing.B, m marshaller) {
	var by []byte
	var err error
	j := 0
	for i := 0; i < b.N; i++ {
		if by, err = m.Marshal(uint8(randomNumbers[j])); err != nil {
			b.Fatalf("%v, %p", err, &by)
		}
		j += 1
		if j >= len(randomNumbers) {
			j = 0
		}
	}
}

func runBenchUint64Encode(b *testing.B, m marshaller) {
	var by []byte
	var err error
	j := 0
	for i := 0; i < b.N; i++ {
		if by, err = m.Marshal(randomNumbers[j]); err != nil {
			b.Fatalf("%v, %p", err, &by)
		}
		j += 1
		if j >= len(randomNumbers) {
			j = 0
		}
	}
}

func runBenchStringEncode(b *testing.B, m marshaller) {
	s := "a string of mediocre length but at least a few characters"
	var by []byte
	var err error
	for i := 0; i < b.N; i++ {
		if by, err = m.Marshal(s); err != nil {
			b.Fatalf("%v, %p", err, &by)
		}
	}
}

type bsonMarshaller struct{}

func (self bsonMarshaller) Marshal(i interface{}) ([]byte, error) {
	return bson.Marshal(i)
}

func (self bsonMarshaller) Unmarshal(b []byte, i interface{}) error {
	return bson.Unmarshal(b, i)
}

type msgpackMarshaller struct{}

func (self msgpackMarshaller) Marshal(i interface{}) ([]byte, error) {
	return msgpack.Marshal(i)
}

func (self msgpackMarshaller) Unmarshal(b []byte, i interface{}) error {
	return msgpack.Unmarshal(b, i)
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

func BenchmarkGobNestedStruct(b *testing.B) {
	runBenchNestedStruct(b, gobMarshaller{}, doEncode|doDecode)
}

func BenchmarkMsgpackNestedStruct(b *testing.B) {
	runBenchNestedStruct(b, msgpackMarshaller{}, doEncode|doDecode)
}

/*
func BenchmarkBincNestedStruct(b *testing.B) {
	runBenchNestedStruct(b, bincMarshaller{}, doEncode|doDecode)
}
*/

func BenchmarkGodecNestedStruct(b *testing.B) {
	runBenchNestedStruct(b, godecMarshaller{}, doEncode|doDecode)
}

func BenchmarkJSONFlatStruct(b *testing.B) {
	runBenchFlatStruct(b, jsonMarshaller{}, doEncode|doDecode)
}

func BenchmarkGobFlatStruct(b *testing.B) {
	runBenchFlatStruct(b, gobMarshaller{}, doEncode|doDecode)
}

func BenchmarkBincFlatStruct(b *testing.B) {
	runBenchFlatStruct(b, bincMarshaller{}, doEncode|doDecode)
}

func BenchmarkMsgpackFlatStruct(b *testing.B) {
	runBenchFlatStruct(b, msgpackMarshaller{}, doEncode|doDecode)
}

func BenchmarkGodecFlatStruct(b *testing.B) {
	runBenchFlatStruct(b, godecMarshaller{}, doEncode|doDecode)
}

func BenchmarkJSONStringSlice(b *testing.B) {
	runBenchStringSlice(b, jsonMarshaller{}, doEncode|doDecode)
}

func BenchmarkGobStringSlice(b *testing.B) {
	runBenchStringSlice(b, gobMarshaller{}, doEncode|doDecode)
}

func BenchmarkBincStringSlice(b *testing.B) {
	runBenchStringSlice(b, bincMarshaller{}, doEncode|doDecode)
}

func BenchmarkGodecStringSlice(b *testing.B) {
	runBenchStringSlice(b, godecMarshaller{}, doEncode|doDecode)
}

func BenchmarkBincUint8Encode(b *testing.B) {
	runBenchUint8Encode(b, bincMarshaller{})
}

func BenchmarkGodecUint8Encode(b *testing.B) {
	runBenchUint8Encode(b, godecMarshaller{})
}

func BenchmarkBincUint64Encode(b *testing.B) {
	runBenchUint64Encode(b, bincMarshaller{})
}

func BenchmarkGodecUint64Encode(b *testing.B) {
	runBenchUint64Encode(b, godecMarshaller{})
}

func BenchmarkBincStringEncode(b *testing.B) {
	runBenchStringEncode(b, bincMarshaller{})
}

func BenchmarkGodecStringEncode(b *testing.B) {
	runBenchStringEncode(b, godecMarshaller{})
}

func BenchmarkBincIntSliceEncode(b *testing.B) {
	runBenchIntSlice(b, bincMarshaller{}, doEncode)
}

func BenchmarkGodecIntSliceEncode(b *testing.B) {
	runBenchIntSlice(b, godecMarshaller{}, doEncode)
}

func BenchmarkBincStringSliceEncode(b *testing.B) {
	runBenchStringSlice(b, bincMarshaller{}, doEncode)
}

func BenchmarkGodecStringSliceEncode(b *testing.B) {
	runBenchStringSlice(b, godecMarshaller{}, doEncode)
}

func BenchmarkBincStringSliceDecode(b *testing.B) {
	runBenchStringSlice(b, bincMarshaller{}, doDecode)
}

func BenchmarkGodecStringSliceDecode(b *testing.B) {
	runBenchStringSlice(b, godecMarshaller{}, doDecode)
}
