package config

import (
	httpsend "forms/internal/transport/http/send"
	httpcreate "forms/internal/transport/http/create"
	"forms/pkg/migrator"
	"forms/pkg/postgres"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	PG postgres.Config
	Migrator migrator.Config
	HTTPCreate httpcreate.Config
	HTTPSend httpsend.Config
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