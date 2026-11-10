package utils

import (
	"bytes"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Sn0wo2/go-afdian-api/v2/pkg/payload"
)

type APIResponder interface {
	SetRawResponse(resp *http.Response)
	payload.Checker
}

type Response struct{}

func (Response) Process(resp *http.Response, p APIResponder) error {
	p.SetRawResponse(resp)

	rawBody, readErr := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	resp.Body = io.NopCloser(bytes.NewReader(rawBody))

	if readErr != nil {
		return readErr
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if len(rawBody) == 0 {
		return errors.New("empty response")
	}

	if err := json.Unmarshal(rawBody, p); err != nil {
		return err
	}

	if ec := p.GetEC(); ec != http.StatusOK {
		return fmt.Errorf("afdian api error: ec=%d, em=%s", ec, p.GetEM())
	}

	return nil
}
