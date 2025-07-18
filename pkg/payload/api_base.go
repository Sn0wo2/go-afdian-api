package payload

import (
	"net/http"
)

type APIBase struct {
	// --- INJECTED RAW ---
	RawResponse *http.Response `json:"-"`

	// --- PAYLOAD ---
	Base
}

// APIDataBase contains the common fields for all API responses.
type APIDataBase struct {
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
}
