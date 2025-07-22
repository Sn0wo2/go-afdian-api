package afdian

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWebHook(t *testing.T) {
	t.Parallel()

	client := NewClient(&Config{})
	wh := NewWebHook(client)

	assert.Equal(t, client, wh.client)
	assert.Equal(t, wh, client.WebHook)
}

func TestWebHook_Start(t *testing.T) {
	t.Parallel()

	t.Run("no listen addr", func(t *testing.T) {
		t.Parallel()

		client := NewClient(&Config{})
		wh := NewWebHook(client)

		err := wh.Start()
		assert.Error(t, err)
	})
}

func TestWebHook_resolve(t *testing.T) {
	t.Parallel()

	const (
		userID   = "user_id123"
		apiToken = "api_token123"
		path     = "/webhook"
	)

	wh := NewWebHook(NewClient(&Config{
		UserID:            userID,
		APIToken:          apiToken,
		WebHookPath:       path,
		WebHookListenAddr: ":8080",
	}))

	t.Run("invalid path", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		wh.resolve()(w, httptest.NewRequest(http.MethodPost, "/invalid", nil))

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid method", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		wh.resolve()(w, httptest.NewRequest(http.MethodGet, path, nil))

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("empty body", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		wh.resolve()(w, httptest.NewRequest(http.MethodPost, path, nil))

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("read body error", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		wh.resolve()(w, httptest.NewRequest(http.MethodPost, path, &errReader{}))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		wh.resolve()(w, httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString("invalid json")))

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid sign", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		wh.resolve()(w, httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(`{"ec":200,"em":"ok","data":{"sign":"invalid"}}`)))

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("successful request", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		wh.resolve()(w, httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(`{"ec":200,"em":"ok","data":{"type":"order","order":{"out_trade_no":"202106232138371083454010626","plan_title":"\u6d4b\u8bd5webhook","user_private_id":"345357fe8374811eaacee52540025c377","user_id":"adf397fe8374811eaacee52540025c377","plan_id":"a45353328af911eb973052540025c377","title":"\u8d5e\u52a9 5 = 5.00 x 1 \u6708","month":1,"total_amount":"5.00","show_amount":"5.00","status":2,"remark":"","redeem_id":"","product_type":1,"discount":"0.00","sku_detail":[{"sku_id":"2172ea4e3a2311edbcaa52540025c377","count":2,"name":"sku1","album_id":"","pic":"","stock":"","post_id":""},{"sku_id":"2172ea4e3a2311edbcaa52540025c378","count":3,"name":"sku2","album_id":"","pic":"","stock":"","post_id":""}],"address_person":"","address_phone":"","address_address":""},"sign":"EIUrs8kvspg0MqLvpODkQqpyQLppXCNGPDe8+vvfER9VCzGImjvVJiPUd2UWaGP3A9sS\/Ov+hFGwHAdvrm4i9Bte3kNQvwAXcaZx2g06JsngEr4MCe+nn0JHm+mtK9np8N5gJ3DC3GAs6l88SmTnMeJ4no+bPexfqPTYAMWs26e0mexjyLY+8f5l9Zv9rQAz+i\/kVjtGNgXFQF34+hAmfOoSxlTJth41XMeVLHxi46kHb+tJNGkU4vLrnMftos2yMzHO+rDP10N7o7VNGMO37aCWfJ+uwKgcXBo0xVbPPAVxHZ+GtOnlEpasINnmoJW7maaoOpRX+IUIxd\/U4gTnxA=="}}`)))

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestWebHookRunCallback(t *testing.T) {
	t.Parallel()

	var (
		cbCalled bool
		cbErrs   []error
	)

	wh := NewWebHook(NewClient(&Config{
		WebHookCallback: func(p *payload.WebHook, errs ...error) {
			cbCalled = true
			cbErrs = errs
		},
	}))

	t.Run("no callback", func(t *testing.T) {
		t.Parallel()

		NewWebHook(NewClient(&Config{})).runCallback(&payload.WebHook{}, errors.New("test error"))
	})

	//nolint:paralleltest
	t.Run("with errors", func(t *testing.T) {
		cbCalled = false

		wh.runCallback(&payload.WebHook{}, errors.New("test error"))
		assert.True(t, cbCalled)
		assert.Len(t, cbErrs, 1)
	})

	//nolint:paralleltest
	t.Run("no errors", func(t *testing.T) {
		cbCalled = false

		wh.runCallback(&payload.WebHook{})
		assert.True(t, cbCalled)
		assert.Empty(t, cbErrs)
	})
}

func TestWebHookWriteResponse(t *testing.T) {
	t.Parallel()

	wh := &WebHook{}

	t.Run("successful write", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		err := wh.writeResponse(w, &payload.WebHook{Base: payload.Base{EC: http.StatusOK, EM: "OK"}})
		require.NoError(t, err)

		body, err := io.ReadAll(w.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{"ec":200,"em":"OK"}`, string(body))
	})

	t.Run("write with status code", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		err := wh.writeResponse(w, &payload.WebHook{Base: payload.Base{EC: http.StatusBadRequest, EM: "Bad Request"}})
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

type errReader struct{}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}
