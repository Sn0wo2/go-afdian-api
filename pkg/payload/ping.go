package payload

import (
	"net/http"
)

// Ping payload
type Ping struct {
	// --- INJECTED RAW ---
	RawResponse *http.Response `json:"-"`

	// --- PAYLOAD ---
	APIBase
	Data *struct {
		// --- ERROR ---
		Explain string `json:"explain,omitempty"`
		Debug   *struct {
			KvString string `json:"kv_string,omitempty"`
		} `json:"debug,omitempty"`

		// --- NORMAL ---
		Request *struct {
			UserID string `json:"user_id,omitempty"`
			Params string `json:"params,omitempty"`
			Ts     int    `json:"ts,omitempty"`
			Sign   string `json:"sign,omitempty"`
		} `json:"request,omitempty"`

		UID  string `json:"uid,omitempty"`
		Page any    `json:"page,omitempty"`
	} `json:"data,omitempty"`
}
