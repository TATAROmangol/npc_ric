import logging
from confluent_kafka import Producer
import json
import os

KAFKA_BROKER = os.getenv("KAFKA_BROKER", "kafka:9092")
KAFKA_TOPIC = os.getenv("KAFKA_LOG_TOPIC", "log")

producer = Producer({"bootstrap.servers": KAFKA_BROKER})


class KafkaLoggingHandler(logging.Handler):
    def emit(self, record):
        try:
            log_entry = self.format(record)
            log_data = {
                "message": log_entry,
                "level": record.levelname,
                "logger": record.name,
            }
            producer.produce(KAFKA_TOPIC, json.dumps(log_data).encode('utf-8'))
            producer.flush()
            producer.poll(0)
        except Exception as e:
            print(f"Kafka logging error: {e}")
