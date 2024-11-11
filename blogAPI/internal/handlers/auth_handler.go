// handlers/auth_handler.go

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anoying-kid/go-apps/blogAPI/internal/middleware"
)

type AuthHandler struct {}

type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
    var req RefreshTokenRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate the refresh token
    userID, err := middleware.ValidateRefreshToken(req.RefreshToken)
    if err != nil {
        http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
        return
    }

    // Generate new token pair
    tokens, err := middleware.GenerateTokenPair(userID)
    if err != nil {
        http.Error(w, "Failed to generate new tokens", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tokens)
}