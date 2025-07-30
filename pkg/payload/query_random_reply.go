package payload

import (
	"net/http"
)

// QueryRandomReply
//
// WARNING: Do NOT convert amounts directly to float when performing operations on money (due to floating-point precision errors)!
// Use github.com/shopspring/decimal for handling monetary values.
type QueryRandomReply struct {
	// --- PAYLOAD ---
	APIBase
	Data *struct {
		APIDataBase
		// --- NORMAL ---
		List []struct {
			OutTradeNo string `json:"out_trade_no,omitempty"` // Order number
			Content    string `json:"content,omitempty"`      // Content of the random reply
		} `json:"list,omitempty"`
	} `json:"data,omitempty"`
}

func (r *QueryRandomReply) SetRawResponse(resp *http.Response) {
	r.RawResponse = resp
}
