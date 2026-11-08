package afdian

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json/v2"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
)

type CallBack func(p *payload.WebHook, errs ...error)

type WebHook struct {
	client     *Client
	HTTPServer *http.Server
}

var webhookPublicKey = func() *rsa.PublicKey {
	block, _ := pem.Decode([]byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwwdaCg1Bt+UKZKs0R54y
lYnuANma49IpgoOwNmk3a0rhg/PQuhUJ0EOZSowIC44l0K3+fqGns3Ygi4AfmEfS
4EKbdk1ahSxu7Zkp2rHMt+R9GarQFQkwSS/5x1dYiHNVMiR8oIXDgjmvxuNes2Cr
8fw9dEF0xNBKdkKgG2qAawcN1nZrdyaKWtPVT9m2Hl0ddOO9thZmVLFOb9NVzgYf
jEgI+KWX6aY19Ka/ghv/L4t1IXmz9pctablN5S0CRWpJW3Cn0k6zSXgjVdKm4uN7
jRlgSRaf/Ind46vMCm3N2sgwxu/g3bnooW+db0iLo13zzuvyn727Q3UDQ0MmZcEW
MQIDAQAB
-----END PUBLIC KEY-----`))
	if block == nil {
		panic("afdian: invalid embedded webhook public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("afdian: invalid embedded webhook public key")
	}

	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		panic("afdian: invalid embedded webhook public key")
	}

	return pub
}()

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
		Handler:           wh.resolve(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return wh.HTTPServer.ListenAndServe()
}

func (wh *WebHook) resolve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		p := &payload.WebHook{RawRequest: r}

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

		if err := json.Unmarshal(raw, p); err != nil {
			handleErr(err, http.StatusBadRequest, "Bad request")

			return
		}

		if p.Data == nil || p.Data.Order == nil {
			handleErr(errors.New("invalid webhook payload"), http.StatusBadRequest, "Bad request")

			return
		}

		sig, err := base64.StdEncoding.DecodeString(p.Data.Sign)
		if err == nil {
			hashed := crypto.SHA256.New()
			hashed.Write([]byte(p.Data.Order.OutTradeNo +
				p.Data.Order.UserID +
				p.Data.Order.PlanID +
				p.Data.Order.TotalAmount))

			err = rsa.VerifyPKCS1v15(webhookPublicKey, crypto.SHA256, hashed.Sum(nil), sig)
		}

		if err != nil {
			handleErr(err, http.StatusBadRequest, "Bad request")

			return
		}

		go wh.runCallback(p)

		_ = wh.writeResponse(w, &payload.WebHook{Base: payload.Base{EC: http.StatusOK, EM: "OK"}})
	}
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

func (wh *WebHook) writeResponse(w http.ResponseWriter, writePayload *payload.WebHook) error {
	if writePayload.EC != 0 {
		w.WriteHeader(writePayload.EC)
	}

	data, err := json.Marshal(writePayload)
	if err != nil {
		return err
	}

	_, err = w.Write(data)

	return err
}
