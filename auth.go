package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

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

func loginHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "failed to process login", http.StatusBadRequest)
			return
		}

		var id int
		var passwordHash string

		err = db.QueryRow(context.Background(), "SELECT id, password_hash FROM users WHERE email = $1", req.Email).Scan(&id, &passwordHash)
		if err != nil {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
		if err != nil {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		generatedJwt, err := generateJWT(id, os.Getenv("JWT_SECRET"))
		if err != nil {
			http.Error(w, "could not generate JWT", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": generatedJwt})

	}
}
