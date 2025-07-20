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
	t.Run("test non-empty bytes", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "test", BytesToString([]byte("test")))
	})
}

func TestStringToBytes(t *testing.T) {
	t.Parallel()
	assert.Nil(t, StringToBytes(""))
	assert.Equal(t, []byte("hello"), StringToBytes("hello"))
	assert.Equal(t, []byte("你好，世界"), StringToBytes("你好，世界")) //nolint:gosmopolitan
	assert.Equal(t, []byte("!@#$%^&*()_+"), StringToBytes("!@#$%^&*()_+"))

	long := make([]byte, 10000)
	for i := range long {
		long[i] = 'a'
	}

	assert.Equal(t, long, StringToBytes(string(long)))
}
