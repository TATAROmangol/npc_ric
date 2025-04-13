package httpcreate

import (
	"context"
	"net/http"
)

type CreateServer struct{
	ctx context.Context
	srv *http.Server
}

func NewServer(ctx context.Context,cfg Config) *CreateServer {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}

	return &CreateServer{
		ctx: ctx,
		srv: srv,
	}
}

func (s *CreateServer) Run() error {
	return s.srv.ListenAndServe()
}

func (s *CreateServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}