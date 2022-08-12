package signups

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/adamstrickland/dapper-api/internal/security"
	"github.com/adamstrickland/dapper-api/internal/users"

	_ "github.com/mattn/go-sqlite3"
)

type requestPayload struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
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

		log.Printf("Received payload: '%+v'", qp)

		user := &users.User{
			Email:               qp.Email,
			UnencryptedPassword: qp.Password,
			FirstName:           qp.FirstName,
			LastName:            qp.LastName,
		}

		u, err := users.Create(cfg, user)

		if err != nil {
			log.Printf("Unable to create User: %e", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := security.NewTokenPayload(cfg, u.Email)

		if err != nil {
			log.Printf("Unable to create session: %e", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(data)

		w.Header().Set("Content-Type", "application/json")
	}
}
