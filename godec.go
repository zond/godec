package godec

import "io"

type bytesEncodeWriter struct {
	buf []byte
}

func (self *bytesEncodeWriter) WriteBytes(b []byte) (err error) {
	l := len(b)
	self.grow(l)
	copy(self.buf[len(self.buf)-l:], b)
	return
}

func (self *bytesEncodeWriter) WriteUint64(u uint64) (err error) {
	l := len(self.buf)
	self.grow(8)
	self.buf[l] = byte(u >> 56)
	self.buf[l+1] = byte(u >> 48)
	self.buf[l+2] = byte(u >> 40)
	self.buf[l+3] = byte(u >> 32)
	self.buf[l+4] = byte(u >> 24)
	self.buf[l+5] = byte(u >> 16)
	self.buf[l+6] = byte(u >> 8)
	self.buf[l+7] = byte(u)
	return
}

func (self *bytesEncodeWriter) grow(n int) {
	if n > cap(self.buf)-len(self.buf) {
		newBuf := make([]byte, len(self.buf), 2*cap(self.buf)+n)
		copy(newBuf, self.buf)
		self.buf = newBuf
	}
	self.buf = self.buf[:len(self.buf)+n]
}

func (self *bytesEncodeWriter) WriteString(s string) (err error) {
	l := len(s)
	self.grow(l)
	copy(self.buf[len(self.buf)-l:], s)
	return
}

func Marshal(i interface{}) (result []byte, err error) {
	w := &bytesEncodeWriter{
		buf: make([]byte, 0, 1<<6),
	}
	enc := &Encoder{
		EncodeWriter: w,
	}
	if err = enc.Encode(i); err != nil {
		return
	}
	result = w.buf
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
	if self.pos+8 > len(self.buf) {
		err = io.EOF
		return
	}
	result = uint64(self.buf[self.pos])<<56 + uint64(self.buf[self.pos+1])<<48 + uint64(self.buf[self.pos+2])<<40 + uint64(self.buf[self.pos+3])<<32 + uint64(self.buf[self.pos+4])<<24 + uint64(self.buf[self.pos+5])<<16 + uint64(self.buf[self.pos+6])<<8 + uint64(self.buf[self.pos+7])
	self.pos += 8
	return
}

func Unmarshal(b []byte, i interface{}) (err error) {
	dec := &Decoder{
		DecodeReader: &bytesDecodeReader{buf: b},
	}
	if err = dec.Decode(i); err != nil {
		return
	}
	return
}
