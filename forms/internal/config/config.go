package config

import (
	tablegrpc "forms/internal/transport/grpc/table"
	"forms/internal/transport/grpc/verify"
	httpserver "forms/internal/transport/http"
	"forms/pkg/migrator"
	"forms/pkg/postgres"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	PG postgres.Config
	Migrator migrator.Config
	HTTP httpserver.Config
	Verify verify.Config
	GRPC tablegrpc.Config
}

func MustLoad() *Config {
	var cfg Config
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parse .env file: %v", err)
	}

	return &cfg
}