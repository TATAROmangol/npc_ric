package httpsend

import (
	"context"
	"net/http"
)

type SendServer struct{
	ctx context.Context
	srv *http.Server
}

func NewServer(ctx context.Context,cfg Config) *SendServer {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}

	return &SendServer{
		ctx: ctx,
		srv: srv,
	}
}

func (s *SendServer) Run() error {
	return s.srv.ListenAndServe()
}

func (s *SendServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}