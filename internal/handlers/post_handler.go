package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anoying-kid/go-apps/blogAPI/internal/middleware"
	"github.com/anoying-kid/go-apps/blogAPI/internal/models"
	"github.com/anoying-kid/go-apps/blogAPI/internal/repository"
	"github.com/gorilla/mux"
)

type PostHandler struct {
	postRepo *repository.PostRepository
}

type UpdatePostRequest struct {
    Title string `json:"title"`
    Body  string `json:"body"`
}

func NewPostHandler(postRepo *repository.PostRepository) *PostHandler {
	return &PostHandler{postRepo: postRepo}
}

type CreatePostRequest struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Get user ID from context
    userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// authorID := int64(1) // Replace with actual author ID
	
	post := &models.Post{
		Title:     req.Title,
		Body:      req.Body,
		AuthorID:  userID,
	}

	if err := h.postRepo.Create(post); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)

}

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    post, err := h.postRepo.GetByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if post == nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    postID, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    existingPost, err := h.postRepo.GetByID(postID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if existingPost == nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    // Check if the user is the author
    if existingPost.AuthorID != userID {
        http.Error(w, "Unauthorized: you are not the author of this post", http.StatusForbidden)
        return
    }

    // Decode the update request
    var req UpdatePostRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Update the post
    existingPost.Title = req.Title
    existingPost.Body = req.Body

    if err := h.postRepo.Update(existingPost); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(existingPost)
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
    limit := 10
    offset := 0

    if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
        if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
            limit = parsedLimit
        }
    }
    if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
        if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
            offset = parsedOffset
        }
    }

    posts, err := h.postRepo.List(limit, offset)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}