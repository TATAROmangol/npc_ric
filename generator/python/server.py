from concurrent import futures
import grpc
from python.db.database import SessionLocal
from python.db.models import Template
from python.generated import docx_generator_pb2, docx_generator_pb2_grpc
from python.services.docx_service import DocxGeneratorService


class DocxGeneratorServicer(docx_generator_pb2_grpc.DocxGeneratorServicer):
    def UploadTemplate(self, request, context):
        db = SessionLocal()
        try:
            print(f"'{request.name}' {len(request.docx_content)} bytes")

            if not request.docx_content:
                context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                return docx_generator_pb2.Response(status="Empty file content")

            if db.query(Template).filter_by(name=request.name).first():
                context.set_code(grpc.StatusCode.ALREADY_EXISTS)
                return docx_generator_pb2.Response(status="Template exists")

            db.add(Template(
                name=request.name,
                content=request.docx_content
            ))
            db.commit()
            return docx_generator_pb2.Response(status="OK")
        except Exception as e:
            db.rollback()
            print(f"Upload error: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            return docx_generator_pb2.Response(status=f"Error: {str(e)}")
        finally:
            db.close()

    def GenerateDocx(self, request, context):
        db = SessionLocal()
        try:
            template = db.query(Template).filter_by(
                name=request.template_name).first()
            if not template:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return docx_generator_pb2.GenerateResponse()

            data = {
                    "rows": [
                        {
                            col: row.values[i]
                            for i, col in enumerate(request.table.columns)
                        }
                        for row in request.table.rows
                            ]
                    }

            content = DocxGeneratorService.generate_docx(
                template_name=request.template_name,
                data=data
            )

            if not content:
                context.set_code(grpc.StatusCode.INTERNAL)
                return docx_generator_pb2.GenerateResponse()

            return docx_generator_pb2.GenerateResponse(result_docx=content)
        except Exception as e:
            print(f"Error in GenerateDocx: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return docx_generator_pb2.GenerateResponse()
        finally:
            db.close()


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    docx_generator_pb2_grpc.add_DocxGeneratorServicer_to_server(
        DocxGeneratorServicer(),
        server
    )
    server.add_insecure_port("[::]:50051")
    server.start()
    print("Server running on port 50051")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
