package client

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dzt-corp/go-etsy/oauth"
)

const (
	defaultBaseURL = "https://api.etsy.com/v3/application/"
	userAgent      = "go-etsy"
	expiryDelta    = 1 * time.Minute
)

type EtsyClient struct {
	accessToken       string
	accessTokenExpiry time.Time
	cfg               *Config
}

type Config struct {
	APIKey       string
	RefreshToken string
	OAuth        *oauth.OAuthClient
}

func (o Config) IsValid() (bool, error) {
	if o.RefreshToken == "" {
		return false, errors.New("refresh token is required")
	}
	if o.APIKey == "" {
		return false, errors.New("API key is required")
	}
	if o.OAuth == nil {
		return false, errors.New("oauth2 client is required")
	}
	return true, nil
}

func NewEtsyClient(cfg *Config) (*EtsyClient, error) {
	if isValid, err := cfg.IsValid(); !isValid {
		return nil, err
	}

	client := &EtsyClient{}
	client.cfg = cfg
	return client, nil
}

func (etsy *EtsyClient) AuthorizeRequest(r *http.Request) error {
	if etsy.accessToken == "" ||
		etsy.accessTokenExpiry.IsZero() ||
		etsy.accessTokenExpiry.Round(0).Add(-expiryDelta).Before(time.Now().UTC()) {
		if err := etsy.RefreshToken(); err != nil {
			return fmt.Errorf("cannot refresh token. Error: %s", err.Error())
		}
	}

	r.Header.Add("Authorization", "Bearer "+etsy.accessToken)

	return nil
}

func (etsy *EtsyClient) RefreshToken() error {
	resp, err := etsy.cfg.OAuth.RefreshToken(etsy.cfg.RefreshToken)
	if err != nil {
		return err
	}

	etsy.accessToken = resp.AccessToken
	etsy.accessTokenExpiry = time.Now().UTC().Add(time.Duration(resp.ExpiresIn) * time.Second) //set expiration time
	return nil
}
