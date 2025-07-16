package payload

import (
	"net/http"
)

type QueryRandomReply struct {
	// --- INJECTED RAW ---
	RawResponse *http.Response `json:"-"`

	// --- PAYLOAD ---
	APIBase
	Data *struct {
		List []struct {
			OutTradeNo string `json:"out_trade_no,omitempty"`
			Content    string `json:"content,omitempty"`
		} `json:"list,omitempty"`
	} `json:"data,omitempty"`
}
