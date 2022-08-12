package users

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/adamstrickland/dapper-api/internal/security"
)

type UserPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type usersPayload struct {
	Users []UserPayload `json:"users"`
}

func NewGetHandler(cfg *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data bytes.Buffer

		w.Header().Set("Content-Type", "application/json")

		users, err := All(cfg)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ups := make([]UserPayload, 0)

		for _, u := range *users {
			up := UserPayload{
				Email:     u.Email,
				FirstName: u.FirstName,
				LastName:  u.LastName,
			}
			ups = append(ups, up)
		}

		err = json.NewEncoder(&data).Encode(&usersPayload{
			Users: ups,
		})

		if err != nil {
			log.Printf("Unable to generate payload: %e", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(data.Bytes())

		if err != nil {
			log.Printf("Unable to write body: %e", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func NewPutHandler(cfg *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			data bytes.Buffer
			up   UserPayload
		)

		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&up)

		if err != nil {
			log.Printf("Unable to unmarshal payload: %e", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t := r.Header.Get(cfg.GetString("tokenHeader"))

		if t == "" {
			log.Printf("No token found")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		subj, err := security.TokenSubject(cfg, t)

		if err != nil {
			log.Printf("Subject could not be extracted from token")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		if *subj != up.Email {
			log.Printf("Token is not authorized to modify resource at '%s'", up.Email)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		u := &User{
			Email:     up.Email,
			FirstName: up.FirstName,
			LastName:  up.LastName,
		}

		uu, err := Update(cfg, u)

		if err != nil {
			log.Printf("Could not update user with email '%s': %e", u.Email, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		up = UserPayload{
			Email:     uu.Email,
			FirstName: uu.FirstName,
			LastName:  uu.LastName,
		}

		err = json.NewEncoder(&data).Encode(&up)

		if err != nil {
			log.Printf("Unable to generate payload: %e", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(data.Bytes())

		if err != nil {
			log.Printf("Unable to write body: %e", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
