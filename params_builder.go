package afdian

import (
	"fmt"
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/helper"
	"github.com/Sn0wo2/go-afdian-api/internal/sign"
	"github.com/json-iterator/go"
)

type ParamsBuilder struct {
	client *Client
	Params map[any]any
}

func newParamsBuilder(client *Client, params map[any]any) *ParamsBuilder {
	return &ParamsBuilder{
		client: client,
		Params: params,
	}
}

func (b *ParamsBuilder) Build() ([]byte, error) {
	p := make(map[string]string, len(b.Params))
	for k, v := range b.Params {
		if k == nil || k == "" || v == nil || v == "" {
			continue
		}
		p[fmt.Sprint(k)] = fmt.Sprint(v)
	}
	paramsJSON, err := jsoniter.Marshal(p)
	if err != nil {
		return nil, err
	}

	ts := time.Now().Unix()
	s, err := sign.APISignParams(b.client.cfg.UserID, b.client.cfg.APIToken, paramsJSON, ts)
	if err != nil {
		return nil, err
	}

	type params struct {
		UserID string `json:"user_id"`
		Params string `json:"params"`
		Ts     int64  `json:"ts"`
		Sign   string `json:"sign"`
	}
	pJSON, err := jsoniter.Marshal(params{
		UserID: b.client.cfg.UserID,
		Params: helper.BytesToString(paramsJSON),
		Ts:     ts,
		Sign:   s,
	})
	if err != nil {
		return nil, err
	}
	return pJSON, err
}
