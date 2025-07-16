package afdian

import (
	"strconv"
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/utils"
	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
	jsoniter "github.com/json-iterator/go"
)

func (c *Client) Ping() (*payload.Ping, error) {
	resp, err := c.Send("/ping", map[string]string{"unix": strconv.FormatInt(time.Now().Unix(), 10)})
	if err != nil {
		return nil, err
	}

	p := &payload.Ping{}
	p.RawResponse = resp

	raw, err := utils.ReadAPIResponse(resp)
	if err != nil {
		return p, err
	}

	if err := jsoniter.Unmarshal(raw, p); err != nil {
		return p, err
	}
	return p, nil
}
