from fastapi import APIRouter, HTTPException
from fastapi.responses import StreamingResponse
from pydantic import BaseModel
import io
import uuid
from db.database import SessionLocal
from db.models import GeneratedDocument, Template
from services.docx_service import generate_docx_from_template
import logging

logger = logging.getLogger(__name__)

router = APIRouter()


class GenerateRequest(BaseModel):
    institution_id: int # Входная модель: ID организации


@router.post("/generate")
def generate(request: GenerateRequest):
    """
    Генерирует DOCX-документ по шаблону, привязанному к переданному institution_id.
    Сохраняет результат в базе и возвращает ссылку для скачивания.
    """
    db = SessionLocal()
    try:
        # Получаем шаблон по institution_id
        template = db.query(Template).filter_by(
                    institution_id=request.institution_id).first()
        if not template:
            raise HTTPException(status_code=404, detail="Template not found")

        # Генерация документа из шаблона
        generated_bytes = generate_docx_from_template(
                                template.content,
                                request.institution_id)

        # Сохраняем сгенерированный документ в базе
        doc_id = str(uuid.uuid4())
        db_doc = GeneratedDocument(
            template_id=template.id,
            file_content=generated_bytes,
            filename=f"{doc_id}.docx"
        )
        db.add(db_doc)
        db.commit()

        # Возвращаем URL для скачивания
        return {"download_url": f"/documents/download/{db_doc.id}"}
    finally:
        db.close()


@router.get("/download/{doc_id}")
def download(doc_id: int):
    """
    Отдаёт ранее сгенерированный документ по его ID в виде потока (StreamingResponse).
    """
    db = SessionLocal()
    try:
        # Ищем документ по ID
        doc = db.query(GeneratedDocument).filter_by(id=doc_id).first()
        if not doc:
            raise HTTPException(status_code=404,
                                detail="Generated file not found")

        # Возвращаем документ в формате DOCX
        return StreamingResponse(
            io.BytesIO(doc.file_content),
            media_type=(
                "application/vnd.openxmlformats-"
                "officedocument.wordprocessingml.document"),
            headers={
                "Content-Disposition": "attachment; filename=generated.docx"
            }
        )
    finally:
        db.close()
