package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToString(t *testing.T) {
	t.Parallel()
	t.Run("test empty bytes", func(t *testing.T) {
		t.Parallel()
		assert.Empty(t, BytesToString(nil))
	})
	t.Run("test ascii bytes", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "hello", BytesToString([]byte("hello")))
	})
	t.Run("test unicode bytes", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "你好，世界", BytesToString([]byte("你好，世界"))) //nolint:gosmopolitan
	})
	t.Run("test special chars", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "!@#$%^&*()_+", BytesToString([]byte("!@#$%^&*()_+")))
	})
	t.Run("test long bytes", func(t *testing.T) {
		t.Parallel()

		long := make([]byte, 10000)
		for i := range long {
			long[i] = 'a'
		}

		assert.Equal(t, string(long), BytesToString(long))
	})

	t.Run("test default bytes", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "0D000721, Ciallo～(∠・ω< )⌒★", BytesToString([]byte("0D000721, Ciallo～(∠・ω< )⌒★")))
	})
}

func TestStringToBytes(t *testing.T) {
	t.Parallel()

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, StringToBytes(""))
	})

	t.Run("ascii string", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, []byte("hello"), StringToBytes("hello"))
	})

	t.Run("unicode string", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, []byte("你好，世界"), StringToBytes("你好，世界")) //nolint:gosmopolitan
	})

	t.Run("special chars string", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, []byte("!@#$%^&*()_+"), StringToBytes("!@#$%^&*()_+"))
	})

	t.Run("long string", func(t *testing.T) {
		t.Parallel()
		long := make([]byte, 10000)
		for i := range long {
			long[i] = 'a'
		}

		assert.Equal(t, long, StringToBytes(string(long)))
	})

	t.Run("default string", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, []byte("0D000721, Ciallo～(∠・ω< )⌒★"), StringToBytes("0D000721, Ciallo～(∠・ω< )⌒★"))
	})
}
