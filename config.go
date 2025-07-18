package afdian

type Config struct {
	BaseURL string

	// Auth configuration
	UserID   string
	APIToken string

	// WebHook configuration
	WebHookListenAddr string
	WebHookPath       string
}

func (c *Config) Default() {
	if c.BaseURL == "" {
		c.BaseURL = "https://afdian.com/api"
	}

	if c.WebHookListenAddr != "" && c.WebHookPath == "" {
		c.WebHookPath = "/"
	}
}
