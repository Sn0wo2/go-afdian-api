package utils

import (
	"fmt"
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/helper"
	"github.com/Sn0wo2/go-afdian-api/internal/sign"
	jsoniter "github.com/json-iterator/go"
)

// BuildParams builds the final request parameters.
func BuildParams(userID, apiToken string, params map[string]string, ts ...int64) ([]byte, error) {
	for k, v := range params {
		if v == "" {
			delete(params, k)
		}
	}

	paramsJSON, err := jsoniter.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	timestamp := time.Now().Unix()
	if len(ts) > 0 {
		timestamp = ts[0]
	}

	signed, err := sign.APISignParams(userID, apiToken, paramsJSON, timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to sign params: %w", err)
	}

	paramsStr := ""
	if len(params) > 0 { // thanks afdian json string in json
		paramsStr = helper.BytesToString(paramsJSON)
	}

	return jsoniter.Marshal(struct {
		UserID string `json:"user_id,omitempty"`
		Params string `json:"params,omitempty"`
		Ts     int64  `json:"ts,omitempty"`
		Sign   string `json:"sign,omitempty"`
	}{
		UserID: userID,
		Params: paramsStr,
		Ts:     timestamp,
		Sign:   signed,
	})
}
