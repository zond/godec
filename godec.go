package godec

import "bytes"

func Marshal(i interface{}) (result []byte, err error) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)
	if err = enc.Encode(i); err != nil {
		return
	}
	result = buf.Bytes()
	return
}

func Unmarshal(b []byte, i interface{}) (err error) {
	dec := &Decoder{
		DecodeReader: &BytesDecodeReader{Buf: b},
	}
	if err = dec.Decode(i); err != nil {
		return
	}
	return
}
