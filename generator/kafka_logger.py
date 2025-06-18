import logging
from confluent_kafka import Producer
import json
import os

# Получаем настройки Kafka из переменных окружения
KAFKA_BROKER = os.getenv("KAFKA_BROKER", "kafka:9092")
KAFKA_TOPIC = os.getenv("KAFKA_LOG_TOPIC", "log")

# Инициализируем Kafka-производителя
producer = Producer({"bootstrap.servers": KAFKA_BROKER})

# Кастомный логгер, отправляющий логи в Kafka
class KafkaLoggingHandler(logging.Handler):
    def emit(self, record):
        try:
            # Форматируем лог-сообщение и сериализуем его в JSON
            log_entry = self.format(record)
            log_data = {
                "message": log_entry,
                "level": record.levelname,
                "logger": record.name,
            }
             # Отправляем сообщение в Kafka
            producer.produce(KAFKA_TOPIC, json.dumps(log_data).encode('utf-8'))
            producer.flush()
            producer.poll(0)
        except Exception as e:
            print(f"Kafka logging error: {e}")
