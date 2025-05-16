package grpcserver

import (
	"auth/pkg/logger"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct{
	ctx context.Context
	cfg Config
	server *grpc.Server
}

func NewServer(ctx context.Context, cfg Config, ver Verificate) *Server{
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			InitLogger(ctx),
			Operation(),
		),
	)

	Register(server, ver)
	return &Server{
		ctx: ctx,
		cfg: cfg,
		server: server,
	}
}

func (s *Server) Run() error {
	logger.GetFromCtx(s.ctx).InfoContext(s.ctx, "Run grpc", "path", s.cfg.Addr())

	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", s.cfg.Host, s.cfg.Port))
	if err != nil {
		logger.GetFromCtx(s.ctx).ErrorContext(s.ctx, "failed create listener", err)
		return err
	}

	err = s.server.Serve(lis)
	if err != nil {
		logger.GetFromCtx(s.ctx).ErrorContext(s.ctx, "failed create server", err)
		return err
	}

	return nil
}

func (s *Server) GracefulStop() {
	s.server.GracefulStop()
}