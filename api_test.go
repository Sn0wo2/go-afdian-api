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

func TestQueryRandomReply(t *testing.T) {
	t.Parallel()

	reply, err := NewClient(&Config{UserID: "user123", APIToken: "token123"}, &http.Client{
		Transport: &mockRoundTripper{
			roundTrip: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"ec":200,"em":"success","data":{"list":[{"out_trade_no":"202505141538455397541020050","content":"999"}]}}`)),
					Header:     make(http.Header),
				}, nil
			},
		},
	}).QueryRandomReply("202505141538455397541020050")

	require.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, 200, reply.EC)
	assert.Equal(t, "ok", reply.EM)
	assert.NotEmpty(t, reply.Data.List)
}

func TestQuerySponsor(t *testing.T) {
	t.Parallel()

	sponsor, err := NewClient(&Config{UserID: "user123", APIToken: "token123"}, &http.Client{
		Transport: &mockRoundTripper{
			roundTrip: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"ec":200,"em":"","data":{"list":[{"out_trade_no":"202106232138371083454010626","custom_order_id":"Steam12345","user_id":"adf397fe8374811eaacee52540025c377","user_private_id":"33","plan_id":"a45353328af911eb973052540025c377","month":1,"total_amount":"5.00","show_amount":"5.00","status":2,"remark":"","redeem_id":"","product_type":0,"discount":"0.00","sku_detail":[{"sku_id":"b082342c4aba11ebb5cb52540025c377","count":1,"name":"15000 赏金/货币 兑换码","album_id":"","pic":"https://pic1.afdiancdn.com/user/8a8e408a3aeb11eab26352540025c377/common/sfsfsff.jpg"}],"address_person":"","address_phone":"","address_address":""}],"total_count":167,"total_page":11}}`)),
					Header:     make(http.Header),
				}, nil
			},
		},
	}).QuerySponsor(1, 100)

	require.NoError(t, err)
	assert.NotNil(t, sponsor)
	assert.Equal(t, 200, sponsor.EC)
	assert.Equal(t, "ok", sponsor.EM)
	assert.Equal(t, 1, sponsor.Data.TotalCount)
	assert.NotEmpty(t, sponsor.Data.List)
}

func TestQueryOrderError(t *testing.T) {
	t.Parallel()

	order, err := NewClient(&Config{UserID: "user123", APIToken: "token123"}, &http.Client{Transport: &mockRoundTripper{
		roundTrip: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"ec":400,"em":"error message from server"}`)),
				Header:     make(http.Header),
			}, nil
		},
	}}).QueryOrder(1, 10)

	require.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, 400, order.EC)
	assert.Equal(t, "error message from server", order.EM)
}

func TestDoRequestInvalidJSON(t *testing.T) {
	t.Parallel()

	ping, err := NewClient(&Config{UserID: "user123", APIToken: "token123"}, &http.Client{Transport: &mockRoundTripper{
		roundTrip: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"ec":200,"em":"ok","data":}`)),
				Header:     make(http.Header),
			}, nil
		},
	}}).Ping()

	require.Error(t, err)
	assert.NotNil(t, ping)
	assert.Equal(t, 0, ping.EC)
}
