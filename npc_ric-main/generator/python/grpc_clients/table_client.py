import grpc
from generated import table_pb2, table_pb2_grpc


def get_table_data(institution_id: int):
    with grpc.insecure_channel("localhost:50052") as channel:
        stub = table_pb2_grpc.TableServiceStub(channel)
        response = stub.GetTable(
            table_pb2.GetTableRequest(institution_id=institution_id))
        return {
            "columns": list(response.columns),
            "rows": [[val for val in row.values] for row in response.rows]
        }
