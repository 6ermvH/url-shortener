package base63_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/6ermvH/url-shortener/pkg/base63"
)

const stdAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func TestNewEncoding_PanicsOnInvalidAlphabet(t *testing.T) {
	require.Panics(t, func() {
		base63.NewEncoding("tooshort", 10)
	})
}

func TestEncodeToString_OutputLength(t *testing.T) {
	for _, length := range []int{5, 10, 15} {
		enc := base63.NewEncoding(stdAlphabet, length)
		result := enc.EncodeToString([]byte{0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe})
		require.Len(t, result, length)
	}
}

func TestEncodeToString_OnlyAlphabetChars(t *testing.T) {
	enc := base63.NewEncoding(stdAlphabet, 10)
	result := enc.EncodeToString([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	for _, c := range result {
		require.True(t, strings.ContainsRune(stdAlphabet, c), "character %q is not in the alphabet", c)
	}
}

func TestEncodeToString_Deterministic(t *testing.T) {
	enc := base63.NewEncoding(stdAlphabet, 10)
	src := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	require.Equal(t, enc.EncodeToString(src), enc.EncodeToString(src))
}

func TestEncodeToString_DifferentInputs(t *testing.T) {
	enc := base63.NewEncoding(stdAlphabet, 10)
	a := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	b := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}

	require.NotEqual(t, enc.EncodeToString(a), enc.EncodeToString(b))
}

func TestEncodeToString_Zero(t *testing.T) {
	enc := base63.NewEncoding(stdAlphabet, 10)
	result := enc.EncodeToString([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	expected := strings.Repeat(string(stdAlphabet[0]), 10)
	require.Equal(t, expected, result)
}

func TestEncodedLen(t *testing.T) {
	enc := base63.NewEncoding(stdAlphabet, 10)
	require.Equal(t, 10, enc.EncodedLen())
}
