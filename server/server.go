package server

import (
	"log"
	"net/http"
	"os"

	"github.com/Dacode45/addressbook/storage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	config ServerConfig
	router *mux.Router
}

func NewServer(u storage.UserStorage, config ServerConfig) *Server {
	s := Server{router: mux.NewRouter(), config: config}
	NewUserRouter(u, config, s.newSubrouter("/api/v1/users"))
	NewContactRouter(u, config, s.newSubrouter("/api/v1/contacts"))
	return &s
}

func (s *Server) Start() {
	log.Println("Listen on port 8080")
	if err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
		log.Fatal("http.ListenAndServer: ", err)
	}
}

func (s *Server) newSubrouter(path string) *mux.Router {
	return s.router.PathPrefix(path).Subrouter()
}
