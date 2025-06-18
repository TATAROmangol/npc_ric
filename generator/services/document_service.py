import uuid
from db.database import SessionLocal
from db.models import Template, GeneratedDocument
from services.docx_service import generate_docx_from_template
from fastapi import HTTPException
from typing import Optional
import logging

logger = logging.getLogger(__name__)


def generate_document(template_name: str, institution_id: int) -> dict:
    db = SessionLocal()
    try:
        # Ищем шаблон по имени
        template = db.query(Template).filter_by(name=template_name).first()
        if not template:
            raise HTTPException(status_code=404, detail="Template not found")

        # Генерируем документ в байтовом виде на основе шаблона и institution_id
        generated_bytes = generate_docx_from_template(template.content,
                                                      institution_id)

        # Генерируем уникальное имя файла
        doc_id = str(uuid.uuid4())
        filename = f"{doc_id}.docx"

        # Сохраняем сгенерированный документ в базу данны
        new_doc = GeneratedDocument(
            template_id=template.id,
            file_content=generated_bytes,
            filename=filename
        )
        db.add(new_doc)
        db.commit()
        db.refresh(new_doc)

        # Возвращаем информацию о документе
        return {"id": new_doc.id, "filename": new_doc.filename}
    except Exception:
        # Логируем ошибку и откатываем транзакцию
        logger.exception("Error while generating document")
        db.rollback()
        raise
    finally:
        db.close()


def get_document_by_id(doc_id: int) -> Optional[dict]:
    db = SessionLocal()
    try:
        # Ищем документ по ID
        doc = db.query(GeneratedDocument).filter_by(id=doc_id).first()
        if not doc:
            return None
        # Возвращаем содержимое и имя файла
        return {"filename": doc.filename, "content": doc.file_content}
    finally:
        db.close()
