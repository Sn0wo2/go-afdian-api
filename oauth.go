package afdian

import (
	"cmp"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/Sn0wo2/go-afdian-api/internal/utils"
	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
)

func (c *Client) AuthorizeURL(redirectURI, state string) (string, error) {
	if c.cfg.OAuthClientID == "" {
		return "", errors.New("afdian oauth: client_id is required, set Config.OAuthClientID")
	}

	redirectURI = cmp.Or(redirectURI, c.cfg.OAuthRedirectURI)
	if redirectURI == "" {
		return "", errors.New("afdian oauth: redirect_uri is required")
	}

	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("scope", cmp.Or(c.cfg.OAuthScope, "basic"))
	q.Set("client_id", c.cfg.OAuthClientID)
	q.Set("redirect_uri", redirectURI)

	if state != "" {
		q.Set("state", state)
	}

	return strings.TrimSuffix(c.cfg.BaseURL, "/api") + "/oauth2/authorize?" + q.Encode(), nil
}

func (c *Client) ExchangeToken(code, redirectURI string) (*payload.OAuthAccessToken, error) {
	if c.cfg.OAuthClientID == "" || c.cfg.OAuthClientSecret == "" {
		return nil, errors.New("afdian oauth: client_id and client_secret are required, set Config.OAuthClientID and Config.OAuthClientSecret")
	}

	if code == "" {
		return nil, errors.New("afdian oauth: code is required")
	}

	redirectURI = cmp.Or(redirectURI, c.cfg.OAuthRedirectURI)
	if redirectURI == "" {
		return nil, errors.New("afdian oauth: redirect_uri is required")
	}

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", c.cfg.OAuthClientID)
	form.Set("client_secret", c.cfg.OAuthClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest(http.MethodPost, c.cfg.BaseURL+"/oauth2/access_token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	resp, err := c.HTTP.Do(req) //nolint:bodyclose // false positive: utils.Response.Process closes the body
	if err != nil {
		return nil, err
	}

	p := new(payload.OAuthAccessToken)

	return p, (utils.Response{}).Process(resp, p)
}
