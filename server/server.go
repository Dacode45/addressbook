package server

import (
	"log"
	"net/http"
	"os"

	"github.com/Dacode45/addressbook/storage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Server sets up the api, and serves over http
type Server struct {
	config ServerConfig
	router *mux.Router
}

// NewServer creates a new Server given a storage backend and configuration
func NewServer(u storage.UserStorage, config ServerConfig) *Server {
	s := Server{router: mux.NewRouter(), config: config}
	NewUserRouter(u, config, s.newSubrouter("/api/v1/users"))
	NewContactRouter(u, config, s.newSubrouter("/api/v1/contacts"))
	return &s
}

// Start serves the server on port 8080
func (s *Server) Start() error {
	log.Println("Listen on port 8080")
	return http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, s.router))
}

func (s *Server) newSubrouter(path string) *mux.Router {
	return s.router.PathPrefix(path).Subrouter()
}
