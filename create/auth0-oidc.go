package create

import (
	"clerk-portal/shared/env"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Get new *Authenticator.
func GetOidcAuthenticator(envprefix string) (*Auth0Oidc, error) {
	e, err := env.NewOIDCEnvConfig(envprefix)
	if err != nil {
		return nil, err
	}
	provider, err := oidc.NewProvider(context.Background(), "https://"+e.GetDomain()+"/")

	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     e.GetClientID(),
		ClientSecret: e.GetClientSecret(),
		RedirectURL:  e.GetCallbackURL(),
		Endpoint:     provider.Endpoint(),
		Scopes:       e.GetScopes(),
	}

	authObj := &Auth0Oidc{
		Provider: provider,
		Config:   conf,
	}

	return authObj, nil
}

// Auth0Oidc is used to authenticate our users.
type Auth0Oidc struct {
	*oidc.Provider
	oauth2.Config
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Auth0Oidc) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

func (a *Auth0Oidc) GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
