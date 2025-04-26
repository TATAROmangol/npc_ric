package main

import (
	"context"
	"forms/internal/config"
	"forms/internal/service"
	"forms/internal/storage"
	tablegrpc "forms/internal/transport/grpc/table"
	"forms/internal/transport/grpc/verify"
	httpserver "forms/internal/transport/http"
	"forms/internal/transport/http/handlers"
	"forms/internal/transport/http/middlewares"
	"forms/pkg/logger"
	"forms/pkg/migrator"
	"forms/pkg/postgres"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	migrationPath = "file://internal/migrations"
)

func main(){
	cfg := config.MustLoad()

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	db, err := postgres.NewDB(cfg.PG)
	if err != nil {
		l.ErrorContext(ctx, "failed to connect to postgres", err)
		os.Exit(1)
	}
	defer db.Close()
	l.InfoContext(ctx, "Connected to postgres", "path", cfg.PG.GetConnString())

	m, err := migrator.New(migrationPath, cfg.Migrator)
	if err != nil {
		l.ErrorContext(ctx, "failed in create migrator", err)
		os.Exit(1)
	}
	l.InfoContext(ctx, "migrator loaded")

	if err := m.Up(); err != nil {
		l.ErrorContext(ctx, "failed in up migrate", err)
		os.Exit(1)
	}
	l.InfoContext(ctx, "migrations complete")
	
	if err := m.Close(); err != nil {
		l.ErrorContext(ctx, "failed in close migrator", err)
		os.Exit(1)
	}
	l.InfoContext(ctx, "migrator closed")

	repo := storage.NewStorage(db)

	srv := service.NewServices(repo)


	ver := verify.NewClient(cfg.Verify)
	handler := handlers.NewHandlers(srv)
	midleware := middlewares.NewMiddlewares(ver)
	httpServer := httpserver.NewServer(ctx, cfg.HTTP, handler, midleware)

	grpcServer := tablegrpc.New(ctx, cfg.GRPC, srv)

	go func() {
		if err := grpcServer.Run(); err != nil {
			l.ErrorContext(ctx, "failed to run grpc server", err)
			os.Exit(1)
		}
	}()
	go func() {
		if err := httpServer.Run(); err != nil {
			l.ErrorContext(ctx, "failed to run http server", err)
			os.Exit(1)
		}
	}()
	
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	l.InfoContext(ctx, "Start graceful stop")

	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	l.InfoContext(ctx, "grpc server stopped")

	if err := httpServer.Shutdown(ctx); err != nil {
		l.ErrorContext(ctx, "failed to shutdown http server", err)
		os.Exit(1)
	}
	l.InfoContext(ctx, "http server stopped")

	if err := ver.Close(); err != nil {
		l.ErrorContext(ctx, "failed to close verify client", err)
		os.Exit(1)
	}
	l.InfoContext(ctx, "verify client closed")

	if err := db.Close(); err != nil {
		l.ErrorContext(ctx, "failed to close db", err)
		os.Exit(1)
	}
	l.InfoContext(ctx, "db closed")

	l.InfoContext(ctx, "all stopped")
}