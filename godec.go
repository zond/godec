package godec

import (
	"encoding/binary"
	"fmt"
	"io"
)

type uvarintBytesEncodeWriter struct {
	buf []byte
	pos int
}

func (self *uvarintBytesEncodeWriter) WriteBytes(b []byte) (err error) {
	l := len(b)
	self.grow(l)
	copy(self.buf[self.pos:], b)
	self.pos += l
	return
}

func (self *uvarintBytesEncodeWriter) WriteUint64(u uint64) (err error) {
	self.grow(binary.MaxVarintLen64)
	wrote := binary.PutUvarint(self.buf[self.pos:], u)
	self.pos += wrote
	return
}

func (self *uvarintBytesEncodeWriter) grow(n int) {
	if n+self.pos > len(self.buf) {
		newBuf := make([]byte, len(self.buf)*2+n)
		copy(newBuf, self.buf)
		self.buf = newBuf
	}
}

func (self *uvarintBytesEncodeWriter) WriteString(s string) (err error) {
	l := len(s)
	self.grow(l)
	copy(self.buf[self.pos:], s)
	self.pos += l
	return
}

type naiveBytesEncodeWriter struct {
	buf []byte
	pos int
}

func (self *naiveBytesEncodeWriter) WriteBytes(b []byte) (err error) {
	l := len(b)
	self.grow(l)
	copy(self.buf[self.pos:], b)
	self.pos += l
	return
}

func (self *naiveBytesEncodeWriter) WriteUint64(u uint64) (err error) {
	self.grow(8)
	binary.BigEndian.PutUint64(self.buf[self.pos:], u)
	self.pos += 8
	return
}

func (self *naiveBytesEncodeWriter) grow(n int) {
	if n+self.pos > len(self.buf) {
		newBuf := make([]byte, len(self.buf)*2+n)
		copy(newBuf, self.buf)
		self.buf = newBuf
	}
}

func (self *naiveBytesEncodeWriter) WriteString(s string) (err error) {
	l := len(s)
	self.grow(l)
	copy(self.buf[self.pos:], s)
	self.pos += l
	return
}

func Marshal(i interface{}) (result []byte, err error) {
	w := &naiveBytesEncodeWriter{
		buf: make([]byte, 1<<6),
	}
	enc := &Encoder{
		EncodeWriter: w,
	}
	if err = enc.Encode(i); err != nil {
		return
	}
	result = w.buf[:w.pos]
	return
}

type uvarintBytesDecodeReader struct {
	buf []byte
	pos int
}

func (self *uvarintBytesDecodeReader) ReadByte() (result byte, err error) {
	if self.pos >= len(self.buf) {
		err = io.EOF
		return
	}
	result = self.buf[self.pos]
	self.pos++
	return
}

func (self *uvarintBytesDecodeReader) ReadBytes(n int) (result []byte, err error) {
	if self.pos+n > len(self.buf) {
		err = io.EOF
		return
	}
	result = self.buf[self.pos : self.pos+n]
	self.pos += n
	return
}

func (self *uvarintBytesDecodeReader) ReadUint64() (result uint64, err error) {
	result, read := binary.Uvarint(self.buf[self.pos:])
	if read <= 0 {
		err = fmt.Errorf("Unable to read uint: %v")
		return
	}
	self.pos += read
	return
}

type naiveBytesDecodeReader struct {
	buf []byte
	pos int
}

func (self *naiveBytesDecodeReader) ReadByte() (result byte, err error) {
	if self.pos >= len(self.buf) {
		err = io.EOF
		return
	}
	result = self.buf[self.pos]
	self.pos++
	return
}

func (self *naiveBytesDecodeReader) ReadBytes(n int) (result []byte, err error) {
	if self.pos+n > len(self.buf) {
		err = io.EOF
		return
	}
	result = self.buf[self.pos : self.pos+n]
	self.pos += n
	return
}

func (self *naiveBytesDecodeReader) ReadUint64() (result uint64, err error) {
	if self.pos+8 > len(self.buf) {
		err = io.EOF
		return
	}
	result = binary.BigEndian.Uint64(self.buf[self.pos:])
	self.pos += 8
	return
}

func Unmarshal(b []byte, i interface{}) (err error) {
	dec := &Decoder{
		DecodeReader: &naiveBytesDecodeReader{
			buf: b,
		},
	}
	if err = dec.Decode(i); err != nil {
		return
	}
	return
}
