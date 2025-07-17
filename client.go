package afdian

import (
	"bytes"
	"fmt"
	"net/http"
)

type Client struct {
	cfg     *Config
	WebHook *WebHook
	HTTP    *http.Client
}

func NewClient(cfg *Config, hc ...*http.Client) *Client {
	h := http.DefaultClient
	if len(hc) > 0 {
		h = hc[0]
	}
	return &Client{cfg: cfg, HTTP: h}
}

// Send an API request
// WARNING: Be aware of potential resource leaks
func (c *Client) Send(path string, params map[any]any) (*http.Response, error) {
	p, err := NewParamsBuilder(c, params).Build()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", c.cfg.BaseURL, path), bytes.NewBuffer(p))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.HTTP.Do(req)
}
