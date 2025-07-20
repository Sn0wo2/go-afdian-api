package utils

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadAPIResponse(t *testing.T) {
	t.Parallel()
	t.Run("successful read", func(t *testing.T) {
		t.Parallel()

		raw, err := ReadAPIResponse(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"ec":200,"em":"ok","data":{"type":"order","order":{"out_trade_no":"202106232138371083454010626","plan_title":"\u6d4b\u8bd5webhook","user_private_id":"345357fe8374811eaacee52540025c377","user_id":"adf397fe8374811eaacee52540025c377","plan_id":"a45353328af911eb973052540025c377","title":"\u8d5e\u52a9 5 = 5.00 x 1 \u6708","month":1,"total_amount":"5.00","show_amount":"5.00","status":2,"remark":"","redeem_id":"","product_type":1,"discount":"0.00","sku_detail":[{"sku_id":"2172ea4e3a2311edbcaa52540025c377","count":2,"name":"sku1","album_id":"","pic":"","stock":"","post_id":""},{"sku_id":"2172ea4e3a2311edbcaa52540025c378","count":3,"name":"sku2","album_id":"","pic":"","stock":"","post_id":""}],"address_person":"","address_phone":"","address_address":""},"sign":"EIUrs8kvspg0MqLvpODkQqpyQLppXCNGPDe8+vvfER9VCzGImjvVJiPUd2UWaGP3A9sS\/Ov+hFGwHAdvrm4i9Bte3kNQvwAXcaZx2g06JsngEr4MCe+nn0JHm+mtK9np8N5gJ3DC3GAs6l88SmTnMeJ4no+bPexfqPTYAMWs26e0mexjyLY+8f5l9Zv9rQAz+i\/kVjtGNgXFQF34+hAmfOoSxlTJth41XMeVLHxi46kHb+tJNGkU4vLrnMftos2yMzHO+rDP10N7o7VNGMO37aCWfJ+uwKgcXBo0xVbPPAVxHZ+GtOnlEpasINnmoJW7maaoOpRX+IUIxd\/U4gTnxA=="}}`)),
		})
		require.NoError(t, err)
		assert.JSONEq(t, `{"ec":200,"em":"ok","data":{"type":"order","order":{"out_trade_no":"202106232138371083454010626","plan_title":"\u6d4b\u8bd5webhook","user_private_id":"345357fe8374811eaacee52540025c377","user_id":"adf397fe8374811eaacee52540025c377","plan_id":"a45353328af911eb973052540025c377","title":"\u8d5e\u52a9 5 = 5.00 x 1 \u6708","month":1,"total_amount":"5.00","show_amount":"5.00","status":2,"remark":"","redeem_id":"","product_type":1,"discount":"0.00","sku_detail":[{"sku_id":"2172ea4e3a2311edbcaa52540025c377","count":2,"name":"sku1","album_id":"","pic":"","stock":"","post_id":""},{"sku_id":"2172ea4e3a2311edbcaa52540025c378","count":3,"name":"sku2","album_id":"","pic":"","stock":"","post_id":""}],"address_person":"","address_phone":"","address_address":""},"sign":"EIUrs8kvspg0MqLvpODkQqpyQLppXCNGPDe8+vvfER9VCzGImjvVJiPUd2UWaGP3A9sS\/Ov+hFGwHAdvrm4i9Bte3kNQvwAXcaZx2g06JsngEr4MCe+nn0JHm+mtK9np8N5gJ3DC3GAs6l88SmTnMeJ4no+bPexfqPTYAMWs26e0mexjyLY+8f5l9Zv9rQAz+i\/kVjtGNgXFQF34+hAmfOoSxlTJth41XMeVLHxi46kHb+tJNGkU4vLrnMftos2yMzHO+rDP10N7o7VNGMO37aCWfJ+uwKgcXBo0xVbPPAVxHZ+GtOnlEpasINnmoJW7maaoOpRX+IUIxd\/U4gTnxA=="}}`, string(raw))
	})

	t.Run("non-200 status code", func(t *testing.T) {
		t.Parallel()

		_, err := ReadAPIResponse(&http.Response{
			StatusCode: http.StatusNotFound,
			// ReadAPIResponse first check status code
			Body: io.NopCloser(bytes.NewBufferString("")),
		})
		assert.EqualError(t, err, "unexpected status code: 404")
	})

	t.Run("empty response body", func(t *testing.T) {
		t.Parallel()

		_, err := ReadAPIResponse(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer(nil)),
		})
		assert.EqualError(t, err, "empty response")
	})

	t.Run("read error", func(t *testing.T) {
		t.Parallel()

		readErr := errors.New("read error")
		_, err := ReadAPIResponse(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(&mockReader{reader: bytes.NewBufferString(""), err: readErr}),
		})
		require.Error(t, err)
		// if we don't require error, it has chance panic(?)
		assert.Equal(t, readErr, err)
	})
}

// mockReader implementation of io.Reader return error
type mockReader struct {
	reader io.Reader
	err    error
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.err != nil {
		return 0, m.err
	}

	return m.reader.Read(p)
}
