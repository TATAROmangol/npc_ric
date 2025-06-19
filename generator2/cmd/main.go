package main

import (
	"context"
	"generator/internal/config"
	"generator/internal/docx"
	"generator/internal/service"
	"generator/internal/storage"
	"generator/internal/transport/grpc/table"
	"generator/internal/transport/grpc/verify"
	"generator/internal/transport/httpvo1"
	"generator/internal/transport/httpvo1/handlers"
	"generator/internal/transport/httpvo1/midlewares"
	"generator/pkg/kafka"
	"generator/pkg/logger"
	"generator/pkg/mongodb"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	kafka := kafka.New(cfg.Kafka)
	w := io.MultiWriter(os.Stdout, kafka)
	l := logger.New(w)
	ctx := context.Background()
	logger.InitFromCtx(ctx, l)
	logger.AppendCtx(ctx, "service", "generator")
	l.InfoContext(ctx, "Config loaded", cfg)

	db, err := mongodb.NewMongoDB(ctx, cfg.Mongo)
	if err != nil {
		l.ErrorContext(ctx, "failed to connect to MongoDB", err)
		os.Exit(1)
	}
	defer db.Close()
	l.InfoContext(ctx, "Connected to MongoDB", "path", cfg.Mongo.Addr())

	repo := storage.NewStorage(db.DB)

	generator, err := docx.NewGenerator(cfg.Docx)
	if err != nil {
		l.ErrorContext(ctx, "failed to create docx generator", err)
		os.Exit(1)
	}

	tabler := table.New(cfg.GRPCTable)
	l.InfoContext(ctx, "Connected to gRPC table service", "addr", cfg.GRPCTable.Addr())
	defer tabler.Close()

	auth := verify.New(cfg.GRPCAuth)
	l.InfoContext(ctx, "Connected to gRPC auth service", "addr", cfg.GRPCAuth.Addr())
	defer auth.Close()

	m := midlewares.NewMidelware(ctx, cfg.HTTP.AuthCookieName, auth)

	service := service.New(repo, generator, tabler)
	h := handlers.New(service)

	httpServer := httpvo1.New(cfg.HTTP, m, h)
	go func(){
		if err := httpServer.Run(); err != nil {
			l.ErrorContext(ctx, "failed to run HTTP server", err)
			os.Exit(1)
		}
	}()
	
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	l.InfoContext(ctx, "Start graceful stop")

	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()

	if err := httpServer.Stop(); err != nil {
		l.ErrorContext(ctx, "failed to shutdown HTTP server", err)
		os.Exit(1)
	}

	if err := tabler.Close(); err != nil {
		l.ErrorContext(ctx, "failed to close gRPC table client", err)
		os.Exit(1)
	}

	if err := auth.Close(); err != nil {
		l.ErrorContext(ctx, "failed to close gRPC auth client", err)
		os.Exit(1)
	}

	if err := db.Close(); err != nil {
		l.ErrorContext(ctx, "failed to close MongoDB connection", err)
		os.Exit(1)
	}

	if err := kafka.Close(); err != nil {
		l.ErrorContext(ctx, "failed to close Kafka writer", err)
		os.Exit(1)
	}

	l.InfoContext(ctx, "Graceful stop complete")
}