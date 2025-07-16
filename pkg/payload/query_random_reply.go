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
		// --- ERROR ---
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

		// --- NORMAL ---
		List []struct {
			OutTradeNo string `json:"out_trade_no,omitempty"`
			Content    string `json:"content,omitempty"`
		} `json:"list,omitempty"`
	} `json:"data,omitempty"`
}
