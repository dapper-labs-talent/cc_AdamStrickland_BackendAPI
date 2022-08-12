package logins

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/adamstrickland/dapper-api/internal/security"
	"github.com/adamstrickland/dapper-api/internal/users"
)

type requestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewPostHandler(cfg *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var qp requestPayload

		err := json.NewDecoder(r.Body).Decode(&qp)

		if err != nil {
			log.Printf("Unable to unmarshal payload: %e", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := users.FindByEmail(cfg, qp.Email)

		if err != nil {
			log.Printf("Unable to identify user: %e", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.UnencryptedPassword != qp.Password {
			log.Println("Unable to authenticate password!")
			http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
			return
		}

		data, err := security.NewTokenPayload(cfg, user.Email)

		if err != nil {
			log.Printf("Unable to create session: %e", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(data)

		if err != nil {
			log.Printf("Unable to write body: %e", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}
