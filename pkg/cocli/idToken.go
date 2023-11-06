package cocli

import "golang.org/x/oauth2"

// These types & fns are used to facilitate id_token use
type TokenWithIdToken struct {
	*oauth2.Token
	IDToken string `json:"id_token,omitempty"`
}

type TokenWithIdTokenSource struct {
	oauth2.TokenSource
}

func (tidts *TokenWithIdTokenSource) Token() (*TokenWithIdToken, error) {
	token, err := tidts.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	// This means token was not refreshed
	if token.Extra("id_token") == nil {
		return GetIdTokenFromCredsFile(), nil
	}

	tokenWithIdToken := &TokenWithIdToken{
		Token:   token,
		IDToken: token.Extra("id_token").(string),
	}

	return tokenWithIdToken, nil
}
