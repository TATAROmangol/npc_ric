package config

import (
	"log"
	"stats/kafka"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
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