package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/anoying-kid/go-apps/blogAPI/internal/handlers"
	"github.com/anoying-kid/go-apps/blogAPI/internal/middleware"
	"github.com/anoying-kid/go-apps/blogAPI/internal/repository"
	"github.com/anoying-kid/go-apps/blogAPI/pkg/config"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// dsn := "host=localhost user=postgres password=mysecretpassword dbname=userdb port=5432 sslmode=disable"
	dbURL := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Fail to connect to the database: ", err)
	}
	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to database")
	defer db.Close()

	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	postRepo := repository.NewPostRepository(db)
	postHandler := handlers.NewPostHandler(postRepo)

	resetRepo := repository.NewPasswordResetRepository(db)
	resetHandler := handlers.NewPasswordResetHandler(userRepo, resetRepo, *cfg)

	// Setup router
	r := mux.NewRouter()

	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	// Protect routes with middleware
	r.HandleFunc("/api/posts", middleware.AuthMiddleware(postHandler.Create)).Methods("POST")
	r.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(postHandler.Update)).Methods("PUT")
	r.HandleFunc("/api/posts/{id}", postHandler.Get).Methods("GET")
	r.HandleFunc("/api/posts", postHandler.List).Methods("GET")

	r.HandleFunc("/api/password-reset", resetHandler.RequestReset).Methods("POST")
	r.HandleFunc("/api/password-reset/confirm", resetHandler.ConfirmReset).Methods("POST")

	// Start server
	log.Printf("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
