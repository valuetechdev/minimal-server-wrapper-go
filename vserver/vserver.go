package vserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"
)

type Middleware = func(http.Handler) http.Handler

type Server struct {
	port        string
	mux         *http.ServeMux
	middlewares []Middleware // The order of middlewares are important
}

func New(port string) *Server {
	mux := http.NewServeMux()

	return &Server{
		port: port,
		mux:  mux,
	}
}

func (s *Server) Serve() error {
	serveMux := http.Handler(s.mux)

	for _, middleware := range slices.Backward(s.middlewares) {
		serveMux = middleware(serveMux)
	}

	addr := fmt.Sprintf(":%s", s.port)
	slog.Info("Listening on " + addr)

	return http.ListenAndServe(addr, serveMux)
}

func (s *Server) AddRoute(route string, handler http.Handler) {
	s.mux.Handle(route, handler)
}

func (s *Server) AddMiddleware(middleware Middleware) {
	s.middlewares = append(s.middlewares, middleware)
}
