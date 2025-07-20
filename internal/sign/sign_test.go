package sign

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"testing"
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/helper"
	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPISignParams(t *testing.T) {
	t.Parallel()

	const userID = "user123"

	const apiToken = "token123"

	params := []byte(`{"key":"value"}`)
	ts := time.Now().Unix()

	actualSign, err := APISignParams(userID, apiToken, params, ts)
	require.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("%x", md5.Sum(helper.StringToBytes(fmt.Sprintf("%sparams%sts%duser_id%s", apiToken, string(params), ts, userID)))), actualSign) //nolint:gosec
}

func TestWebHookSignVerify(t *testing.T) {
	t.Parallel()
	t.Run("empty sign", func(t *testing.T) {
		t.Parallel()
		assert.EqualError(t, WebHookSignVerify(newTestWebHook(t, "")), "sign is empty")
	})

	t.Run("invalid base64 sign", func(t *testing.T) {
		t.Parallel()
		assert.Error(t, WebHookSignVerify(newTestWebHook(t, "invalid base64")))
	})

	t.Run("invalid public key", func(t *testing.T) {
		t.Parallel()

		rb := publicKeyPEM
		publicKeyPEM = `-----BEGIN PUBLIC KEY-----\ninvalid\n-----END PUBLIC KEY-----`

		assert.Error(t, WebHookSignVerify(newTestWebHook(t, "EIUrs8kvspg0MqLvpODkQqpyQLppXCNGPDe8+vvfER9VCzGImjvVJiPUd2UWaGP3A9sS/Ov+hFGwHAdvrm4i9Bte3kNQvwAXcaZx2g06JsngEr4MCe+nn0JHm+mtK9np8N5gJ3DC3GAs6l88SmTnMeJ4no+bPexfqPTYAMWs26e0mexjyLY+8f5l9Zv9rQAz+i/kVjtGNgXFQF34+hAmfOoSxlTJth41XMeVLHxi46kHb+tJNGkU4vLrnMftos2yMzHO+rDP10N7o7VNGMO37aCWfJ+uwKgcXBo0xVbPPAVxHZ+GtOnlEpasINnmoJW7maaoOpRX+IUIxd/U4gTnxA==")))

		publicKeyPEM = rb
	})

	//nolint:paralleltest
	t.Run("valid public key", func(t *testing.T) {
		assert.NoError(t, WebHookSignVerify(newTestWebHook(t, "EIUrs8kvspg0MqLvpODkQqpyQLppXCNGPDe8+vvfER9VCzGImjvVJiPUd2UWaGP3A9sS/Ov+hFGwHAdvrm4i9Bte3kNQvwAXcaZx2g06JsngEr4MCe+nn0JHm+mtK9np8N5gJ3DC3GAs6l88SmTnMeJ4no+bPexfqPTYAMWs26e0mexjyLY+8f5l9Zv9rQAz+i/kVjtGNgXFQF34+hAmfOoSxlTJth41XMeVLHxi46kHb+tJNGkU4vLrnMftos2yMzHO+rDP10N7o7VNGMO37aCWfJ+uwKgcXBo0xVbPPAVxHZ+GtOnlEpasINnmoJW7maaoOpRX+IUIxd/U4gTnxA==")))
	})
}

func newTestWebHook(t *testing.T, sign string) *payload.WebHook {
	t.Helper()

	// afdian private key signature
	const rawJSON = `{"ec":200,"em":"ok","data":{"type":"order","order":{"out_trade_no":"202106232138371083454010626","plan_title":"\u6d4b\u8bd5webhook","user_private_id":"345357fe8374811eaacee52540025c377","user_id":"adf397fe8374811eaacee52540025c377","plan_id":"a45353328af911eb973052540025c377","title":"\u8d5e\u52a9 5 = 5.00 x 1 \u6708","month":1,"total_amount":"5.00","show_amount":"5.00","status":2,"remark":"","redeem_id":"","product_type":1,"discount":"0.00","sku_detail":[{"sku_id":"2172ea4e3a2311edbcaa52540025c377","count":2,"name":"sku1","album_id":"","pic":"","stock":"","post_id":""},{"sku_id":"2172ea4e3a2311edbcaa52540025c378","count":3,"name":"sku2","album_id":"","pic":"","stock":"","post_id":""}],"address_person":"","address_phone":"","address_address":""},"sign":"EIUrs8kvspg0MqLvpODkQqpyQLppXCNGPDe8+vvfER9VCzGImjvVJiPUd2UWaGP3A9sS\/Ov+hFGwHAdvrm4i9Bte3kNQvwAXcaZx2g06JsngEr4MCe+nn0JHm+mtK9np8N5gJ3DC3GAs6l88SmTnMeJ4no+bPexfqPTYAMWs26e0mexjyLY+8f5l9Zv9rQAz+i\/kVjtGNgXFQF34+hAmfOoSxlTJth41XMeVLHxi46kHb+tJNGkU4vLrnMftos2yMzHO+rDP10N7o7VNGMO37aCWfJ+uwKgcXBo0xVbPPAVxHZ+GtOnlEpasINnmoJW7maaoOpRX+IUIxd\/U4gTnxA=="}}`

	var p payload.WebHook

	assert.NoError(t, jsoniter.Unmarshal(helper.StringToBytes(rawJSON), &p))

	p.Data.Sign = sign

	return &p
}
