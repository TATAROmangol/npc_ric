package kafka

type Config struct{
	Host string `env:"KAFKA_HOST"`
	Port string `env:"KAFKA_PORT"`
	LogTopic string `env:"KAFKA_LOG_TOPIC"`
}