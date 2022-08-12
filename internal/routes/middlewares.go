package routes

import (
	"log"
	"net/http"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/adamstrickland/dapper-api/internal/security"
)

func LoggingMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("RECV: %-8s %-25s %8d", r.Method, r.URL, r.ContentLength)
			next.ServeHTTP(w, r)
		})
	}
}

func ContentTypeMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, "", http.StatusNotFound)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func AuthnMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token := r.Header.Get(cfg.GetString("tokenHeader")); token != "" {
				if ok, err := security.IsValidToken(cfg, token); ok && err == nil {
					next.ServeHTTP(w, r)
				} else {
					http.Error(w, "", http.StatusUnauthorized)
				}
			} else {
				http.Error(w, "", http.StatusNotFound)
			}
		})
	}
}
