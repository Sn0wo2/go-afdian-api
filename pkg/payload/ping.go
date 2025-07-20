package payload

import (
	"net/http"
)

// Ping payload
type Ping struct {
	// --- INJECTED RAW ---
	RawResponse *http.Response `json:"-"` // Raw HTTP response

	// --- PAYLOAD ---
	APIBase
	Data *struct {
		APIDataBase
		UID  string `json:"uid,omitempty"`
		Page any    `json:"page,omitempty"`
	} `json:"data,omitempty"`
}

func (p *Ping) SetRawResponse(resp *http.Response) {
	p.RawResponse = resp
}
