import logging
import sys
from kafka_logger import KafkaLoggingHandler


def setup_logger():
    logger = logging.getLogger()
    # Устанавливаем уровень логирования
    logger.setLevel(logging.DEBUG)

    # Формат вывода логов
    formatter = logging.Formatter(
        '[%(asctime)s] [%(levelname)s] [%(name)s] %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S'
    )

    # Обработчик логов в консоль
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setFormatter(formatter)

    # Обработчик логов в файл
    file_handler = logging.FileHandler("app.log", mode='a', encoding='utf-8')
    file_handler.setFormatter(formatter)

    # Обработчик логов в Kafka
    kafka_handler = KafkaLoggingHandler()
    kafka_handler.setFormatter(formatter)

    # Очищаем старые обработчики, если они уже есть
    if logger.hasHandlers():
        logger.handlers.clear()

    logger.addHandler(console_handler)
    logger.addHandler(file_handler)
    logger.addHandler(kafka_handler)

    return logger
