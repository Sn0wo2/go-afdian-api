package afdian

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/utils"
	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
	jsoniter "github.com/json-iterator/go"
)

type RawResponder interface {
	SetRawResponse(resp *http.Response)
}

func doRequest[T any, P interface {
	*T
	RawResponder
}](c *Client, endpoint string, params map[string]string) (*T, error) {
	resp, err := c.Send(endpoint, params) //nolint:bodyclose
	if err != nil {
		return nil, err
	}

	p := new(T)
	P(p).SetRawResponse(resp)

	raw, err := utils.ReadAPIResponse(resp)
	if err != nil {
		return p, err
	}

	if err := jsoniter.Unmarshal(raw, p); err != nil {
		return p, err
	}

	if checker, ok := any(p).(payload.Checker); ok {
		if ec := checker.GetEC(); ec != http.StatusOK {
			return p, fmt.Errorf("afdian api error: ec=%d, em=%s", ec, checker.GetEM())
		}
	}

	return p, nil
}

func (c *Client) Ping() (*payload.Ping, error) {
	return doRequest[payload.Ping, *payload.Ping](c, "/open/ping", map[string]string{"unix": strconv.FormatInt(time.Now().Unix(), 10)})
}

// QueryRandomReply need one or more outTradeNo
func (c *Client) QueryRandomReply(outTradeNo ...string) (*payload.QueryRandomReply, error) {
	return doRequest[payload.QueryRandomReply, *payload.QueryRandomReply](c, "/open/query-random-reply", map[string]string{"out_trade_no": strings.Join(outTradeNo, ",")})
}

func (c *Client) QueryOrder(page, perPage int, outTradeNo ...string) (*payload.QueryOrder, error) {
	return doRequest[payload.QueryOrder, *payload.QueryOrder](c, "/open/query-order", map[string]string{"page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "out_trade_no": strings.Join(outTradeNo, ",")})
}

func (c *Client) QuerySponsor(page, perPage int, outTradeNo ...string) (*payload.QuerySponsor, error) {
	return doRequest[payload.QuerySponsor, *payload.QuerySponsor](c, "/open/query-sponsor", map[string]string{"page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "out_trade_no": strings.Join(outTradeNo, ",")})
}
