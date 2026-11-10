package afdian

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Sn0wo2/go-afdian-api/v2/internal/utils"
	"github.com/Sn0wo2/go-afdian-api/v2/pkg/payload"
)

func (c *Client) doRequest[P utils.APIResponder](p P, endpoint string, params map[string]string) (P, error) {
	resp, err := c.Send(endpoint, params) //nolint:bodyclose // false positive: utils.Response.Process closes the body
	if err != nil {
		return p, err
	}

	if err := (utils.Response{}).Process(resp, p); err != nil {
		return p, err
	}

	return p, nil
}

func (c *Client) Ping() (*payload.Ping, error) {
	return c.doRequest(new(payload.Ping), "/open/ping", map[string]string{"unix": strconv.FormatInt(time.Now().Unix(), 10)})
}

func (c *Client) QueryRandomReply(outTradeNo string, outTradeNos ...string) (*payload.QueryRandomReply, error) {
	return c.doRequest(new(payload.QueryRandomReply), "/open/query-random-reply", map[string]string{"out_trade_no": strings.Join(slices.Concat([]string{outTradeNo}, outTradeNos), ",")})
}

func (c *Client) QueryOrder(page, perPage int, outTradeNo ...string) (*payload.QueryOrder, error) {
	return c.doRequest(new(payload.QueryOrder), "/open/query-order", map[string]string{"page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "out_trade_no": strings.Join(outTradeNo, ",")})
}

func (c *Client) QuerySponsor(page, perPage int, userIDs ...string) (*payload.QuerySponsor, error) {
	return c.doRequest(new(payload.QuerySponsor), "/open/query-sponsor", map[string]string{"page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "user_id": strings.Join(userIDs, ",")})
}

func (c *Client) QueryPlan(planID string) (*payload.QueryPlan, error) {
	return c.doRequest(new(payload.QueryPlan), "/open/query-plan", map[string]string{"plan_id": planID})
}

func (c *Client) SendMsg(recipient, content string) (*payload.SendMsg, error) {
	return c.doRequest(new(payload.SendMsg), "/open/send-msg", map[string]string{"recipient": recipient, "content": content})
}

const (
	UpdateRandomReplyTypeAppend    = "append"
	UpdateRandomReplyTypeOverwrite = "overwrite"
)

type UpdatePlanReplyOptions struct {
	// PlanID is the subscription plan ID; mutually exclusive with SkuID.
	PlanID string
	// SkuID is the product SKU ID; mutually exclusive with PlanID.
	SkuID string
	// AutoReply overwrites the fixed auto reply when non-empty; empty means no update.
	AutoReply string
	// AutoRandomReply updates the random reply when non-empty; empty means no update.
	AutoRandomReply string
	// UpdateRandomReplyType is required when AutoRandomReply is set.
	UpdateRandomReplyType string
}

func (c *Client) UpdatePlanReply(opts UpdatePlanReplyOptions) (*payload.UpdatePlanReply, error) {
	if (opts.PlanID == "") == (opts.SkuID == "") {
		return nil, errors.New("exactly one of PlanID or SkuID must be set")
	}

	switch opts.UpdateRandomReplyType {
	case "", UpdateRandomReplyTypeAppend, UpdateRandomReplyTypeOverwrite:
	default:
		return nil, fmt.Errorf("invalid UpdateRandomReplyType: %q", opts.UpdateRandomReplyType)
	}

	return c.doRequest(new(payload.UpdatePlanReply), "/open/update-plan-reply", map[string]string{
		"plan_id":                  opts.PlanID,
		"sku_id":                   opts.SkuID,
		"auto_reply":               opts.AutoReply,
		"auto_random_reply":        opts.AutoRandomReply,
		"update_random_reply_type": opts.UpdateRandomReplyType,
	})
}
