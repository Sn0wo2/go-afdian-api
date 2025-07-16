package payload

type QueryRandomReply struct {
	APIBase
	Data *struct {
		List []struct {
			OutTradeNo string `json:"out_trade_no,omitempty"`
			Content    string `json:"content,omitempty"`
		} `json:"list,omitempty"`
	} `json:"data,omitempty"`
}
