package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/anoying-kid/go-apps/blogAPI/internal/handlers"
	"github.com/anoying-kid/go-apps/blogAPI/internal/repository"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main(){
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=userdb port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Fail to connect to the database: ",err)
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

	// Setup router
	r := mux.NewRouter()
	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/posts", postHandler.Create).Methods("POST")
    r.HandleFunc("/api/posts/{id}", postHandler.Get).Methods("GET")
    r.HandleFunc("/api/posts", postHandler.List).Methods("GET")

	// Start server
	log.Printf("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}