package godec

import (
	"encoding/binary"
	"fmt"
	"io"
)

type bytesEncodeWriter struct {
	buf []byte
	pos int
	tmp []byte
}

func (self *bytesEncodeWriter) WriteBytes(b []byte) (err error) {
	l := len(b)
	self.grow(l)
	copy(self.buf[self.pos-l:], b)
	return
}

func (self *bytesEncodeWriter) WriteUint64(u uint64) (err error) {
	wrote := binary.PutUvarint(self.tmp, u)
	self.grow(wrote)
	copy(self.buf[self.pos-wrote:], self.tmp[:wrote])
	return
}

func (self *bytesEncodeWriter) grow(n int) {
	if n+self.pos > len(self.buf) {
		newBuf := make([]byte, len(self.buf)*2+n)
		copy(newBuf, self.buf)
		self.buf = newBuf
	}
	self.pos += n
}

func (self *bytesEncodeWriter) WriteString(s string) (err error) {
	l := len(s)
	self.grow(l)
	copy(self.buf[self.pos-l:], s)
	return
}

func Marshal(i interface{}) (result []byte, err error) {
	w := &bytesEncodeWriter{
		buf: make([]byte, 1<<6),
		tmp: make([]byte, binary.MaxVarintLen64),
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

type bytesDecodeReader struct {
	buf []byte
	pos int
}

func (self *bytesDecodeReader) ReadByte() (result byte, err error) {
	if self.pos >= len(self.buf) {
		err = io.EOF
		return
	}
	result = self.buf[self.pos]
	self.pos++
	return
}

func (self *bytesDecodeReader) ReadBytes(n int) (result []byte, err error) {
	if self.pos+n > len(self.buf) {
		err = io.EOF
		return
	}
	result = self.buf[self.pos : self.pos+n]
	self.pos += n
	return
}

func (self *bytesDecodeReader) ReadUint64() (result uint64, err error) {
	result, read := binary.Uvarint(self.buf[self.pos:])
	if read <= 0 {
		err = fmt.Errorf("Unable to read uint: %v")
		return
	}
	self.pos += read
	return
}

func Unmarshal(b []byte, i interface{}) (err error) {
	dec := &Decoder{
		DecodeReader: &bytesDecodeReader{
			buf: b,
		},
	}
	if err = dec.Decode(i); err != nil {
		return
	}
	return
}
