package afdian

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigSetDefaults(t *testing.T) {
	t.Parallel()

	t.Run("default base url", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		cfg.setDefaults()
		assert.Equal(t, "https://afdian.com/api", cfg.BaseURL)
	})

	t.Run("custom base url", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{BaseURL: "https://custom.url"}
		cfg.setDefaults()
		assert.Equal(t, "https://custom.url", cfg.BaseURL)
	})

	t.Run("default webhook path", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{WebHookListenAddr: ":8080"}
		cfg.setDefaults()
		assert.Equal(t, "/", cfg.WebHookPath)
	})

	t.Run("custom webhook path", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{WebHookListenAddr: ":8080", WebHookPath: "/hook"}
		cfg.setDefaults()
		assert.Equal(t, "/hook", cfg.WebHookPath)
	})

	t.Run("no default webhook path", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		cfg.setDefaults()
		assert.Empty(t, cfg.WebHookPath)
	})
}
