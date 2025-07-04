package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"moveshare/internal/auth"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
)

// AuthMiddleware authenticates requests using JWT
func AuthMiddleware(jwtAuth auth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			log.Printf("Received Authorization header: %q", authHeader)
			if authHeader == "" {
				log.Printf("No Authorization header found")
				http.Error(w, `{"error": "authorization header required"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			log.Printf("Extracted token: %q", tokenStr)
			if tokenStr == authHeader {
				log.Printf("Invalid Authorization header format: %q", authHeader)
				http.Error(w, `{"error": "invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			userID, err := jwtAuth.ValidateToken(tokenStr)
			if err != nil {
				log.Printf("Token validation error for token %q: %v", tokenStr, err)
				http.Error(w, `{"error": "invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			log.Printf("Token validated successfully for user ID: %d", userID)
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			log.Printf("Context updated with user ID: %d", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
