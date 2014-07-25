package common

import "math"

var tab64 = [64]uint8{
	63, 0, 58, 1, 59, 47, 53, 2,
	60, 39, 48, 27, 54, 33, 42, 3,
	61, 51, 37, 40, 49, 18, 28, 20,
	55, 30, 34, 11, 43, 14, 22, 4,
	62, 57, 46, 52, 38, 26, 32, 41,
	50, 36, 17, 19, 29, 10, 13, 21,
	56, 45, 25, 31, 35, 16, 9, 12,
	44, 24, 15, 8, 23, 7, 6, 5,
}

const (
	MaxVarintLen64 = 9
)

var tab8 = [math.MaxUint8 + 1]uint8{}

func init() {
	for i := uint64(0); i < math.MaxUint8+1; i++ {
		tab8[i] = log2_64(i)
	}
}

func log2_8(value uint8) uint8 {
	return tab8[int(value)]
}

func Log2_64(value uint64) uint8 {
	return log2_64(value)
}

func log2_64(value uint64) uint8 {
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	value |= value >> 32
	return tab64[((uint64)((value-(value>>1))*0x07EDD5E59A4E28C2))>>58]
}

func ReadUvarint(b []byte) (u uint64) {
	numBytes := uint8(9)
	if b[0] != 0 {
		numBytes = 7 - log2_8(b[0]) + 1
	}
	switch numBytes {
	case 1:
		u = uint64(b[0])
		u &= ^(uint64(1) << 7)
	case 2:
		u = uint64(b[0])<<8 + uint64(b[1])
		u &= ^(uint64(1) << 14)
	case 3:
		u = uint64(b[0])<<16 + uint64(b[1])<<8 + uint64(b[2])
		u &= ^(uint64(1) << 21)
	case 4:
		u = uint64(b[0])<<24 + uint64(b[1])<<16 + uint64(b[2])<<8 + uint64(b[3])
		u &= ^(uint64(1) << 28)
	case 5:
		u = uint64(b[0])<<32 + uint64(b[1])<<24 + uint64(b[2])<<16 + uint64(b[3])<<8 + uint64(b[4])
		u &= ^(uint64(1) << 35)
	case 6:
		u = uint64(b[0])<<40 + uint64(b[1])<<32 + uint64(b[2])<<24 + uint64(b[3])<<16 + uint64(b[4])<<8 + uint64(b[5])
		u &= ^(uint64(1) << 42)
	case 7:
		u = uint64(b[0])<<48 + uint64(b[1])<<40 + uint64(b[2])<<32 + uint64(b[3])<<24 + uint64(b[4])<<16 + uint64(b[5])<<8 + uint64(b[6])
		u &= ^(uint64(1) << 49)
	case 8:
		u = uint64(b[1])<<48 + uint64(b[2])<<40 + uint64(b[3])<<32 + uint64(b[4])<<24 + uint64(b[5])<<16 + uint64(b[6])<<8 + uint64(b[7])
	case 9:
		u = uint64(b[1])<<56 + uint64(b[2])<<48 + uint64(b[3])<<40 + uint64(b[4])<<32 + uint64(b[5])<<24 + uint64(b[6])<<16 + uint64(b[7])<<8 + uint64(b[8])
	}
	return
}

func PutUvarint(b []byte, u uint64) (numBytes uint8) {
	numBytes = 1
	if u > 0 {
		numBytes = log2_64(u)/7 + 1
	}
	switch numBytes {
	case 1:
		u |= uint64(1) << 7
		b[0] = byte(u)
	case 2:
		u |= uint64(1) << 14
		b[1] = byte(u)
		u >>= 8
		b[0] = byte(u)
	case 3:
		u |= uint64(1) << 21
		b[2] = byte(u)
		u >>= 8
		b[1] = byte(u)
		u >>= 8
		b[0] = byte(u)
	case 4:
		u |= uint64(1) << 28
		b[3] = byte(u)
		u >>= 8
		b[2] = byte(u)
		u >>= 8
		b[1] = byte(u)
		u >>= 8
		b[0] = byte(u)
	case 5:
		u |= uint64(1) << 35
		b[4] = byte(u)
		u >>= 8
		b[3] = byte(u)
		u >>= 8
		b[2] = byte(u)
		u >>= 8
		b[1] = byte(u)
		u >>= 8
		b[0] = byte(u)
	case 6:
		u |= uint64(1) << 42
		b[5] = byte(u)
		u >>= 8
		b[4] = byte(u)
		u >>= 8
		b[3] = byte(u)
		u >>= 8
		b[2] = byte(u)
		u >>= 8
		b[1] = byte(u)
		u >>= 8
		b[0] = byte(u)
	case 7:
		u |= uint64(1) << 49
		b[6] = byte(u)
		u >>= 8
		b[5] = byte(u)
		u >>= 8
		b[4] = byte(u)
		u >>= 8
		b[3] = byte(u)
		u >>= 8
		b[2] = byte(u)
		u >>= 8
		b[1] = byte(u)
		u >>= 8
		b[0] = byte(u)
	case 8:
		u |= uint64(1) << 56
		b[7] = byte(u)
		u >>= 8
		b[6] = byte(u)
		u >>= 8
		b[5] = byte(u)
		u >>= 8
		b[4] = byte(u)
		u >>= 8
		b[3] = byte(u)
		u >>= 8
		b[2] = byte(u)
		u >>= 8
		b[1] = byte(u)
		u >>= 8
		b[0] = byte(u)
	case 9:
		fallthrough
	case 10:
		numBytes = 9
		b[8] = byte(u)
		u >>= 8
		b[7] = byte(u)
		u >>= 8
		b[6] = byte(u)
		u >>= 8
		b[5] = byte(u)
		u >>= 8
		b[4] = byte(u)
		u >>= 8
		b[3] = byte(u)
		u >>= 8
		b[2] = byte(u)
		u >>= 8
		b[1] = byte(u)
		b[0] = 0
	}
	return
}
