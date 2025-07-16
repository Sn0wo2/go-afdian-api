package afdian

import (
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/helper"
	"github.com/json-iterator/go"
)

type ParamsBuilder struct {
	client *Client
	Params map[string]string
}

func NewParamsBuilder(client *Client, params map[string]string) *ParamsBuilder {
	return &ParamsBuilder{
		client: client,
		Params: params,
	}
}

type params struct {
	UserID string `json:"user_id"`
	Params string `json:"params"`
	Ts     int64  `json:"ts"`
	Sign   string `json:"sign"`
}

func (b *ParamsBuilder) Build() ([]byte, error) {
	paramsJSON, err := jsoniter.Marshal(b.Params)
	if err != nil {
		return nil, err
	}

	ts := time.Now().Unix()
	sign, err := APISignParams(b.client.cfg.UserID, b.client.cfg.APIToken, paramsJSON, ts)
	if err != nil {
		return nil, err
	}
	pJSON, err := jsoniter.Marshal(params{
		UserID: b.client.cfg.UserID,
		Params: helper.BytesToString(paramsJSON),
		Ts:     ts,
		Sign:   sign,
	})
	if err != nil {
		return nil, err
	}
	return pJSON, err
}
