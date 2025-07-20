package afdian

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/sign"
	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
	jsoniter "github.com/json-iterator/go"
)

// CallBack for WebHook.
type CallBack func(p *payload.WebHook, errs ...error)

// WebHook handler
type WebHook struct {
	client     *Client
	HTTPServer *http.Server
}

// NewWebHook creates a new WebHook handler.
func NewWebHook(client *Client) *WebHook {
	wh := &WebHook{client: client}
	client.WebHook = wh

	return wh
}

// Start the WebHook server
func (wh *WebHook) Start() error {
	if wh.client.cfg.WebHookListenAddr == "" {
		return errors.New("WebHookListenAddr is empty")
	}

	server := http.Server{
		Addr:              wh.client.cfg.WebHookListenAddr,
		Handler:           wh.resolve(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	wh.HTTPServer = &server

	return server.ListenAndServe()
}

func (wh *WebHook) resolve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		p := &payload.WebHook{}
		p.RawRequest = r

		handleErr := func(err error, code int, msg string) {
			go wh.runCallback(p, err)

			_ = wh.writeResponse(w, &payload.WebHook{Base: payload.Base{EC: code, EM: msg}})
		}

		if r.URL.Path != wh.client.cfg.WebHookPath {
			handleErr(fmt.Errorf("invalid path: %s", r.URL.Path), http.StatusNotFound, "Not found")

			return
		}

		if r.Method != http.MethodPost {
			handleErr(fmt.Errorf("invalid method: %s", r.Method), http.StatusMethodNotAllowed, "Method not allowed")

			return
		}

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			handleErr(err, http.StatusInternalServerError, "Internal server error")

			return
		}

		if len(raw) == 0 {
			handleErr(errors.New("empty request body"), http.StatusBadRequest, "Bad request")

			return
		}

		r.Body = io.NopCloser(bytes.NewReader(raw))

		if err := jsoniter.Unmarshal(raw, p); err != nil {
			handleErr(err, http.StatusBadRequest, "Bad request")

			return
		}

		if err := sign.WebHookSignVerify(p); err != nil {
			handleErr(fmt.Errorf("invalid sign: %s", p.Data.Sign), http.StatusBadRequest, "Bad request")

			return
		}

		go wh.runCallback(p, nil)

		_ = wh.writeResponse(w, &payload.WebHook{Base: payload.Base{EC: http.StatusOK, EM: "OK"}})
	}
}

func (wh *WebHook) runCallback(p *payload.WebHook, errs ...error) {
	if wh.client.cfg.WebHookCallback == nil {
		return
	}

	filtered := make([]error, 0, len(errs))

	for _, e := range errs {
		if e != nil {
			filtered = append(filtered, e)
		}
	}

	wh.client.cfg.WebHookCallback(p, filtered...)
}

func (wh *WebHook) writeResponse(w http.ResponseWriter, writePayload *payload.WebHook) error {
	if writePayload.EC != 0 {
		w.WriteHeader(writePayload.EC)
	}

	data, err := jsoniter.Marshal(writePayload)
	if err != nil {
		return err
	}

	_, err = w.Write(data)

	return err
}
