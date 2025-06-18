import grpc
from generated import table_pb2, table_pb2_grpc
import os
import logging

logger = logging.getLogger(__name__)


def get_table_data(institution_id: int):
    """
    Делает gRPC-запрос к внешнему сервису, чтобы получить табличные данные
    (столбцы и строки) по ID организации.
    """
    logger.info(f"Requesting table data for institution {institution_id}")

    # Получаем адрес и порт gRPC-сервиса из переменных окружени
    host = os.getenv("FORMS_GRPC_HOST", "forms")
    port = os.getenv("FORMS_GRPC_PORT", "50050")

    # Устанавливаем соединение и делаем запрос
    with grpc.insecure_channel(f"{host}:{port}") as channel:
        stub = table_pb2_grpc.TableServiceStub(channel)
        response = stub.GetTable(
            table_pb2.GetTableRequest(institution_id=institution_id))
        
        # Преобразуем ответ в словарь с колонками и строками
        return {
            "columns": list(response.columns),
            "rows": [[val for val in row.values] for row in response.rows]
        }
