package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	done chan struct{}
	cfg   Config
	reader *kafka.Reader
}

func New(cfg Config) (*Kafka, error) {
	k := &Kafka{
		done: make(chan struct{}),
		cfg:  cfg,
	}

	err := k.connectWithRetry(5, 10*time.Second)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (k *Kafka) connectWithRetry(maxAttempts int, interval time.Duration) error {
	var lastErr error
	
	for i := 0; i < maxAttempts; i++ {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{k.cfg.Addr()},
			Topic:   k.cfg.LogTopic,
		})

		conn, err := kafka.DialLeader(context.Background(), "tcp", k.cfg.Addr(), k.cfg.LogTopic, 0)
		if err == nil {
			conn.Close()
			k.reader = reader
			return nil
		}
		
		lastErr = err
		reader.Close()
		
		if i < maxAttempts-1 {
			log.Printf("Failed to connect to Kafka")
			time.Sleep(interval)
		}
	}
	
	return lastErr
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

