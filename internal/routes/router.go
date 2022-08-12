package routes

import (
	"net/http"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/adamstrickland/dapper-api/internal/logins"
	"github.com/adamstrickland/dapper-api/internal/signups"
	"github.com/adamstrickland/dapper-api/internal/users"
	"github.com/gorilla/mux"
)

func NewRouter(cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/signup", signups.NewPostHandler(cfg)).
		Methods(http.MethodPost)

	router.HandleFunc("/login", logins.NewPostHandler(cfg)).
		Methods(http.MethodPost)

	router.Use(LoggingMiddleware(cfg))

	router.Use(ContentTypeMiddleware(cfg))

	srouter := router.
		Name("secured").
		Subrouter()

	srouter.HandleFunc("/users", users.NewGetHandler(cfg)).
		Methods(http.MethodGet)

	srouter.HandleFunc("/users", users.NewPutHandler(cfg)).
		Methods(http.MethodPut)

	srouter.Use(AuthnMiddleware(cfg))

	return router
}
