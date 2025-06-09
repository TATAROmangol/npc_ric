package config

import (
	tablegrpc "forms/internal/transport/grpc/table"
	"forms/internal/transport/grpc/verify"
	httpserver "forms/internal/transport/http"
	"forms/pkg/kafka"
	"forms/pkg/postgres"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	PG postgres.Config
	HTTP httpserver.Config
	Verify verify.Config
	GRPC tablegrpc.Config
	Kafka kafka.Config
}

func MustLoad() *Config {
	var cfg Config
	
	godotenv.Load()

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parse .env file: %v", err)
	}

	var test Config
	if cfg == test{
		log.Fatalf("Error load cfg file")
	}

	return &cfg
}