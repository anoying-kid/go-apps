package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/anoying-kid/go-apps/blogAPI/internal/models"
	"github.com/anoying-kid/go-apps/blogAPI/internal/repository"
	"github.com/anoying-kid/go-apps/blogAPI/pkg/config"
	"github.com/anoying-kid/go-apps/blogAPI/pkg/utils"
)

type PasswordResetHandler struct {
	userRepo *repository.UserRepository
	resetRepo *repository.PasswordResetRepository
	config config.Config
}

func NewPasswordResetHandler(
	userRepo *repository.UserRepository,
	resetRepo *repository.PasswordResetRepository,
	config config.Config) *PasswordResetHandler {

	return &PasswordResetHandler{
		userRepo: userRepo,
		resetRepo: resetRepo,
		config: config}
}

type PasswordResetRequest struct {
    Email string `json:"email"`
}

func (h *PasswordResetHandler) RequestReset(w http.ResponseWriter, r *http.Request) {
	var req PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if user == nil {
		// Don't reveal email that doesn't exist
		w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{
            "message": "If your email exists in our system, you will receive reset instructions",
        })
        return
		}

	// Generate secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Create reset token record
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: time.Now().Add(time.Hour), // Token expires in 1 hour
	}
	if err := h.resetRepo.Create(resetToken); err != nil {
		http.Error(w, "Error creating reset token", http.StatusInternalServerError)
		log.Fatal("Error creating reset token: ", err)
		return
	}

	// Send reset email
	if err := utils.SendPasswordResetEmail(user.Email, token, h.config); err != nil {
		http.Error(w, "Error sending email", http.StatusInternalServerError)
		log.Fatal("Error sending email: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If your email exists in our system, you will receive reset instructions",
	})
}

type PasswordResetConfirmRequest struct {
    Token    string `json:"token"`
    Password string `json:"password"`
}

func (h *PasswordResetHandler) ConfirmReset(w http.ResponseWriter, r *http.Request) {
    var req PasswordResetConfirmRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Validate token
    resetToken, err := h.resetRepo.GetByToken(req.Token)
    if err != nil {
        http.Error(w, "Error validating token", http.StatusInternalServerError)
        return
    }
    if resetToken == nil || resetToken.Used || time.Now().After(resetToken.ExpiredAt) {
        http.Error(w, "Invalid or expired token", http.StatusBadRequest)
        return
    }

    // Hash new password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        http.Error(w, "Error processing password", http.StatusInternalServerError)
        return
    }

    // Update user's password
    if err := h.userRepo.UpdatePassword(resetToken.UserID, hashedPassword); err != nil {
        http.Error(w, "Error updating password", http.StatusInternalServerError)
        return
    }

    // Mark token as used
    if err := h.resetRepo.MarkAsUsed(resetToken.ID); err != nil {
        http.Error(w, "Error updating token status", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Password has been successfully reset",
    })
}