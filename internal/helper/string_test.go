package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "nil bytes",
			input: nil,
			want:  "",
		},
		{
			name:  "empty bytes",
			input: []byte{},
			want:  "",
		},
		{
			name:  "ascii bytes",
			input: []byte("hello"),
			want:  "hello",
		},
		{
			name:  "unicode bytes",
			input: []byte("你好，世界"), //nolint:gosmopolitan
			want:  "你好，世界",
		},
		{
			name:  "special chars",
			input: []byte("!@#$%^&*()_+"),
			want:  "!@#$%^&*()_+",
		},
		{
			name:  "default bytes",
			input: []byte("0D000721, Ciallo～(∠・ω< )⌒★"),
			want:  "0D000721, Ciallo～(∠・ω< )⌒★",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, BytesToString(tc.input))
		})
	}

	t.Run("test long bytes", func(t *testing.T) {
		t.Parallel()

		long := make([]byte, 10000)
		for i := range long {
			long[i] = 'a'
		}

		assert.Equal(t, string(long), BytesToString(long))
	})
}

func TestStringToBytes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		input string
		want  []byte
	}{
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "ascii string",
			input: "hello",
			want:  []byte("hello"),
		},
		{
			name:  "unicode string",
			input: "你好，世界", //nolint:gosmopolitan
			want:  []byte("你好，世界"),
		},
		{
			name:  "special chars string",
			input: "!@#$%^&*()_+",
			want:  []byte("!@#$%^&*()_+"),
		},
		{
			name:  "default string",
			input: "0D000721, Ciallo～(∠・ω< )⌒★",
			want:  []byte("0D000721, Ciallo～(∠・ω< )⌒★"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.input == "" {
				assert.Nil(t, StringToBytes(tc.input))
			} else {
				assert.Equal(t, tc.want, StringToBytes(tc.input))
			}
		})
	}

	t.Run("long string", func(t *testing.T) {
		t.Parallel()

		long := make([]byte, 10000)
		for i := range long {
			long[i] = 'a'
		}

		assert.Equal(t, long, StringToBytes(string(long)))
	})
}
