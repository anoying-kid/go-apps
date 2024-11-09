package middleware

import (
	"context"
	"net/http"
	"strings"
)

// Custom type for context keys
type contextKey string

const UserIDKey contextKey = "user_id"

type Claims struct {
	UserID int64 `json:"user_id"`
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get the Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        // Check if the header starts with "Bearer "
        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
            http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
            return
        }

        // Validate the JWT token
        claims, err := ValidateToken(bearerToken[1])
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add the user ID to the request context
        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}
