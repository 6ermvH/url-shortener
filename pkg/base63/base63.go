package base63

import "encoding/binary"

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	base     = uint64(len(alphabet))
)

func Encode(b [8]byte, length int) string {
	num := binary.BigEndian.Uint64(b[:])

	result := make([]byte, length)
	for i := range length {
		result[i] = alphabet[num%base]
		num /= base
	}

	return string(result)
}
