package main

import (
	"net/http"
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

// Handles application shutdown (closes DB connection, etc)
// Make sure the process does not exit before this is called
func (app *App) Shutdown() {
	// Noop
}

// Starts listening
func (app *App) Listen() error {
	handler := app.CreateHandler()

	return http.ListenAndServe(":80", handler)
}

// Creates and configures the router
func (app *App) CreateHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/authenticate", app.service.Authenticate).Methods("POST")

	router.HandleFunc("/_status/healthz", app.HealthStatus).Methods("GET")
	router.HandleFunc("/_status/readiness", app.ReadinessStatus).Methods("GET")

	return router
}

// Checks if the app is up and running
func (app *App) HealthStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Checks if the app is ready for accepting request (eg. database is available as well)
func (app *App) ReadinessStatus(w http.ResponseWriter, r *http.Request) {
	// Todo add check for user model service

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
