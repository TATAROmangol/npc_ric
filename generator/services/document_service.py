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
        template = db.query(Template).filter_by(name=template_name).first()
        if not template:
            raise HTTPException(status_code=404, detail="Template not found")

        generated_bytes = generate_docx_from_template(template.content,
                                                      institution_id)

        doc_id = str(uuid.uuid4())
        filename = f"{doc_id}.docx"

        new_doc = GeneratedDocument(
            template_id=template.id,
            file_content=generated_bytes,
            filename=filename
        )
        db.add(new_doc)
        db.commit()
        db.refresh(new_doc)

        return {"id": new_doc.id, "filename": new_doc.filename}
    except Exception:
        logger.exception("Error while generating document")
        db.rollback()
        raise
    finally:
        db.close()


def get_document_by_id(doc_id: int) -> Optional[dict]:
    db = SessionLocal()
    try:
        doc = db.query(GeneratedDocument).filter_by(id=doc_id).first()
        if not doc:
            return None
        return {"filename": doc.filename, "content": doc.file_content}
    finally:
        db.close()
