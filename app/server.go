package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jwoglom/goscout/app/db"
	"github.com/jwoglom/goscout/app/endpointsv1"

	"github.com/gorilla/mux"
)

// Server contains the port used to start up the server
// and a private mux.Router instance
type Server struct {
	Port   int
	Db     *db.Db
	router *mux.Router

	Package *Package
}

// Package provides accessors to different subpackages
type Package struct {
	EndpointsV1 *endpointsv1.EndpointsV1
}

// NewServer creates a Server
func NewServer(port int) *Server {
	server := &Server{
		Port:   port,
		Db:     db.NewDb(),
		router: mux.NewRouter(),
	}

	server.genPackage()
	return server
}

func (s *Server) genPackage() {
	s.Package = &Package{
		EndpointsV1: &endpointsv1.EndpointsV1{
			Db: s.Db,
		},
	}
}

// Run adds the necessary routes and starts the server
func (s *Server) Run() {
	s.addRoutes()

	port := fmt.Sprintf(":%d", s.Port)
	fmt.Printf("Listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, s.router))
}
