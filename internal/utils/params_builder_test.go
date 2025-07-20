package utils

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"testing"
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/helper"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildParams(t *testing.T) {
	t.Parallel()

	userID := "test_user_id"
	apiToken := "test_api_token"
	ts := time.Now().Unix()
	params := map[string]string{
		"key1": "value1",
		"key2": "", // this should be removed
	}

	result, err := BuildParams(userID, apiToken, params, ts)
	require.NoError(t, err)
	assert.NotNil(t, result)

	const paramsJSON = `{"key1":"value1"}`

	assert.Equal(t, userID, jsoniter.Get(result, "user_id").ToString())
	assert.JSONEq(t, paramsJSON, jsoniter.Get(result, "params").ToString())
	assert.Equal(t, ts, jsoniter.Get(result, "ts").ToInt64())

	assert.Equal(t, fmt.Sprintf("%x", md5.Sum(helper.StringToBytes(fmt.Sprintf("%sparams%sts%duser_id%s", apiToken, paramsJSON, ts, userID)))), jsoniter.Get(result, "sign").ToString()) //nolint:gosec
}
