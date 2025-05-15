from fastapi import FastAPI, HTTPException, UploadFile, File, Form
from fastapi.responses import StreamingResponse
from pydantic import BaseModel
import io
import uuid
from db.database import SessionLocal
from db.models import Template, GeneratedDocument
from services.docx_service import generate_docx_from_template

app = FastAPI()


class GenerateRequest(BaseModel):
    template_name: str
    institution_id: int


@app.post("/generate")
def generate(request: GenerateRequest):
    db = SessionLocal()
    try:
        template = db.query(Template).filter_by(
            name=request.template_name).first()
        if not template:
            raise HTTPException(status_code=404, detail="Template not found")

        generated_bytes = generate_docx_from_template(template.content,
                                                      request.institution_id)

        doc_id = str(uuid.uuid4())
        db_doc = GeneratedDocument(
            template_id=template.id,
            file_content=generated_bytes,
            filename=f"{doc_id}.docx"
        )
        db.add(db_doc)
        db.commit()

        return {"download_url": f"/download/{db_doc.id}"}
    finally:
        db.close()


@app.get("/download/{doc_id}")
def download(doc_id: int):
    db = SessionLocal()
    try:
        doc = db.query(GeneratedDocument).filter_by(id=doc_id).first()
        if not doc:
            raise HTTPException(status_code=404,
                                detail="Generated file not found")

        return StreamingResponse(
            io.BytesIO(doc.file_content),
            media_type="application/vnd.openxmlformats-officedocument.wordprocessingml.document",
            headers={
                "Content-Disposition": f"attachment; filename={doc.filename}"}
        )
    finally:
        db.close()


@app.post("/upload-template/")
async def upload_template(name: str = Form(...), file: UploadFile = File(...)):
    db = SessionLocal()
    try:
        content = await file.read()
        template = Template(name=name, content=content)
        db.add(template)
        db.commit()
        return {"message": "Template uploaded to DB"}
    finally:
        db.close()
