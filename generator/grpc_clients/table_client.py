import grpc
from generated import table_pb2, table_pb2_grpc
import os


def get_table_data(institution_id: int):
    host = os.getenv("FORMS_GRPC_HOST", "forms")
    port = os.getenv("FORMS_GRPC_PORT", "50050")
    with grpc.insecure_channel(f"{host}:{port}") as channel:
        stub = table_pb2_grpc.TableServiceStub(channel)
        response = stub.GetTable(
            table_pb2.GetTableRequest(institution_id=institution_id))
        return {
            "columns": list(response.columns),
            "rows": [[val for val in row.values] for row in response.rows]
        }
