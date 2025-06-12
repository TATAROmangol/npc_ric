import logging
from db.database import SessionLocal
from db.models import Template
from fastapi import HTTPException

logger = logging.getLogger(__name__)


def save_template(name: str, content: bytes):
    db = SessionLocal()
    try:
        logger.info(f"Attempting to save template with name: {name}")
        existing = db.query(Template).filter_by(name=name).first()
        if existing:
            logger.warning(f"Template with name '{name}' already exists")
            raise HTTPException(
                status_code=400,
                detail="Template with this name already exists")
        template = Template(name=name, content=content)
        db.add(template)
        db.commit()
        logger.info(f"Template '{name}' saved successfully (ID={template.id})")
    except Exception as e:
        logger.exception(f"Error while saving template '{name}': {e}")
        raise
    finally:
        db.close()


def get_all_templates():
    db = SessionLocal()
    try:
        templates = db.query(Template).all()
        logger.info(f"Retrieved {len(templates)} templates from database")
        return [{"id": t.id, "name": t.name} for t in templates]
    except Exception:
        logger.exception("Error while retrieving templates")
        raise
    finally:
        db.close()
