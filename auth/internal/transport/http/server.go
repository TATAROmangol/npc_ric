package httpserver

import (
	"auth/pkg/logger"
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type Service interface {
	Loginer
}

type Server struct {
	ctx context.Context
	server *http.Server
}

func NewServer(ctx context.Context, cfg Config, srv Service) *Server {
	mux := mux.NewRouter()
	mux.Use(InitLoggerContextMiddleware(ctx))
	mux.Use(Operation())
	mux.Handle("/login", LoginHandler(srv)).Methods(http.MethodPost)
	mux.Handle("/logout", LogoutHandler()).Methods(http.MethodPost)

	return &Server{
		ctx: ctx,
		server: &http.Server{
			Addr:    cfg.Addr(),
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	logger.GetFromCtx(s.ctx).InfoContext(s.ctx, "Run http", "path", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.GetFromCtx(s.ctx).ErrorContext(s.ctx, "failed to start server", err)
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}