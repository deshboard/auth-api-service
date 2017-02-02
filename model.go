package main

import "errors"

type Tokens struct {
	AccessToken *Token `json:"access_token"`
}

type Token struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

type Credentials map[string]string

var ErrUnauthorized = errors.New("auth: unauthorized")

type Authenticator interface {
	Authenticate(credentials Credentials) (*Tokens, error)
}
