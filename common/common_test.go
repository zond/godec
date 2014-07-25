package common

import (
	"encoding/binary"
	"math"
	"math/rand"
	"testing"
	"time"
)

var biguints = []uint64{}
var smalluints = []uint64{}

func init() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		biguints = append(biguints, uint64(rand.Int63()))
		smalluints = append(smalluints, uint64(rand.Int63n(1024)))
	}
}

func BenchmarkStdBigLog2(b *testing.B) {
	j := 0
	for i := 0; i < b.N; i++ {
		math.Ilogb(float64(biguints[j]))
		j += 1
		if j >= len(biguints) {
			j = 0
		}
	}
}

func BenchmarkGodecBigLog2(b *testing.B) {
	j := 0
	for i := 0; i < b.N; i++ {
		Log2_64(biguints[j])
		j += 1
		if j >= len(biguints) {
			j = 0
		}
	}
}

func BenchmarkStdSmallLog2(b *testing.B) {
	j := 0
	for i := 0; i < b.N; i++ {
		math.Ilogb(float64(smalluints[j]))
		j += 1
		if j >= len(biguints) {
			j = 0
		}
	}
}

func sizeCheck(u uint64) uint8 {
	if u < 1<<7 {
		return 1
	}
	if u < 1<<14 {
		return 2
	}
	if u < 1<<21 {
		return 3
	}
	if u < 1<<28 {
		return 4
	}
	if u < 1<<35 {
		return 5
	}
	if u < 1<<42 {
		return 6
	}
	if u < 1<<49 {
		return 7
	}
	if u < 1<<56 {
		return 8
	}
	return 9
}

func BenchmarkGodecSmallLog2(b *testing.B) {
	j := 0
	for i := 0; i < b.N; i++ {
		Log2_64(smalluints[j])
		j += 1
		if j >= len(biguints) {
			j = 0
		}
	}
}

func BenchmarkSizeCheckBig(b *testing.B) {
	j := 0
	for i := 0; i < b.N; i++ {
		sizeCheck(biguints[j])
		j += 1
		if j >= len(biguints) {
			j = 0
		}
	}
}

func BenchmarkSizeCheckSmall(b *testing.B) {
	j := 0
	for i := 0; i < b.N; i++ {
		sizeCheck(smalluints[j])
		j += 1
		if j >= len(smalluints) {
			j = 0
		}
	}
}
func TestGodecPutUvarint(t *testing.T) {
	buf := make([]byte, MaxVarintLen64)
	for _, u := range biguints {
		wrote := PutUvarint(buf, u)
		u2 := ReadUvarint(buf[:wrote])
		if u != u2 {
			t.Errorf("Wanted %v was encoded to %v which was decoded to %v", u, buf[:wrote], u2)
		}
	}
	for _, u := range smalluints {
		wrote := PutUvarint(buf, u)
		u2 := ReadUvarint(buf[:wrote])
		if u != u2 {
			t.Errorf("Wanted %v was encoded to %v which was decoded to %v", u, buf[:wrote], u2)
		}
	}
}

func BenchmarkGodecPutUvarintBig(b *testing.B) {
	j := 0
	buf := make([]byte, MaxVarintLen64)
	for i := 0; i < b.N; i++ {
		PutUvarint(buf, biguints[j])
		j += 1
		if j >= len(biguints) {
			j = 0
		}
	}
}

func BenchmarkBinaryPutUvarintBig(b *testing.B) {
	j := 0
	buf := make([]byte, MaxVarintLen64)
	for i := 0; i < b.N; i++ {
		binary.PutUvarint(buf, biguints[j])
		j += 1
		if j >= len(biguints) {
			j = 0
		}
	}
}

func BenchmarkGodecPutUvarintSmall(b *testing.B) {
	j := 0
	buf := make([]byte, MaxVarintLen64)
	for i := 0; i < b.N; i++ {
		PutUvarint(buf, smalluints[j])
		j += 1
		if j >= len(smalluints) {
			j = 0
		}
	}
}

func BenchmarkBinaryPutUvarintSmall(b *testing.B) {
	j := 0
	buf := make([]byte, MaxVarintLen64)
	for i := 0; i < b.N; i++ {
		binary.PutUvarint(buf, smalluints[j])
		j += 1
		if j >= len(smalluints) {
			j = 0
		}
	}
}
