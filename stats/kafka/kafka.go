package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	done chan struct{}
	reader *kafka.Reader
}

func New(cfg Config) *Kafka {
	return &Kafka{
		done: make(chan struct{}),
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{cfg.Addr()},
			Topic:   cfg.LogTopic,
		}),
	}
}

func (k *Kafka) StartReading(ch chan string) {
	go func(){
		for {
			select {
			case k.done <- struct{}{}:
				return
			default:
				msg, err := k.reader.ReadMessage(context.Background())
				if err != nil {
					log.Fatal("Failed to read message from Kafka:", err)
					return
				}
				ch <- string(msg.Value)
			}
		}
	}()
}

func (k *Kafka) Close() error {
	close(k.done)
	return k.reader.Close()
}

