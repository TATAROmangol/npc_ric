package config

import (
	"generator/internal/docx"
	"generator/internal/transport/grpc/table"
	"generator/internal/transport/grpc/verify"
	"generator/internal/transport/httpvo1"
	"generator/pkg/kafka"
	"generator/pkg/mongodb"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct{
	GRPCTable table.Config
	GRPCAuth verify.Config
	HTTP httpvo1.Config
	Kafka kafka.Config
	Mongo mongodb.Config
	Docx docx.Config
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
