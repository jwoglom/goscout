package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const DefaultPort = 3000

type Server struct {
	Port   int
	router *mux.Router
}

func NewServer() *Server {
	return &Server{
		Port:   DefaultPort,
		router: mux.NewRouter(),
	}
}

func (s *Server) Run() {
	s.addRoutes()

	port := fmt.Sprintf(":%d", s.Port)
	log.Fatal(http.ListenAndServe(port, s.router))
}
