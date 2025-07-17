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
