package afdian

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorizeURL(t *testing.T) {
	t.Parallel()

	t.Run("default scope and config redirect uri", func(t *testing.T) {
		t.Parallel()

		got, err := NewClient(&Config{
			OAuthClientID:    "client123",
			OAuthRedirectURI: "https://example.com/callback",
		}).AuthorizeURL("", "state111")

		require.NoError(t, err)

		u, err := url.Parse(got)
		require.NoError(t, err)

		assert.Equal(t, "https://afdian.com/oauth2/authorize", u.Scheme+"://"+u.Host+u.Path)

		q := u.Query()
		assert.Equal(t, "code", q.Get("response_type"))
		assert.Equal(t, "basic", q.Get("scope"))
		assert.Equal(t, "client123", q.Get("client_id"))
		assert.Equal(t, "https://example.com/callback", q.Get("redirect_uri"))
		assert.Equal(t, "state111", q.Get("state"))
	})

	t.Run("param redirect uri overrides config and custom scope", func(t *testing.T) {
		t.Parallel()

		got, err := NewClient(&Config{
			BaseURL:          "https://ifdian.net/api",
			OAuthClientID:    "client123",
			OAuthScope:       "basic",
			OAuthRedirectURI: "https://example.com/config",
		}).AuthorizeURL("https://example.com/callback", "")

		require.NoError(t, err)

		assert.True(t, strings.HasPrefix(got, "https://ifdian.net/oauth2/authorize?"), got)

		q, err := url.ParseQuery(strings.TrimPrefix(got, "https://ifdian.net/oauth2/authorize?"))
		require.NoError(t, err)

		assert.Equal(t, "https://example.com/callback", q.Get("redirect_uri"))
		assert.Empty(t, q.Get("state"))
	})

	t.Run("missing client id", func(t *testing.T) {
		t.Parallel()

		_, err := NewClient(&Config{}).AuthorizeURL("https://example.com/callback", "")

		require.Error(t, err)
		assert.EqualError(t, err, "afdian oauth: client_id is required, set Config.OAuthClientID")
	})

	t.Run("missing redirect uri", func(t *testing.T) {
		t.Parallel()

		_, err := NewClient(&Config{OAuthClientID: "client123"}).AuthorizeURL("", "")

		require.Error(t, err)
		assert.EqualError(t, err, "afdian oauth: redirect_uri is required")
	})
}

func TestExchangeToken(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		token, err := NewClient(&Config{
			OAuthClientID:     "client123",
			OAuthClientSecret: "secret123",
			OAuthRedirectURI:  "https://example.com/callback",
		}, &http.Client{
			Transport: &mockRoundTripper{
				roundTrip: func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodPost, req.Method)
					assert.Equal(t, "https://afdian.com/api/oauth2/access_token", req.URL.String())
					assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
					assert.Equal(t, userAgent(), req.Header.Get("User-Agent"))

					body, err := io.ReadAll(req.Body)
					require.NoError(t, err)

					form, err := url.ParseQuery(string(body))
					require.NoError(t, err)

					assert.Equal(t, "authorization_code", form.Get("grant_type"))
					assert.Equal(t, "client123", form.Get("client_id"))
					assert.Equal(t, "secret123", form.Get("client_secret"))
					assert.Equal(t, "the-code", form.Get("code"))
					assert.Equal(t, "https://example.com/callback", form.Get("redirect_uri"))

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{"ec":200,"em":"ok","data":{"user_id":"adf397fe8374811eaacee52540025c377","user_private_id":"33","name":"SomeOne","avatar":"https://pic1.afdiancdn.com/user/avatar.png"}}`)),
						Header:     make(http.Header),
					}, nil
				},
			},
		}).ExchangeToken("the-code", "")

		require.NoError(t, err)
		require.NotNil(t, token)
		assert.Equal(t, 200, token.EC)
		assert.Equal(t, "ok", token.EM)
		require.NotNil(t, token.Data)
		assert.Equal(t, "adf397fe8374811eaacee52540025c377", token.Data.UserID)
		assert.Equal(t, "33", token.Data.UserPrivateID)
		assert.Equal(t, "SomeOne", token.Data.Name)
		assert.Equal(t, "https://pic1.afdiancdn.com/user/avatar.png", token.Data.Avatar)
		assert.NotNil(t, token.RawResponse)
	})

	t.Run("param redirect uri overrides config", func(t *testing.T) {
		t.Parallel()

		token, err := NewClient(&Config{
			OAuthClientID:     "client123",
			OAuthClientSecret: "secret123",
			OAuthRedirectURI:  "https://example.com/config",
		}, &http.Client{
			Transport: &mockRoundTripper{
				roundTrip: func(req *http.Request) (*http.Response, error) {
					body, err := io.ReadAll(req.Body)
					require.NoError(t, err)

					form, err := url.ParseQuery(string(body))
					require.NoError(t, err)

					assert.Equal(t, "https://example.com/param", form.Get("redirect_uri"))

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{"ec":200,"em":"ok","data":{"user_id":"user1"}}`)),
						Header:     make(http.Header),
					}, nil
				},
			},
		}).ExchangeToken("the-code", "https://example.com/param")

		require.NoError(t, err)
		require.NotNil(t, token)
		assert.Equal(t, "user1", token.Data.UserID)
	})

	t.Run("api error", func(t *testing.T) {
		t.Parallel()

		token, err := NewClient(&Config{
			OAuthClientID:     "client123",
			OAuthClientSecret: "secret123",
			OAuthRedirectURI:  "https://example.com/callback",
		}, &http.Client{Transport: &mockRoundTripper{
			roundTrip: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"ec":400,"em":"invalid code"}`)),
					Header:     make(http.Header),
				}, nil
			},
		}}).ExchangeToken("the-code", "")

		require.Error(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, 400, token.EC)
		assert.Equal(t, "invalid code", token.EM)
		assert.EqualError(t, err, "afdian api error: ec=400, em=invalid code")
	})

	t.Run("missing client credentials", func(t *testing.T) {
		t.Parallel()

		_, err := NewClient(&Config{OAuthClientID: "client123"}).ExchangeToken("the-code", "https://example.com/callback")

		require.Error(t, err)
		assert.EqualError(t, err, "afdian oauth: client_id and client_secret are required, set Config.OAuthClientID and Config.OAuthClientSecret")
	})

	t.Run("missing code", func(t *testing.T) {
		t.Parallel()

		_, err := NewClient(&Config{
			OAuthClientID:     "client123",
			OAuthClientSecret: "secret123",
		}).ExchangeToken("", "https://example.com/callback")

		require.Error(t, err)
		assert.EqualError(t, err, "afdian oauth: code is required")
	})

	t.Run("missing redirect uri", func(t *testing.T) {
		t.Parallel()

		_, err := NewClient(&Config{
			OAuthClientID:     "client123",
			OAuthClientSecret: "secret123",
		}).ExchangeToken("the-code", "")

		require.Error(t, err)
		assert.EqualError(t, err, "afdian oauth: redirect_uri is required")
	})
}
