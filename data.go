package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceAuthenticator struct {
	endpoint              string
	issuer                string
	accessTokenExpiration int
	signingKey            []byte
}

type ServiceUser struct {
	Id                int    `json:"id"`
	Username          string `json:"username"`
	EncryptedPassword string `json:"encrypted_password"`
}

func NewUserServiceAuthenticator(
	endpoint string,
	issuer string,
	accessTokenExpiration int,
	signingKey []byte,
) *UserServiceAuthenticator {
	return &UserServiceAuthenticator{
		endpoint:              endpoint,
		issuer:                issuer,
		accessTokenExpiration: accessTokenExpiration,
		signingKey:            signingKey,
	}
}

func (a *UserServiceAuthenticator) Authenticate(credentials Credentials) (*Tokens, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/%s", a.endpoint, credentials["username"]))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrUnauthorized
	}

	user := new(ServiceUser)

	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(credentials["password"]))
	if err != nil {
		return nil, ErrUnauthorized
	}

	now := time.Now()

	claims := &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Duration(a.accessTokenExpiration)).Unix(),
		Issuer:    a.issuer,
		Subject:   string(user.Id),
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwt.SignedString(a.signingKey)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken: &Token{
			Token:     tokenString,
			ExpiresIn: a.accessTokenExpiration,
		},
	}, nil
}
