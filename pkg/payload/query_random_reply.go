package payload

import (
	"net/http"
)

// QueryRandomReply
//
// WARNING: Do NOT convert amounts directly to float when performing operations on money (due to floating-point precision errors)!
// Use github.com/shopspring/decimal for handling monetary values.
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
