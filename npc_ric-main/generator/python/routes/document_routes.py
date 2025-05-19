from fastapi import APIRouter, HTTPException
from fastapi.responses import StreamingResponse
from pydantic import BaseModel
import io
from services import document_service

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
    doc = document_service.get_document_by_id(doc_id)
    if not doc:
        raise HTTPException(status_code=404, detail="Generated file not found")

    return StreamingResponse(
        io.BytesIO(doc["content"]),
        media_type="application/vnd.openxmlformats-officedocument.wordprocessingml.document",
        headers={"Content-Disposition": f"attachment; filename={doc['filename']}"}
    )
