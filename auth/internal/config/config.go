package config

import (
	"auth/internal/admin"
	grpcserver "auth/internal/transport/grpc"
	httpserver "auth/internal/transport/http"
	"auth/pkg/jwt"
	"auth/pkg/kafka"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP httpserver.Config
	GRPC grpcserver.Config
	Admin admin.Config
	JWT jwt.Config
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