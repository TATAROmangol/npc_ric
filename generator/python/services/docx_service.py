from python.db.models import Template, GeneratedDocument
from python.db.database import SessionLocal
from docxtpl import DocxTemplate
import io
# import json


class DocxGeneratorService:
    @staticmethod
    def generate_docx(template_name: str,
                      data: dict,
                      save_to_db: bool = False):
        db = SessionLocal()
        try:
            template = db.query(Template).filter_by(name=template_name).first()
            if not template:
                return None
            doc = DocxTemplate(io.BytesIO(template.content))
            doc.render(data)
            buffer = io.BytesIO()
            doc.save(buffer)
            file_content = buffer.getvalue()

            if save_to_db:
                new_doc = GeneratedDocument(
                    template_id=template.id,
                    file_content=file_content
                )
                db.add(new_doc)
                db.commit()

            return file_content
        finally:
            db.close()
