package httpserver

import (
	"context"
	"net/http"
)

type Server struct{
	ctx context.Context
	srv *http.Server
}

func NewServer(ctx context.Context,cfg Config) *Server {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}

	return &Server{
		ctx: ctx,
		srv: srv,
	}
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}