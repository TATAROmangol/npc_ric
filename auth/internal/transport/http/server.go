package httpserver

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type Service interface {
	Loginer
}

type Server struct {
	server *http.Server
}

func NewServer(ctx context.Context, cfg Config, srv Service) *Server {
	mux := mux.NewRouter()
	mux.Use(InitLoggerContextMiddleware(ctx))
	mux.Use(Operation())
	mux.Handle("/login", LoginHandler(srv)).Methods(http.MethodPost)

	return &Server{
		server: &http.Server{
			Addr:    cfg.Addr(),
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}