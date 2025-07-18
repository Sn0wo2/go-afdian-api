package payload

import (
	"net/http"
)

// Ping payload
type Ping struct {
	// --- INJECTED RAW ---
	RawResponse *http.Response `json:"-"` // Raw HTTP response

	// --- PAYLOAD ---
	APIBase
	Data *struct {
		// --- ERROR ---
		Explain string `json:"explain,omitempty"` // Error description
		Debug   *struct {
			KvString string `json:"kv_string,omitempty"` // For debugging signature errors
		} `json:"debug,omitempty"`

		// --- NORMAL ---
		Request *struct {
			UserID string `json:"user_id,omitempty"` // User ID
			Params string `json:"params,omitempty"`  // Request parameters
			Ts     int    `json:"ts,omitempty"`      // Timestamp
			Sign   string `json:"sign,omitempty"`    // Signature
		} `json:"request,omitempty"`

		UID  string `json:"uid,omitempty"`  // User ID
		Page any    `json:"page,omitempty"` // Page number
	} `json:"data,omitempty"`
}
