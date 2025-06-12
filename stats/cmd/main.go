package main

import (
	"log"
	"stats/config"
	"stats/kafka"
)

func main() {
	cfg := config.MustLoad()

	kafkaReader, err := kafka.New(cfg.Kafka)
	if err != nil {
		log.Fatal("Failed to create Kafka reader", "error", err)
	}
	defer kafkaReader.Close()

	ch := make(chan string)
	kafkaReader.StartReading(ch)

	for msg := range ch {
		log.Println(msg)
	}
}
