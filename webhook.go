package afdian

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
	"github.com/json-iterator/go"
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

// SetCallback sets the callback function, must implement idempotent logic
func (wh *WebHook) SetCallback(callback CallBack) {
	wh.callback = callback
}

func (wh *WebHook) Start() error {
	if wh.client.cfg.WebHookListenAddr == "" {
		return nil
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

		if r.URL.Path != wh.client.cfg.WebHookPath {
			go runCallback(nil, fmt.Errorf("invalid path: %s", r.URL.Path))
			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusNotFound, Em: "Not found"}})
			return
		}

		if r.URL.Query().Get("auth") != wh.client.cfg.WebHookQueryToken {
			go runCallback(nil, fmt.Errorf("invalid token: %s", r.URL.Query().Get("auth")))
			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusUnauthorized, Em: "Unauthorized"}})
			return
		}

		if r.Method != http.MethodPost {
			go runCallback(nil, fmt.Errorf("invalid method: %s", r.Method))
			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusMethodNotAllowed, Em: "Method not allowed"}})
			return
		}

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			go runCallback(nil, err)
			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusInternalServerError, Em: "Internal server error"}})
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(raw))
		if len(raw) == 0 {
			go runCallback(nil, fmt.Errorf("empty request body"))
			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusBadRequest, Em: "Bad request"}})
			return
		}

		p := &payload.WebHook{}
		if err := jsoniter.Unmarshal(raw, p); err != nil {
			go runCallback(nil, err)
			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusInternalServerError, Em: "Internal server error"}})
			return
		}

		if WebHookSignVerify(p) {
			go runCallback(nil, fmt.Errorf("invalid sign: %s", p.Data.Sign))
			_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusBadRequest, Em: "Bad request"}})
			return
		}

		p.RawRequest = r

		go runCallback(p, nil)
		_ = writeResponse(&payload.WebHook{Base: payload.Base{Ec: http.StatusOK, Em: "OK"}})
	}
}
