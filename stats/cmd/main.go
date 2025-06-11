package main

import (
	"log/slog"
	"stats/config"
	"stats/kafka"
)

func main() {
	cfg := config.MustLoad()

	kafkaReader := kafka.New(cfg.Kafka)
	defer kafkaReader.Close()

	ch := make(chan string)
	kafkaReader.StartReading(ch)

	for msg := range ch {
		slog.Info(msg)
	}
}
