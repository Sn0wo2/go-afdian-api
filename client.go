package afdian

import (
	"bytes"
	"cmp"
	"crypto/md5" //nolint:gosec
	"encoding/json/v2"
	"fmt"
	"maps"
	"net/http"
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/utils"
	"github.com/parquet-go/bitpack/unsafecast"
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

	c := *cfg
	c.BaseURL = cmp.Or(c.BaseURL, "https://afdian.com/api")
	c.UserAgent = cmp.Or(c.UserAgent, utils.UserAgent{}.String())

	if c.WebHookListenAddr != "" && c.WebHookPath == "" {
		c.WebHookPath = "/"
	}

	return &Client{cfg: &c, HTTP: h}
}

func (c *Client) Send(path string, params map[string]string) (*http.Response, error) {
	filtered := maps.Clone(params)
	maps.DeleteFunc(filtered, func(_, v string) bool { return v == "" })

	// afdian computes the sign over the params JSON with keys sorted
	// (doc: md5(token+请求数据按key排序拼接key和value)), and encoding/json/v2
	// traverses maps non-deterministically by default, so force deterministic
	// (lexicographically sorted) encoding.
	paramsJSON, err := json.Marshal(filtered, json.Deterministic(true))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	ts := time.Now().Unix()

	paramsStr := unsafecast.String(paramsJSON)
	if len(filtered) == 0 {
		paramsStr = ""
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
		Sign:   fmt.Sprintf("%x", md5.Sum(unsafecast.Bytes(fmt.Sprintf("%sparams%sts%duser_id%s", c.cfg.APIToken, paramsStr, ts, c.cfg.UserID)))), //nolint:gosec
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
