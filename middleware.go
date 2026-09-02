package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "userID"

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if !strings.HasPrefix(bearer, "Bearer ") {
			http.Error(w, "bearer not valid", http.StatusUnauthorized)
			return
		}
		bearerToken := strings.TrimPrefix(bearer, "Bearer ")
		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !parsedToken.Valid {
			http.Error(w, "token is not valid", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["userID"].(float64)

		if !ok {
			http.Error(w, "claim missing or wrong type", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
