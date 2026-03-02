package base63

const (
	base = uint64(63)
)

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

func (e *Encoding) Encode(num uint64) string {
	result := make([]byte, e.length)
	for i := range e.length {
		result[i] = e.alphabet[num%base]
		num /= base
	}

	return string(result)
}
