package afdian

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
	jsoniter "github.com/json-iterator/go"
)

func (c *Client) Ping() (*payload.Ping, error) {
	resp, err := c.Send("/ping", map[string]string{"ping": strconv.FormatInt(time.Now().UnixNano(), 10)})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	p := &payload.Ping{}
	p.RawResponse = resp

	if resp.StatusCode != http.StatusOK {
		return p, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return p, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(raw))
	if len(raw) == 0 {
		return p, fmt.Errorf("empty response")
	}

	err = jsoniter.Unmarshal(raw, p)
	if err != nil {
		return p, err
	}

	p.RawResponse = resp

	return p, nil
}
