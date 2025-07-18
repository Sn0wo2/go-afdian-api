package afdian

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Sn0wo2/go-afdian-api/internal/utils"
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

	cfg.setDefaults()

	return &Client{cfg: cfg, HTTP: h}
}

// Send an API request
// WARNING: Be aware of potential resource leaks
func (c *Client) Send(path string, params map[string]string) (*http.Response, error) {
	p, err := utils.BuildParams(c.cfg.UserID, c.cfg.APIToken, params)
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
