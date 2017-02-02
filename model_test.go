package main

import "github.com/stretchr/testify/mock"

type AuthenticatorMock struct {
	mock.Mock
}

func (a *AuthenticatorMock) Authenticate(credentials Credentials) (*Tokens, error) {
	args := a.Called(credentials)

	tokens := args.Get(0)
	if tokens == nil {
		return nil, args.Error(1)
	}

	return tokens.(*Tokens), args.Error(1)
}
