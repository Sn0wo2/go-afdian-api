package payload

import (
	"net/http"
)

type APIBase struct {
	RawResponse *http.Response `json:"-"`

	Base
}

type APIDataBase struct {
	Explain string `json:"explain,omitempty"`
	Debug   *struct {
		KvString string `json:"kv_string,omitempty"`
	} `json:"debug,omitempty"`

	Request *struct {
		UserID string `json:"user_id,omitempty"`
		Params string `json:"params,omitempty"`
		Ts     int    `json:"ts,omitempty"`
		Sign   string `json:"sign,omitempty"`
	} `json:"request,omitempty"`
}

func (b *APIBase) SetRawResponse(resp *http.Response) {
	b.RawResponse = resp
}
