package server

import (
	"fmt"
	"net/http"

	"github.com/lucasmenendez/psi/internal/session"
)

type Server struct {
	uri      string
	port     int
	sessions *session.Manager
}

func Init(port int) *Server {
	return &Server{
		uri:      fmt.Sprintf("0.0.0.0:%d", port),
		port:     port,
		sessions: session.NewManager(),
	}
}

func (server *Server) Start() {
	http.HandleFunc("/newIntersection", handleIntersectionRequest)

	http.ListenAndServe(server.uri, nil)
}
