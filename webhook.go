package afdian

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Sn0wo2/go-afdian-api/internal/sign"
	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
	jsoniter "github.com/json-iterator/go"
)

type CallBack func(p *payload.WebHook, errs ...error)

type WebHook struct {
	client     *Client
	HttpServer *http.Server
	callback   CallBack
}

func NewWebHook(client *Client) *WebHook {
	wh := &WebHook{client: client}
	client.WebHook = wh

	return wh
}

// SetCallback
// Must implement idempotent logic
func (wh *WebHook) SetCallback(callback CallBack) {
	wh.callback = callback
}

func (wh *WebHook) Start() error {
	if wh.client.cfg.WebHookListenAddr == "" {
		return errors.New("WebHookListenAddr is empty")
	}

	server := http.Server{
		Addr:    wh.client.cfg.WebHookListenAddr,
		Handler: wh.resolve(),
	}

	wh.HttpServer = &server

	return server.ListenAndServe()
}

func (wh *WebHook) resolve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		writeResponse := func(wp *payload.WebHook) error {
			if wp.Ec != 0 {
				w.WriteHeader(wp.Ec)
			}

			data, err := jsoniter.Marshal(wp)
			if err != nil {
				return err
			}

			_, err = w.Write(data)

			return err
		}

		runCallback := func(p *payload.WebHook, errs ...error) {
			if wh.callback == nil {
				return
			}

			filtered := make([]error, 0, len(errs))

			for _, e := range errs {
				if e != nil {
					filtered = append(filtered, e)
				}
			}

			wh.callback(p, filtered...)
		}

		p := &payload.WebHook{}
		p.RawRequest = r

		if r.URL.Path != wh.client.cfg.WebHookPath {
			go runCallback(p, fmt.Errorf("invalid path: %s", r.URL.Path))

			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusNotFound, Em: "Not found"}})

			return
		}

		if r.Method != http.MethodPost {
			go runCallback(p, fmt.Errorf("invalid method: %s", r.Method))

			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusMethodNotAllowed, Em: "Method not allowed"}})

			return
		}

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			go runCallback(p, err)

			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusInternalServerError, Em: "Internal server error"}})

			return
		}

		if len(raw) == 0 {
			go runCallback(p, errors.New("empty request body"))

			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusBadRequest, Em: "Bad request"}})

			return
		}

		r.Body = io.NopCloser(bytes.NewReader(raw))

		if err := jsoniter.Unmarshal(raw, p); err != nil {
			go runCallback(p, err)

			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusInternalServerError, Em: "Internal server error"}})

			return
		}

		if err := sign.WebHookSignVerify(p); err != nil {
			go runCallback(p, fmt.Errorf("invalid sign: %s", p.Data.Sign), err)

			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusBadRequest, Em: "Bad request"}})

			return
		}

		go runCallback(p, nil)

		_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusOK, Em: "OK"}})
	}
}
