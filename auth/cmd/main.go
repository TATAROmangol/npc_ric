package main

import (
	"auth/internal/admin"
	"auth/internal/config"
	"auth/internal/service"
	grpcserver "auth/internal/transport/grpc"
	httpserver "auth/internal/transport/http"
	"auth/pkg/jwt"
	"auth/pkg/kafka"
	"auth/pkg/logger"
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	kafkaWriter := kafka.New(cfg.Kafka)
	defer kafkaWriter.Close()
	writer := io.MultiWriter(os.Stdout, kafkaWriter)

	ctx := context.Background()
	l := logger.New(writer)
	ctx = logger.InitFromCtx(ctx, l)

	ctx = logger.AppendCtx(ctx, "service", "auth")

	jwt := jwt.New(cfg.JWT)
	admin := admin.New(cfg.Admin)
	service := service.New(admin, jwt)

	grpc := grpcserver.NewServer(ctx, cfg.GRPC, service)
	http := httpserver.NewServer(ctx, cfg.HTTP, service)

	go func() {
		if err := grpc.Run(); err != nil {
			l.ErrorContext(ctx, "failed to run grpc server", err)
		}
	}()

	go func() {
		if err := http.Run(); err != nil {
			l.ErrorContext(ctx, "failed to run http server", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch 

	l.InfoContext(ctx, "Shutdown signal received")
	grpc.GracefulStop()
	if err := http.Shutdown(ctx); err != nil {
		l.ErrorContext(ctx, "failed to shutdown http server", err)
	}
	l.InfoContext(ctx, "Servers stopped")
}