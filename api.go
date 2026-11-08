package afdian

import (
	"bytes"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
)

type rawResponder interface {
	SetRawResponse(resp *http.Response)
}

// processResponse reads the response body, restores it for later use,
// unmarshals it into target and checks the HTTP status code as well as
// the afdian-specific EC business code.
func processResponse(resp *http.Response, raw rawResponder, target any) error {
	raw.SetRawResponse(resp)

	rawBody, readErr := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	resp.Body = io.NopCloser(bytes.NewReader(rawBody))

	if readErr != nil {
		return readErr
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if len(rawBody) == 0 {
		return errors.New("empty response")
	}

	if err := json.Unmarshal(rawBody, target); err != nil {
		return err
	}

	if checker, ok := target.(payload.Checker); ok {
		if ec := checker.GetEC(); ec != http.StatusOK {
			return fmt.Errorf("afdian api error: ec=%d, em=%s", ec, checker.GetEM())
		}
	}

	return nil
}

func doRequest[T any, P interface {
	*T
	rawResponder
}](c *Client, endpoint string, params map[string]string) (*T, error) {
	resp, err := c.Send(endpoint, params)
	if err != nil {
		return nil, err
	}

	p := P(new(T))

	if err := processResponse(resp, p, p); err != nil {
		return p, err
	}

	return p, nil
}

func (c *Client) Ping() (*payload.Ping, error) {
	return doRequest[payload.Ping, *payload.Ping](c, "/open/ping", map[string]string{"unix": strconv.FormatInt(time.Now().Unix(), 10)})
}

func (c *Client) QueryRandomReply(outTradeNo string, outTradeNos ...string) (*payload.QueryRandomReply, error) {
	return doRequest[payload.QueryRandomReply, *payload.QueryRandomReply](c, "/open/query-random-reply", map[string]string{"out_trade_no": strings.Join(append([]string{outTradeNo}, outTradeNos...), ",")})
}

func (c *Client) QueryOrder(page, perPage int, outTradeNo ...string) (*payload.QueryOrder, error) {
	return doRequest[payload.QueryOrder, *payload.QueryOrder](c, "/open/query-order", map[string]string{"page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "out_trade_no": strings.Join(outTradeNo, ",")})
}

func (c *Client) QuerySponsor(page, perPage int, userIDs ...string) (*payload.QuerySponsor, error) {
	return doRequest[payload.QuerySponsor, *payload.QuerySponsor](c, "/open/query-sponsor", map[string]string{"page": strconv.Itoa(page), "per_page": strconv.Itoa(perPage), "user_id": strings.Join(userIDs, ",")})
}

func (c *Client) QueryPlan(planID string) (*payload.QueryPlan, error) {
	return doRequest[payload.QueryPlan, *payload.QueryPlan](c, "/open/query-plan", map[string]string{"plan_id": planID})
}

func (c *Client) SendMsg(recipient, content string) (*payload.SendMsg, error) {
	return doRequest[payload.SendMsg, *payload.SendMsg](c, "/open/send-msg", map[string]string{"recipient": recipient, "content": content})
}

const (
	UpdateRandomReplyTypeAppend    = "append"    // append to the existing random reply, afdian prepends a newline
	UpdateRandomReplyTypeOverwrite = "overwrite" // overwrite the existing random reply
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

	return doRequest[payload.UpdatePlanReply, *payload.UpdatePlanReply](c, "/open/update-plan-reply", map[string]string{
		"plan_id":                  opts.PlanID,
		"sku_id":                   opts.SkuID,
		"auto_reply":               opts.AutoReply,
		"auto_random_reply":        opts.AutoRandomReply,
		"update_random_reply_type": opts.UpdateRandomReplyType,
	})
}
