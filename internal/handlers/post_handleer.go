package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anoying-kid/go-apps/blogAPI/internal/models"
	"github.com/anoying-kid/go-apps/blogAPI/internal/repository"
	"github.com/gorilla/mux"
)

type PostHandler struct {
	postRepo *repository.PostRepository
}

func NewPostHandler(postRepo *repository.PostRepository) *PostHandler {
	return &PostHandler{postRepo: postRepo}
}

type CreatePostRequest struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authorID := int64(1)
	
	post := &models.Post{
		Title:     req.Title,
		Body:      req.Body,
		AuthorID:  authorID,
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