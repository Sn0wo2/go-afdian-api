package payload

// OAuthAccessToken payload
//
// NOTE: Afdian's OAuth2 implementation is not fully OAuth2-compliant:
// exchanging the authorization code does not return an access_token or
// refresh_token, it directly returns the identity of the authorized user.
type OAuthAccessToken struct {
	// --- PAYLOAD ---
	APIBase
	Data *struct {
		UserID        string `json:"user_id,omitempty"`         // Website user ID
		UserPrivateID string `json:"user_private_id,omitempty"` // Same as user_private_id of the open API, ignore if unused
		Name          string `json:"name,omitempty"`            // Nickname (not unique)
		Avatar        string `json:"avatar,omitempty"`          // Avatar URL
	} `json:"data,omitempty"`
}
