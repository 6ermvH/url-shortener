package base63

const StdAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

var StdEncoding = NewEncoding(StdAlphabet, 10)

type Encoding struct {
	alphabet string
	base     uint64
	length   int
}

func NewEncoding(alphabet string, length int) *Encoding {
	return &Encoding{
		alphabet: alphabet,
		base:     uint64(len(alphabet)),
		length:   length,
	}
}

func (e *Encoding) Encode(num uint64) string {
	result := make([]byte, e.length)
	for i := range e.length {
		result[i] = e.alphabet[num%e.base]
		num /= e.base
	}
	return string(result)
}
