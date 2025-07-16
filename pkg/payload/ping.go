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
		Uid     string `json:"uid,omitempty"`
		Request *struct {
			UserId string `json:"user_id,omitempty"`
			Params string `json:"params,omitempty"`
			Ts     int    `json:"ts,omitempty"`
			Sign   string `json:"sign,omitempty"`
		} `json:"request,omitempty"`
	} `json:"data,omitempty"`
}
