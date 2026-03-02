package base63

import "encoding/binary"

const base = uint64(63)

type Encoding struct {
	alphabet string
	length   int
}

func NewEncoding(alphabet string, length int) *Encoding {
	if len(alphabet) != int(base) {
		panic("base63: alphabet must be exactly 63 characters")
	}
	return &Encoding{
		alphabet: alphabet,
		length:   length,
	}
}

func (e *Encoding) EncodedLen() int {
	return e.length
}

func (e *Encoding) Encode(dst, src []byte) {
	num := binary.BigEndian.Uint64(src)
	for i := range e.length {
		dst[i] = e.alphabet[num%base]
		num /= base
	}
}

func (e *Encoding) EncodeToString(src []byte) string {
	dst := make([]byte, e.length)
	e.Encode(dst, src)
	return string(dst)
}
