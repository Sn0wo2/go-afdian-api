package afdian

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
	"github.com/json-iterator/go"
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

// Send sends an API request
// WARNING 注意资源泄漏, 调用后请及时关闭
func (c *Client) Send(path string, params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", c.cfg.BaseURL, path), bytes.NewBuffer(NewParamsBuilder(c, params).Build()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Ping() (*payload.Ping, error) {
	resp, err := c.Send("/ping", map[string]string{"ping": strconv.FormatInt(time.Now().UnixNano(), 10)})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	p := &payload.Ping{}
	err = jsoniter.Unmarshal(raw, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
