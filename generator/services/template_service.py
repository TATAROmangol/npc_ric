from db.database import SessionLocal
from db.models import Template
from fastapi import HTTPException


def save_template(name: str, content: bytes):
    db = SessionLocal()
    try:
        existing = db.query(Template).filter_by(name=name).first()
        if existing:
            raise HTTPException(
                status_code=400,
                detail="Template with this name already exists")
        template = Template(name=name, content=content)
        db.add(template)
        db.commit()
    finally:
        db.close()


def get_all_templates():
    db = SessionLocal()
    try:
        return [{"id": t.id, "name": t.name} for t in db.query(Template).all()]
    finally:
        db.close()
