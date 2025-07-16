package payload

import (
	"net/http"
)

type APIBase struct {
	// --- INJECTED RAW ---
	RawResponse *http.Response `json:"-"`

	// --- PAYLOAD ---
	Base
	Explain string `json:"explain,omitempty"`
	Debug   *struct {
		KvString string `json:"kv_string,omitempty"`
	} `json:"debug,omitempty"`
	Request *struct {
		UserId string `json:"user_id,omitempty"`
		Params string `json:"params,omitempty"`
		Ts     int    `json:"ts,omitempty"`
		Sign   string `json:"sign,omitempty"`
	} `json:"request,omitempty"`
}
