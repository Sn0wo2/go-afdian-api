package afdian

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/helper"
	"github.com/json-iterator/go"
)

type ParamsBuilder interface {
	Build() []byte
}

type builder struct {
	client *Client
	params map[string]string
}

func NewParamsBuilder(client *Client, params map[string]string) ParamsBuilder {
	return &builder{
		client: client,
		params: params,
	}
}

func (b *builder) Build() []byte {
	paramsJSON, err := jsoniter.Marshal(b.params)
	if err != nil {
		return nil
	}

	ts := time.Now().Unix()

	type Pl struct {
		UserID string `json:"user_id"`
		Params string `json:"params"`
		Ts     int64  `json:"ts"`
		Sign   string `json:"sign"`
	}
	plJSON, err := jsoniter.Marshal(Pl{
		UserID: b.client.cfg.UserID,
		Params: helper.BytesToString(paramsJSON),
		Ts:     ts,
		Sign:   fmt.Sprintf("%x", md5.Sum(helper.StringToBytes(fmt.Sprintf("%sparams%sts%duser_id%s", b.client.cfg.APIToken, helper.BytesToString(paramsJSON), ts, b.client.cfg.UserID)))),
	})
	if err != nil {
		return nil
	}
	return plJSON
}
