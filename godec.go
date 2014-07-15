package godec

import "io"

type encodeBytesWriter struct {
	buf []byte
}

func (self *encodeBytesWriter) Write(b []byte) (written int, err error) {
	written = len(b)
	if cap(self.buf)-len(self.buf) < written {
		newBuf := make([]byte, len(self.buf)+written, 2*cap(self.buf)+written)
		copy(newBuf, self.buf)
		self.buf = newBuf
	}
	oldLen := len(self.buf)
	self.buf = self.buf[0 : len(self.buf)+written]
	copy(self.buf[oldLen:], b)
	return
}

func Marshal(i interface{}) (result []byte, err error) {
	w := &encodeBytesWriter{
		buf: make([]byte, 0, 1<<6),
	}
	enc := NewEncoder(w)
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

func Unmarshal(b []byte, i interface{}) (err error) {
	dec := &Decoder{
		DecodeReader: &bytesDecodeReader{buf: b},
	}
	if err = dec.Decode(i); err != nil {
		return
	}
	return
}
