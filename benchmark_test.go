package godec

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	binc "github.com/ugorji/go/codec"
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

/*
I would love to run this test, but binc seems broken?
func BenchmarkBincStringMap(b *testing.B) {
	runBenchMap(b, bincMarshaller{})
}
*/

func BenchmarkGobStringMap(b *testing.B) {
	runBenchMap(b, gobMarshaller{})
}

func BenchmarkGodecStringMap(b *testing.B) {
	runBenchMap(b, godecMarshaller{})
}

func BenchmarkJSONStringSlice(b *testing.B) {
	runBenchStringSlice(b, jsonMarshaller{}, doEncode|doDecode)
}

func BenchmarkBincStringSlice(b *testing.B) {
	runBenchStringSlice(b, bincMarshaller{}, doEncode|doDecode)
}

func BenchmarkGobStringSlice(b *testing.B) {
	runBenchStringSlice(b, gobMarshaller{}, doEncode|doDecode)
}

func BenchmarkGodecStringSlice(b *testing.B) {
	runBenchStringSlice(b, godecMarshaller{}, doEncode|doDecode)
}

func BenchmarkGodecUint8Encode(b *testing.B) {
	runBenchUint8Encode(b, godecMarshaller{})
}

func BenchmarkBincUint8Encode(b *testing.B) {
	runBenchUint8Encode(b, bincMarshaller{})
}

func BenchmarkGodecUint64Encode(b *testing.B) {
	runBenchUint64Encode(b, godecMarshaller{})
}

func BenchmarkBincUint64Encode(b *testing.B) {
	runBenchUint64Encode(b, bincMarshaller{})
}

func BenchmarkGodecStringEncode(b *testing.B) {
	runBenchStringEncode(b, godecMarshaller{})
}

func BenchmarkBincStringEncode(b *testing.B) {
	runBenchStringEncode(b, bincMarshaller{})
}

func BenchmarkGodecIntSliceEncode(b *testing.B) {
	runBenchIntSlice(b, godecMarshaller{}, doEncode)
}

func BenchmarkBincIntSliceEncode(b *testing.B) {
	runBenchIntSlice(b, bincMarshaller{}, doEncode)
}

func BenchmarkGodecStringSliceEncode(b *testing.B) {
	runBenchStringSlice(b, godecMarshaller{}, doEncode)
}

func BenchmarkBincStringSliceEncode(b *testing.B) {
	runBenchStringSlice(b, bincMarshaller{}, doEncode)
}

func BenchmarkGodecStringSliceDecode(b *testing.B) {
	runBenchStringSlice(b, godecMarshaller{}, doDecode)
}

func BenchmarkBincStringSliceDecode(b *testing.B) {
	runBenchStringSlice(b, bincMarshaller{}, doDecode)
}
