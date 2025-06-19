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
    institution_id: int  # ID организации


@router.post("/generate", response_class=StreamingResponse)
def generate(request: GenerateRequest):
    """
    Генерирует DOCX-документ по шаблону, сразу возвращает его пользователю.
    """
    db = SessionLocal()
    try:
        # Получаем шаблон по institution_id
        template = db.query(Template).filter_by(institution_id=request.institution_id).first()
        if not template:
            raise HTTPException(status_code=404, detail="Template not found")

        # Генерация документа из шаблона
        generated_bytes = generate_docx_from_template(template.content, request.institution_id)

        # Сохраняем сгенерированный документ 
        doc_id = str(uuid.uuid4())
        db_doc = GeneratedDocument(
            template_id=template.id,
            file_content=generated_bytes,
            filename=f"{doc_id}.docx"
        )
        db.add(db_doc)
        db.commit()

        # Возвращаем документ 
        return StreamingResponse(
            io.BytesIO(generated_bytes),
            media_type="application/vnd.openxmlformats-officedocument.wordprocessingml.document",
            headers={
                "Content-Disposition": f"attachment; filename={doc_id}.docx"
            }
        )
    finally:
        db.close()
