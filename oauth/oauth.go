package oauth

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type OAuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	HTTPClient   *http.Client
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func NewOAuthClient(clientID, redirectURI string) *OAuthClient {
	return &OAuthClient{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		HTTPClient:  http.DefaultClient,
	}
}

// Generate a secure random code_verifier
func GenerateCodeVerifier() (string, error) {
	bytes := make([]byte, 32) // 32 bytes => ~43 base64 characters
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func (c *OAuthClient) Connect(codeVerifier string, scopes []string) string {
	baseURL := "https://www.etsy.com/oauth/connect"

	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal("Invalid base URL:", err)
	}

	codeChallenge := GenerateCodeChallenge(codeVerifier)
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", c.ClientID)
	q.Set("redirect_uri", c.RedirectURI)
	q.Set("scope", strings.Join(scopes, " "))
	q.Set("state", "superstate") //state=superstate assigns the state string to superstate which Etsy.com should return with the authorization code
	q.Set("code_challenge", codeChallenge)
	q.Set("code_challenge_method", "S256")

	u.RawQuery = q.Encode()
	return u.String()
}

func (c *OAuthClient) ExchangeCode(code, codeVerifier string) (*AccessTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", c.ClientID)
	form.Set("redirect_uri", c.RedirectURI)
	form.Set("code", code)
	form.Set("code_verifier", codeVerifier)

	req, err := http.NewRequest("POST", "https://api.etsy.com/v3/public/oauth/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (c *OAuthClient) RefreshToken(refreshToken string) (*AccessTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	form.Set("client_id", c.ClientID)

	req, err := http.NewRequest("POST", "https://api.etsy.com/v3/public/oauth/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}
