package afdian

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/Sn0wo2/go-afdian-api/internal/utils"
	"github.com/parquet-go/bitpack/unsafecast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper struct {
	roundTrip func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTrip(req)
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("default http client", func(t *testing.T) {
		t.Parallel()

		client := NewClient(&Config{
			UserID:   "test_user_id",
			APIToken: "test_api_token",
		})
		assert.Equal(t, http.DefaultClient, client.HTTP)
	})

	t.Run("custom http client", func(t *testing.T) {
		t.Parallel()

		customClient := &http.Client{}
		client := NewClient(&Config{
			UserID:   "test_user_id",
			APIToken: "test_api_token",
		}, customClient)
		assert.Equal(t, customClient, client.HTTP)
	})
}

func TestSend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		params         map[string]string
		mockResponse   *http.Response
		mockError      error
		expectedBody   string
		expectError    bool
		expectedStatus int
	}{
		{
			name: "successful request",
			path: "/open/test",
			params: map[string]string{
				"key": "value",
			},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"ec":200,"em":"ok"}`)),
				Header:     make(http.Header),
			},
			mockError:      nil,
			expectedBody:   `{"ec":200,"em":"ok"}`,
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:        "new request error",
			path:        "\n",
			params:      nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &Config{
				UserID:   "test_user_id",
				APIToken: "test_api_token",
			}

			mockClient := &http.Client{
				Transport: &mockRoundTripper{
					roundTrip: func(req *http.Request) (*http.Response, error) {
						assert.Equal(t, http.MethodPost, req.Method)
						assert.Equal(t, "https://afdian.com/api"+tt.path, req.URL.String())
						assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
						assert.Equal(t, utils.UserAgent{}.String(), req.Header.Get("User-Agent"))

						return tt.mockResponse, tt.mockError
					},
				},
			}

			client := NewClient(cfg, mockClient)
			resp, err := client.Send(tt.path, tt.params)

			if tt.expectError {
				assert.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, readErr := io.ReadAll(resp.Body)
			require.NoError(t, readErr)
			assert.Equal(t, tt.expectedBody, unsafecast.String(body))
			assert.NoError(t, resp.Body.Close())
		})
	}
}
