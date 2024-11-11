package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/anoying-kid/go-apps/blogAPI/internal/handlers"
	"github.com/anoying-kid/go-apps/blogAPI/internal/middleware"
	"github.com/anoying-kid/go-apps/blogAPI/internal/repository"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	router *mux.Router
	db     *sql.DB
)

type TestUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Post struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	AuthorID int64  `json:"author_id"`
}

func TestMain(m *testing.M) {
	// Setup
	var err error
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=userdb_test port=5432 sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}

	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	postRepo := repository.NewPostRepository(db)
	postHandler := handlers.NewPostHandler(postRepo)

	router = mux.NewRouter()
	router.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/api/posts", middleware.AuthMiddleware(postHandler.Create)).Methods("POST")
	router.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(postHandler.Update)).Methods("PUT")
	router.HandleFunc("/api/posts/{id}", postHandler.Get).Methods("GET")
	router.HandleFunc("/api/posts", postHandler.List).Methods("GET")

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupDatabase()
	db.Close()

	os.Exit(code)
}

func cleanupDatabase() {
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM users")
}

func TestUserRegistrationAndLogin(t *testing.T) {
	cleanupDatabase()

	// Test user registration
	user := TestUser{
		Username: "testuser",
		Password: "testpass123",
	}

	// Register user
	t.Run("Register User", func(t *testing.T) {
		body, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	// Login user
	t.Run("Login User", func(t *testing.T) {
		body, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var loginResp LoginResponse
		err := json.NewDecoder(rr.Body).Decode(&loginResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, loginResp.Token)
	})
}

func TestPostOperations(t *testing.T) {
    cleanupDatabase()

    // First register and login to get token
    user := TestUser{
        Username: "postuser",
        Password: "postpass123",
    }

    // Register
    body, _ := json.Marshal(user)
    req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    // Login to get token
    req = httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    rr = httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    var loginResp LoginResponse
    json.NewDecoder(rr.Body).Decode(&loginResp)
    token := loginResp.Token

    // Test creating a post
    t.Run("Create Post", func(t *testing.T) {
        post := map[string]string{
            "title": "Test Post",
            "body":  "This is a test post body",
        }
        body, _ := json.Marshal(post)
        req := httptest.NewRequest("POST", "/api/posts", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+token)

        rr := httptest.NewRecorder()
        router.ServeHTTP(rr, req)

        assert.Equal(t, http.StatusCreated, rr.Code)

        var createdPost Post
        err := json.NewDecoder(rr.Body).Decode(&createdPost)
        assert.NoError(t, err)
        assert.Equal(t, post["title"], createdPost.Title)
        assert.Equal(t, post["body"], createdPost.Body)
    })

	t.Run("Update Post", func(t *testing.T) {
		// Create a post first to ensure we're the owner
		createPost := map[string]string{
			"title": "Initial Test Post",
			"body":  "This is the initial body",
		}
		createBody, _ := json.Marshal(createPost)
		createReq := httptest.NewRequest("POST", "/api/posts", bytes.NewBuffer(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createReq.Header.Set("Authorization", "Bearer "+token)
		
		createRr := httptest.NewRecorder()
		router.ServeHTTP(createRr, createReq)
		
		var createdPost Post
		json.NewDecoder(createRr.Body).Decode(&createdPost)
		postID := createdPost.ID
	
		// Now update the post we just created
		updatePost := map[string]string{
			"title": "Updated Test Post",
			"body":  "This is an updated test post body",
		}
		updateBody, _ := json.Marshal(updatePost)
	
		// Send update request
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/posts/%d", postID), bytes.NewBuffer(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	
		assert.Equal(t, http.StatusOK, rr.Code)
	
		var updatedPost Post
		err := json.NewDecoder(rr.Body).Decode(&updatedPost)
		assert.NoError(t, err)
		assert.Equal(t, postID, updatedPost.ID)
		assert.Equal(t, updatePost["title"], updatedPost.Title)
		assert.Equal(t, updatePost["body"], updatedPost.Body)
	})

    // Test listing posts
    t.Run("List Posts", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/posts", nil)
        rr := httptest.NewRecorder()
        router.ServeHTTP(rr, req)

        assert.Equal(t, http.StatusOK, rr.Code)

        var posts []Post
        err := json.NewDecoder(rr.Body).Decode(&posts)
        assert.NoError(t, err)
        assert.NotEmpty(t, posts)
    })

    // Test getting a specific post
    t.Run("Get Post", func(t *testing.T) {
        // First, list posts to get an ID
        req := httptest.NewRequest("GET", "/api/posts", nil)
        rr := httptest.NewRecorder()
        router.ServeHTTP(rr, req)

        var posts []Post
        json.NewDecoder(rr.Body).Decode(&posts)
        postID := posts[0].ID

        // Now get the specific post
        req = httptest.NewRequest("GET", fmt.Sprintf("/api/posts/%d", postID), nil)
        rr = httptest.NewRecorder()
        router.ServeHTTP(rr, req)

        assert.Equal(t, http.StatusOK, rr.Code)

        var post Post
        err := json.NewDecoder(rr.Body).Decode(&post)
        assert.NoError(t, err)
        assert.Equal(t, postID, post.ID)
    })
}

func TestInvalidOperations(t *testing.T) {
	cleanupDatabase()

	// Test creating post without authentication
	t.Run("Create Post Without Auth", func(t *testing.T) {
		post := map[string]string{
			"title": "Test Post",
			"body":  "This is a test post body",
		}
		body, _ := json.Marshal(post)
		req := httptest.NewRequest("POST", "/api/posts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	// Test invalid login
	t.Run("Invalid Login", func(t *testing.T) {
		user := TestUser{
			Username: "nonexistent",
			Password: "wrongpass",
		}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	// Test getting non-existent post
	t.Run("Get Non-existent Post", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/posts/99999", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
