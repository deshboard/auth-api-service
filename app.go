package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type App struct {
	service *Service
}

// Create a new app
func NewApp() (*App, error) {
	rawAccessTokenExpiration := DefaultEnv("ACCESS_TOKEN_EXPIRATION", "15000")
	accessTokenExpiration, err := strconv.Atoi(rawAccessTokenExpiration)
	if err != nil {
		// TODO: add logging
		panic(err)
	}

	authenticator := NewUserServiceAuthenticator(
		RequireEnv("SERVICE_MODEL_USER"),
		DefaultEnv("ISSUER", "service.api.auth"),
		accessTokenExpiration,
		[]byte(RequireEnv("SIGNING_KEY")),
	)

	service := NewService(authenticator)
	service.getParams = func(r *http.Request) map[string]string {
		return mux.Vars(r)
	}

	return &App{
		service: service,
	}, nil
}

func (app *App) Shutdown() {
	// Noop
}

func (app *App) Listen() error {
	handler := app.CreateHandler()

	return http.ListenAndServe(":80", handler)
}

func (app *App) CreateHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/authenticate", app.service.Authenticate).Methods("POST")

	return router
}

func RequireEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Environment variable %s is mandatory", key))
	}

	return value
}

func DefaultEnv(key string, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	return value
}
