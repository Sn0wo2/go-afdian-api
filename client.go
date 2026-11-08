package afdian

import (
	"bytes"
	"crypto/md5" //nolint:gosec
	"encoding/json/v2"
	"fmt"
	"net/http"
	"time"
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

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://afdian.com/api"
	}

	if cfg.UserAgent == "" {
		cfg.UserAgent = userAgent()
	}

	if cfg.WebHookListenAddr != "" && cfg.WebHookPath == "" {
		cfg.WebHookPath = "/"
	}

	return &Client{cfg: cfg, HTTP: h}
}

func (c *Client) Send(path string, params map[string]string) (*http.Response, error) {
	filtered := make(map[string]string, len(params))

	for k, v := range params {
		if v != "" {
			filtered[k] = v
		}
	}

	// afdian computes the sign over the params JSON with keys sorted
	// (doc: md5(token+请求数据按key排序拼接key和value)), and encoding/json/v2
	// traverses maps non-deterministically by default, so force deterministic
	// (lexicographically sorted) encoding.
	paramsJSON, err := json.Marshal(filtered, json.Deterministic(true))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	ts := time.Now().Unix()

	paramsStr := ""
	if len(filtered) > 0 {
		paramsStr = string(paramsJSON)
	}

	body, err := json.Marshal(struct {
		UserID string `json:"user_id,omitempty"`
		Params string `json:"params,omitempty"`
		Ts     int64  `json:"ts,omitempty"`
		Sign   string `json:"sign,omitempty"`
	}{
		UserID: c.cfg.UserID,
		Params: paramsStr,
		Ts:     ts,
		Sign:   fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%sparams%sts%duser_id%s", c.cfg.APIToken, paramsStr, ts, c.cfg.UserID)))), //nolint:gosec
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.cfg.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	return c.HTTP.Do(req)
}
