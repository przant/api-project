package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/przant/api-project/handlers"
	"github.com/przant/api-project/repositories"
)

type Server struct {
	router *mux.Router
	db     *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{
		router: mux.NewRouter(),
		db:     db,
	}
}

func (s *Server) setupRoutes() {
	// Create dependencies
	productRepo := repositories.NewProductRepository(s.db)
	productHandler := handlers.NewProductHandler(productRepo)

	// Health check endpoint
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}).Methods("GET")

	// Product routes
	s.router.HandleFunc("/products", productHandler.Create).Methods("POST")
	s.router.HandleFunc("/products", productHandler.List).Methods("GET")
	s.router.HandleFunc("/products/{id}", productHandler.Get).Methods("GET")
	s.router.HandleFunc("/products/{id}", productHandler.Update).Methods("PUT")
	s.router.HandleFunc("/products/{id}", productHandler.Delete).Methods("DELETE")
}

func (s *Server) Run(addr string) error {
	s.setupRoutes()
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/products?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	server := NewServer(db)
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
