package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func ReadAPIResponse(resp *http.Response) ([]byte, error) {
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resp.Body = io.NopCloser(bytes.NewReader(raw))

	if len(raw) == 0 {
		return nil, errors.New("empty response")
	}

	return raw, nil
}
