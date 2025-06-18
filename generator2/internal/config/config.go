package config

import (
	"generator/internal/transport/grpc/table"
	"generator/internal/transport/grpc/verify"
	"generator/internal/transport/http"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct{
	GRPCTable table.Config
	GRPCAuth verify.Config
	HTTP http.Config
}

func MustLoad() *Config {
	var cfg Config 

	godotenv.Load()

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	if cfg == (Config{}) {
		log.Fatalf("failed to load config: %s", "empty")
	}

	return &cfg
}
