package afdian

type Config struct {
	BaseURL   string
	UserAgent string
	UserID    string
	APIToken  string

	OAuthClientID     string
	OAuthClientSecret string
	OAuthScope        string // OAuth2 scope, only "basic" is documented so far, defaults to "basic" when empty
	OAuthRedirectURI  string // Default OAuth2 redirect_uri, can be overridden per call

	WebHookListenAddr string
	WebHookPath       string
	WebHookCallback   CallBack
}
