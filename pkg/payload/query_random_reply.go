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
	RawResponse *http.Response `json:"-"` // Raw HTTP response

	// --- PAYLOAD ---
	APIBase
	Data *struct {
		// --- ERROR ---
		Explain string `json:"explain,omitempty"` // Error description
		Debug   *struct {
			KvString string `json:"kv_string,omitempty"` // For debugging signature errors
		} `json:"debug,omitempty"`
		Request *struct {
			UserID string `json:"user_id,omitempty"` // User ID
			Params string `json:"params,omitempty"`  // Request parameters
			Ts     int    `json:"ts,omitempty"`      // Timestamp
			Sign   string `json:"sign,omitempty"`    // Signature
		} `json:"request,omitempty"`

		// --- NORMAL ---
		List []struct {
			OutTradeNo string `json:"out_trade_no,omitempty"` // Order number
			Content    string `json:"content,omitempty"`      // Content of the random reply
		} `json:"list,omitempty"`
	} `json:"data,omitempty"`
}
