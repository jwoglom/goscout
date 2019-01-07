package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server contains the port used to start up the server
// and a private mux.Router instance
type Server struct {
	Port   int
	router *mux.Router
}

// NewServer creates a Server
func NewServer(port int) *Server {
	return &Server{
		Port:   port,
		router: mux.NewRouter(),
	}
}

// Run adds the necessary routes and starts the server
func (s *Server) Run() {
	s.addRoutes()

	port := fmt.Sprintf(":%d", s.Port)
	fmt.Printf("Listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, s.router))
}
