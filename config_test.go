package afdian

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_setDefaults(t *testing.T) {
	t.Parallel()

	t.Run("default base url", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		cfg.setDefaults()
		assert.Equal(t, "https://afdian.com/api", cfg.BaseURL)
	})

	t.Run("default webhook path", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{WebHookListenAddr: ":8080"}
		cfg.setDefaults()
		assert.Equal(t, "/", cfg.WebHookPath)
	})

	t.Run("no default webhook path", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		cfg.setDefaults()
		assert.Empty(t, cfg.WebHookPath)
	})
}
