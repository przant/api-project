package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	return &Server{
		router: mux.NewRouter(),
	}
}

func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}).Methods("GET")
}

func (s *Server) Run(addr string) error {
	s.setupRoutes()
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}

func main() {
	server := NewServer()
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
