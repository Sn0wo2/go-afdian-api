package afdian

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	t.Parallel()

	ping, err := NewClient(&Config{UserID: "user123", APIToken: "token123"}, &http.Client{
		Transport: &mockRoundTripper{
			roundTrip: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"ec":200,"em":"pong","data":{"uid":"user123","request":{"user_id":"user123","params":"{\"unix\":\"1752932834\"}","ts":1752932834,"sign":"93b6b59e1d007c8b045dd86e8bad117b"},"page":null}}`)),
					Header:     make(http.Header),
				}, nil
			},
		},
	}).Ping()

	require.NoError(t, err)
	assert.NotNil(t, ping)
	assert.Equal(t, 200, ping.EC)
	assert.Equal(t, "user123", ping.Data.UID)
	assert.Equal(t, "pong", ping.EM)
	assert.Equal(t, 1752932834, ping.Data.Request.Ts)
}

func TestQueryOrder_Error(t *testing.T) {
	t.Parallel()

	mockRT := &mockRoundTripper{
		roundTrip: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"ec":400,"em":"error message from server"}`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	mockHTTPClient := &http.Client{Transport: mockRT}
	cfg := &Config{UserID: "user123", APIToken: "token123"}
	client := NewClient(cfg, mockHTTPClient)

	order, err := client.QueryOrder(1, 10)

	require.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, 400, order.EC)
	assert.Equal(t, "error message from server", order.EM)
}
