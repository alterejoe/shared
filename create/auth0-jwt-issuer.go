package create

import (
	"github.com/alterejoe/shared/env"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type JwtIssuer struct {
	Domain, ClientID, ClientSecret, Audience, Scope string
	token                                           string
	expiry                                          time.Time
	mu                                              sync.Mutex
}

type Token struct {
	Token string `json:"token"`
}

func GetJwtIssuer(targetjwtprefix string) (*JwtIssuer, error) {
	e, err := env.NewJwtIssuerENV(targetjwtprefix)
	if err != nil {
		return nil, err
	}

	return &JwtIssuer{
		Domain:       e.GetDomain(),
		ClientID:     e.GetClientID(),
		ClientSecret: e.GetClientSecret(),
		Audience:     e.GetAudience(),
		Scope:        e.GetScope(),
	}, nil
}

func (a *JwtIssuer) validateFields() error {
	if a.Domain == "" || a.ClientID == "" || a.ClientSecret == "" || a.Audience == "" {
		return fmt.Errorf("missing Auth0 params: domain=%q, client_id=%q, client_secret=%q, audience=%q, scope=%q",
			a.Domain, a.ClientID, a.ClientSecret, a.Audience, a.Scope)
	}
	return nil
}

func (a *JwtIssuer) GetJwtToken() (string, error) {
	if a.token != "" && time.Until(a.expiry) > 30*time.Second {
		return a.token, nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.token != "" && time.Until(a.expiry) > 30*time.Second {
		return a.token, nil
	}

	if err := a.validateFields(); err != nil {
		return "", err
	}

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {a.ClientID},
		"client_secret": {a.ClientSecret},
		"audience":      {a.Audience},
	}
	if a.Scope != "" {
		form.Set("scope", a.Scope)
	}

	req, err := http.NewRequest("POST", "https://"+a.Domain+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token: %s", strings.TrimSpace(string(body)))
	}

	var resp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", errors.New("invalid JSON response from Auth0")
	}

	a.token = resp.AccessToken
	a.expiry = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
	return a.token, nil
}
