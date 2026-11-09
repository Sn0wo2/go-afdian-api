package afdian

import (
	"testing"

	"github.com/Sn0wo2/go-afdian-api/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestConfigSetDefaults(t *testing.T) {
	t.Parallel()

	t.Run("default base url", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		client := NewClient(cfg)
		assert.Equal(t, "https://afdian.com/api", client.cfg.BaseURL)
		assert.Empty(t, cfg.BaseURL)
	})

	t.Run("custom base url", func(t *testing.T) {
		t.Parallel()

		client := NewClient(&Config{BaseURL: "https://custom.url"})
		assert.Equal(t, "https://custom.url", client.cfg.BaseURL)
	})

	t.Run("default webhook path", func(t *testing.T) {
		t.Parallel()

		client := NewClient(&Config{WebHookListenAddr: ":8080"})
		assert.Equal(t, "/", client.cfg.WebHookPath)
	})

	t.Run("custom webhook path", func(t *testing.T) {
		t.Parallel()

		client := NewClient(&Config{WebHookListenAddr: ":8080", WebHookPath: "/hook"})
		assert.Equal(t, "/hook", client.cfg.WebHookPath)
	})

	t.Run("default user agent", func(t *testing.T) {
		t.Parallel()

		client := NewClient(&Config{})
		assert.Equal(t, utils.UserAgent{}.String(), client.cfg.UserAgent)
	})

	t.Run("custom user agent", func(t *testing.T) {
		t.Parallel()

		client := NewClient(&Config{UserAgent: "my-app/1.0"})
		assert.Equal(t, "my-app/1.0", client.cfg.UserAgent)
	})

	t.Run("no default webhook path", func(t *testing.T) {
		t.Parallel()

		client := NewClient(&Config{})
		assert.Empty(t, client.cfg.WebHookPath)
	})
}
