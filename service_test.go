package main

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	authenticator := new(AuthenticatorMock)
	service := NewService(authenticator)

	assert.Equal(t, authenticator, service.authenticator)
}

func TestService_Authenticate(t *testing.T) {
	authenticator := new(AuthenticatorMock)
	service := NewService(authenticator)

	credentials := Credentials{
		"username": "user",
		"password": "pass",
	}

	tokens := &Tokens{
		AccessToken: &Token{
			Token:     "token",
			ExpiresIn: 12345,
		},
	}

	authenticator.On("Authenticate", credentials).Return(tokens, nil)

	b := []byte(`{"username": "user", "password": "pass"}`)

	req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader(b))
	w := httptest.NewRecorder()

	service.Authenticate(w, req)

	returnedTokens := new(Tokens)

	err := json.NewDecoder(w.Body).Decode(returnedTokens)
	require.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, tokens, returnedTokens)

	authenticator.AssertExpectations(t)
}

func TestService_Authenticate_Unauthorized(t *testing.T) {
	authenticator := new(AuthenticatorMock)
	service := NewService(authenticator)

	credentials := Credentials{
		"username": "invalid",
		"password": "pass",
	}

	authenticator.On("Authenticate", credentials).Return(nil, ErrUnauthorized)

	b := []byte(`{"username": "invalid", "password": "pass"}`)

	req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader(b))
	w := httptest.NewRecorder()

	service.Authenticate(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, 0, w.Body.Len())

	authenticator.AssertExpectations(t)
}

func TestService_Authenticate_Failure(t *testing.T) {
	authenticator := new(AuthenticatorMock)
	service := NewService(authenticator)

	credentials := Credentials{
		"username": "invalid",
		"password": "pass",
	}

	authenticator.On("Authenticate", credentials).Return(nil, errors.New("something went wrong"))

	b := []byte(`{"username": "invalid", "password": "pass"}`)

	req := httptest.NewRequest("POST", "/authenticate", bytes.NewReader(b))
	w := httptest.NewRecorder()

	service.Authenticate(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, 0, w.Body.Len())

	authenticator.AssertExpectations(t)
}
