package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"

	"golang.org/x/crypto/bcrypt"
)

func registerHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		// decode r.Body into req
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "bad json request", http.StatusBadRequest)
			return
		}

		// hash req.Password using bcrypt.GenerateFromPassword
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "failed to process registration", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(context.Background(), "INSERT INTO users (email, password_hash) VALUES ($1, $2)", req.Email, string(hash))
		if err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
