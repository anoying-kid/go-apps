package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/anoying-kid/go-apps/blogAPI/internal/middleware"
	"github.com/anoying-kid/go-apps/blogAPI/internal/models"
	"github.com/anoying-kid/go-apps/blogAPI/internal/repository"

	"github.com/anoying-kid/go-apps/blogAPI/pkg/utils"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var registerReq RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    hashedPassword, err := utils.HashPassword(registerReq.Password)
    if err != nil {
        http.Error(w, "Error processing password", http.StatusInternalServerError)
        return
    }

    user := &models.User{
        Username: registerReq.Username,
        Email:    registerReq.Email,
        Password: hashedPassword,
		CreatedAt: time.Now().Format(time.RFC3339),
    	UpdatedAt: time.Now().Format(time.RFC3339),
    }

    if err := h.userRepo.Create(user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "id":      user.ID,
        "username": user.Username,
        "email":   user.Email,
        "message": "User created successfully",
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
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
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}