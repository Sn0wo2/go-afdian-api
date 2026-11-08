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
		NewClient(cfg)
		assert.Equal(t, "https://afdian.com/api", cfg.BaseURL)
	})

	t.Run("custom base url", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{BaseURL: "https://custom.url"}
		NewClient(cfg)
		assert.Equal(t, "https://custom.url", cfg.BaseURL)
	})

	t.Run("default webhook path", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{WebHookListenAddr: ":8080"}
		NewClient(cfg)
		assert.Equal(t, "/", cfg.WebHookPath)
	})

	t.Run("custom webhook path", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{WebHookListenAddr: ":8080", WebHookPath: "/hook"}
		NewClient(cfg)
		assert.Equal(t, "/hook", cfg.WebHookPath)
	})

	t.Run("default user agent", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		NewClient(cfg)
		assert.Equal(t, userAgent(), cfg.UserAgent)
	})

	t.Run("custom user agent", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{UserAgent: "my-app/1.0"}
		NewClient(cfg)
		assert.Equal(t, "my-app/1.0", cfg.UserAgent)
	})

	t.Run("no default webhook path", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		NewClient(cfg)
		assert.Empty(t, cfg.WebHookPath)
	})
}
