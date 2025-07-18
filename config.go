package afdian

type Config struct {
	BaseURL string

	// Auth configuration
	UserID   string
	APIToken string

	// WebHook configuration
	WebHookListenAddr string
	WebHookPath       string
	WebHookCallback   CallBack
}

func (c *Config) setDefaults() {
	if c.BaseURL == "" {
		c.BaseURL = "https://afdian.com/api"
	}

	if c.WebHookListenAddr != "" && c.WebHookPath == "" {
		c.WebHookPath = "/"
	}
}
