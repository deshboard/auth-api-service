package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserServiceAuthenticator(t *testing.T) {
	signingKey := []byte("signing_key")

	authenticator := NewUserServiceAuthenticator(
		"endpoint",
		"issuer",
		1234,
		signingKey,
	)

	assert.Equal(t, "endpoint", authenticator.endpoint)
	assert.Equal(t, "issuer", authenticator.issuer)
	assert.Equal(t, 1234, authenticator.accessTokenExpiration)
	assert.Equal(t, signingKey, authenticator.signingKey)
}

func TestUserServiceAuthenticator_Authenticate_Integration(t *testing.T) {
	if !*integration {
		t.Skip()
	}

	signingKey := []byte("signing_key")

	authenticator := NewUserServiceAuthenticator(
		RequireEnv("SERVICE_MODEL_USER"),
		"issuer",
		1234,
		signingKey,
	)

	b := `{"username": "user", "password": "password"}`

	resp, err := http.Post(
		fmt.Sprintf("http://%s", authenticator.endpoint),
		"application/json",
		strings.NewReader(b),
	)

	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 201, resp.StatusCode)

	credentials := Credentials{
		"username": "user",
		"password": "password",
	}

	tokens, err := authenticator.Authenticate(credentials)

	assert.NoError(t, err)

	assert.Equal(t, 1234, tokens.AccessToken.ExpiresIn)

	token, err := jwt.Parse(tokens.AccessToken.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	assert.True(t, ok)
	assert.True(t, token.Valid)
	assert.Equal(t, "issuer", claims["iss"])
}

func TestUserServiceAuthenticator_Authenticate_UserNotFound_Integration(t *testing.T) {
	if !*integration {
		t.Skip()
	}

	signingKey := []byte("signing_key")

	authenticator := NewUserServiceAuthenticator(
		RequireEnv("SERVICE_MODEL_USER"),
		"issuer",
		1234,
		signingKey,
	)

	credentials := Credentials{
		"username": "user_not_found",
		"password": "password",
	}

	tokens, err := authenticator.Authenticate(credentials)

	assert.Error(t, err)

	assert.Nil(t, tokens)
	assert.Equal(t, ErrUnauthorized, err)
}

func TestUserServiceAuthenticator_Authenticate_WrongPassword_Integration(t *testing.T) {
	if !*integration {
		t.Skip()
	}

	signingKey := []byte("signing_key")

	authenticator := NewUserServiceAuthenticator(
		RequireEnv("SERVICE_MODEL_USER"),
		"issuer",
		1234,
		signingKey,
	)

	b := `{"username": "user_wrong_password", "password": "password"}`

	resp, err := http.Post(
		fmt.Sprintf("http://%s", authenticator.endpoint),
		"application/json",
		strings.NewReader(b),
	)

	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, 201, resp.StatusCode)

	credentials := Credentials{
		"username": "user_wrong_password",
		"password": "wrong_password",
	}

	tokens, err := authenticator.Authenticate(credentials)

	assert.Error(t, err)

	assert.Nil(t, tokens)
	assert.Equal(t, ErrUnauthorized, err)
}
