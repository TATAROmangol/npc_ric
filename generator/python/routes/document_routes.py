from fastapi import APIRouter, HTTPException
from fastapi.responses import StreamingResponse
from pydantic import BaseModel
import io
from services import document_service
from db.database import SessionLocal
from db.models import GeneratedDocument



router = APIRouter()


class GenerateRequest(BaseModel):
    template_name: str
    institution_id: int


@router.post("/generate")
def generate(request: GenerateRequest):
    try:
        doc = document_service.generate_document(request.template_name,
                                                 request.institution_id)
        return {"download_url": f"/documents/download/{doc['id']}"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))



@router.get("/download/{doc_id}")
def download(doc_id: int):
    db = SessionLocal()
    try:
        doc = db.query(GeneratedDocument).filter_by(id=doc_id).first()
        if not doc:
            raise HTTPException(status_code=404, detail="Generated file not found")

        return StreamingResponse(
            io.BytesIO(doc.file_content),
            media_type="application/vnd.openxmlformats-officedocument.wordprocessingml.document",
            headers={
                "Content-Disposition": f"attachment; filename={doc.filename}"
            }
        )
    finally:
        db.close()
