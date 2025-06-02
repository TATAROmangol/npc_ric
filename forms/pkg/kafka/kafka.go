package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Kafka struct{
	writer *kafka.Writer
}

func New(cfg Config) *Kafka {
	return &Kafka{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{cfg.Addr()},
			Topic:  cfg.LogTopic,
		}),
	}
}

func (k *Kafka) Write(message []byte) (int, error) {
	if err := k.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: message,
		},
	); err != nil {
		return 0, err
	}

	return len(message), nil
}

func (k *Kafka) Close() error {
	return k.writer.Close()
}