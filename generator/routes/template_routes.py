from fastapi import APIRouter, UploadFile, File, Form, HTTPException
from services import template_service
from db.database import SessionLocal
from db.models import Template

router = APIRouter()


@router.post("/upload")
async def upload_template(
    institution_id: int = Form(...),
    file: UploadFile = File(...)
                        ):
    db = SessionLocal()
    try:
        existing = db.query(Template).filter_by(
            institution_id=institution_id).first()
        if existing:
            db.delete(existing)
            db.commit()

        content = await file.read()
        template = Template(name=file.filename,
                            content=content,
                            institution_id=institution_id)
        db.add(template)
        db.commit()
        return {"message": "Template uploaded for institution"}
    finally:
        db.close()


@router.delete("/template/{institution_id}")
def delete_template(institution_id: int):
    db = SessionLocal()
    try:
        template = db.query(Template).filter_by(
                        institution_id=institution_id).first()
        if not template:
            raise HTTPException(status_code=404, detail="Template not found")
        db.delete(template)
        db.commit()
        return {"message": "Template deleted"}
    finally:
        db.close()


@router.get("/")
def list_templates():
    return template_service.get_all_templates()
