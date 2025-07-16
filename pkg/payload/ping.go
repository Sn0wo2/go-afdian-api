package payload

type Ping struct {
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
