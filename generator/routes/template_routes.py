from fastapi import APIRouter, UploadFile, File, Form, HTTPException
from services import template_service
from db.database import SessionLocal
from db.models import Template

import logging

router = APIRouter()


@router.post("/upload")
async def upload_template(
    institution_id: int = Form(...),
    file: UploadFile = File(...)
                        ):
    """
    Загружает шаблон для указанной организации. 
    Если для организации уже есть шаблон — он удаляется и заменяется новым.
    """
    db = SessionLocal()
    logger = logging.getLogger(__name__)
    try:
        # Удаляем старый шаблон, если существует
        existing = db.query(Template).filter_by(
            institution_id=institution_id).first()
        if existing:
            db.delete(existing)
            db.commit()

        # Чтение и сохранение нового шаблона
        content = await file.read()
        template = Template(name=file.filename,
                            content=content,
                            institution_id=institution_id)
        db.add(template)
        db.commit()
        logger.info("Template uploaded successfully for institution %s",
                    institution_id)
        return {"message": "Template uploaded for institution"}
    finally:
        db.close()


@router.delete("/template/{institution_id}")
def delete_template(institution_id: int):
    """
    Удаляет шаблон, привязанный к указанной организации.
    Если шаблон не найден, возвращает 404.
    """
    db = SessionLocal()
    logger = logging.getLogger(__name__)
    try:
        template = db.query(Template).filter_by(
                        institution_id=institution_id).first()
        if not template:
            logger.warning("Template not found for institution %s",
                           institution_id)
            raise HTTPException(status_code=404, detail="Template not found")
        db.delete(template)
        db.commit()
        return {"message": "Template deleted"}
    finally:
        db.close()


@router.get("/")
def list_templates():
    """
    Возвращает список всех шаблонов (id и имя) из базы данных.
    """
    return template_service.get_all_templates()
