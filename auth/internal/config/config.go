package config

import (
	"auth/internal/admin"
	grpcserver "auth/internal/transport/grpc"
	httpserver "auth/internal/transport/http"
	"auth/pkg/jwt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP httpserver.Config
	GRPC grpcserver.Config
	Admin admin.Config
	JWT jwt.Config
}

func MustLoad() *Config {
	var cfg Config
	
	if err := godotenv.Load(); err != nil {
		log.Fatalf("not found env: %v", err)
	}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parse .env file: %v", err)
	}
	return &cfg
}