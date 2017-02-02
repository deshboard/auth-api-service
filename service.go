package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Returns parameters from the request
// (decouples the service from the router implementation)
type ParamFetcher func(r *http.Request) map[string]string

type Service struct {
	authenticator Authenticator
	getParams     ParamFetcher
}

func NewService(authenticator Authenticator) *Service {
	return &Service{
		authenticator: authenticator,
		getParams: func(r *http.Request) map[string]string {
			return make(map[string]string)
		},
	}
}

func (s *Service) Authenticate(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		// TODO: add logging
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	// TODO: add validation
	tokens, err := s.authenticator.Authenticate(credentials)
	if err != nil {
		switch err {
		case ErrUnauthorized:
			w.WriteHeader(http.StatusUnauthorized)

			return
		default:
			// TODO: add logging
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	var body bytes.Buffer

	err = json.NewEncoder(&body).Encode(tokens)
	if err != nil {
		// TODO: add logging
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body.WriteTo(w)
}
