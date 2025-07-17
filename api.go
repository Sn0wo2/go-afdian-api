package afdian

import (
	"strconv"
	"strings"
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/utils"
	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
	jsoniter "github.com/json-iterator/go"
)

func (c *Client) Ping() (*payload.Ping, error) {
	resp, err := c.Send("/ping", map[any]any{"unix": strconv.FormatInt(time.Now().Unix(), 10)})
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

func (c *Client) QueryRandomReply(outTradeNo ...string) (*payload.QueryRandomReply, error) {
	resp, err := c.Send("/query-random-reply", map[any]any{"out_trade_no": strings.Join(outTradeNo, ",")})
	if err != nil {
		return nil, err
	}

	qrr := &payload.QueryRandomReply{}
	qrr.RawResponse = resp

	raw, err := utils.ReadAPIResponse(resp)
	if err != nil {
		return qrr, err
	}

	if err := jsoniter.Unmarshal(raw, qrr); err != nil {
		return qrr, err
	}
	return qrr, nil
}

func (c *Client) QueryOrder(page, perPage int, outTradeNo ...string) (*payload.QueryOrder, error) {
	resp, err := c.Send("/query-order", map[any]any{"page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "out_trade_no": strings.Join(outTradeNo, ",")})
	if err != nil {
		return nil, err
	}

	qo := &payload.QueryOrder{}
	qo.RawResponse = resp

	raw, err := utils.ReadAPIResponse(resp)
	if err != nil {
		return qo, err
	}

	if err := jsoniter.Unmarshal(raw, qo); err != nil {
		return qo, err
	}
	return qo, nil
}

func (c *Client) QuerySponsor(page, perPage int, outTradeNo ...string) (*payload.QuerySponsor, error) {
	resp, err := c.Send("/query-sponsor", map[any]any{"page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "out_trade_no": strings.Join(outTradeNo, ",")})
	if err != nil {
		return nil, err
	}

	qs := &payload.QuerySponsor{}
	qs.RawResponse = resp

	raw, err := utils.ReadAPIResponse(resp)
	if err != nil {
		return qs, err
	}

	if err := jsoniter.Unmarshal(raw, qs); err != nil {
		return qs, err
	}
	return qs, nil
}
