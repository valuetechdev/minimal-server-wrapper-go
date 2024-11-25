package vserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"slices"
)

type Middleware = func(http.Handler) http.Handler

type Server struct {
	// mux         *http.ServeMux
	server     *http.Server
	middleware []Middleware
	// options     *ServerOptions
}

type ServerOptions struct {
	// defaults to ":8080"
	Addr string

	// defaults to context.Background
	BaseCtx context.Context
}

// fill in the missing options
func (so *ServerOptions) addOptionDefaults() {
	if so.BaseCtx == nil {
		so.BaseCtx = context.Background()
	}

	if so.Addr == "" {
		so.Addr = ":8080"
	}
}

func New(options *ServerOptions) *Server {
	mux := http.NewServeMux()

	options.addOptionDefaults()

	server := &http.Server{
		Addr: options.Addr,
		BaseContext: func(l net.Listener) context.Context {
			return options.BaseCtx
		},
		Handler: mux,
	}

	return &Server{server, []Middleware{}}
}

func (s *Server) Serve() error {
	// apply middleware to server handler
	for _, middleware := range slices.Backward(s.middleware) {
		s.server.Handler = middleware(s.server.Handler)
	}

	slog.Info("Listening on " + s.server.Addr)

	return s.server.ListenAndServe()
}

// Allows for graceful shutdown of the server. See http package for more details
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) AddRoute(route string, handler http.Handler) {
	// assert type as we know what the handler type is
	s.server.Handler.(*http.ServeMux).Handle(route, handler)
}

func (s *Server) AddMiddleware(middleware Middleware) {
	s.server.Handler = middleware(s.server.Handler)
}
