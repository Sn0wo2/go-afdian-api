package afdian

import (
	"bytes"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Sn0wo2/go-afdian-api/internal/webhook"
	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
)

type CallBack func(p *payload.WebHook, errs ...error)

type WebHook struct {
	client     *Client
	HTTPServer *http.Server
}

func NewWebHook(client *Client) *WebHook {
	client.WebHook = &WebHook{client: client}

	return client.WebHook
}

func (wh *WebHook) Start() error {
	if wh.client.cfg.WebHookListenAddr == "" {
		return errors.New("WebHookListenAddr is empty")
	}

	wh.HTTPServer = &http.Server{
		Addr:              wh.client.cfg.WebHookListenAddr,
		Handler:           http.HandlerFunc(wh.handle),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return wh.HTTPServer.ListenAndServe()
}

func (wh *WebHook) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	p := &payload.WebHook{RawRequest: r}

	handleErr := func(err error, status, ec int, msg string) {
		go wh.runCallback(p, err)

		_ = wh.writeResponse(w, status, &payload.WebHook{Base: payload.Base{EC: ec, EM: msg}})
	}

	if r.URL.Path != wh.client.cfg.WebHookPath {
		handleErr(fmt.Errorf("invalid path: %s", r.URL.Path), http.StatusNotFound, http.StatusNotFound, "Not found")

		return
	}

	if r.Method != http.MethodPost {
		handleErr(fmt.Errorf("invalid method: %s", r.Method), http.StatusMethodNotAllowed, http.StatusMethodNotAllowed, "Method not allowed")

		return
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		handleErr(err, http.StatusInternalServerError, http.StatusInternalServerError, "Internal server error")

		return
	}

	if len(raw) == 0 {
		handleErr(errors.New("empty request body"), http.StatusBadRequest, http.StatusBadRequest, "Bad request")

		return
	}

	r.Body = io.NopCloser(bytes.NewReader(raw))

	if err := json.Unmarshal(raw, p); err != nil {
		handleErr(err, http.StatusBadRequest, http.StatusBadRequest, "Bad request")

		return
	}

	if err := webhook.Verify(p); err != nil {
		handleErr(err, http.StatusOK, http.StatusBadRequest, "signature verification failed")

		return
	}

	go wh.runCallback(p)

	_ = wh.writeResponse(w, http.StatusOK, &payload.WebHook{Base: payload.Base{EC: http.StatusOK, EM: "OK"}})
}

func (wh *WebHook) runCallback(p *payload.WebHook, errs ...error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("afdian: recovered from webhook callback panic: %v", r)
		}
	}()

	if cb := wh.client.cfg.WebHookCallback; cb != nil {
		cb(p, errs...)
	}
}

func (wh *WebHook) writeResponse(w http.ResponseWriter, status int, writePayload *payload.WebHook) error {
	w.WriteHeader(status)

	data, err := json.Marshal(writePayload)
	if err != nil {
		return err
	}

	_, err = w.Write(data)

	return err
}
