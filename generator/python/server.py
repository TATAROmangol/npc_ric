from fastapi import FastAPI, HTTPException, UploadFile, File, Form
from fastapi.responses import FileResponse
from pydantic import BaseModel
import uuid
from db.database import SessionLocal
from db.models import Template, GeneratedDocument
from services.docx_service import generate_docx_from_template
from grpc_clients.table_client import get_table_data
import shutil

app = FastAPI()


class GenerateRequest(BaseModel):
    template_name: str
    institution_id: int


@app.post("/generate")
def generate(request: GenerateRequest):
    db = SessionLocal()
    try:
        template = db.query(Template).filter_by(name=request.template_name).first()
        if not template:
            raise HTTPException(status_code=404, detail="Template not found")
        tmp_path = f"/tmp/{uuid.uuid4()}.docx"
        with open(tmp_path, "wb") as f:
            f.write(template.content)

        output_path = f"/tmp/{uuid.uuid4()}_gen.docx"
        generate_docx_from_template(tmp_path, output_path, request.institution_id)

        with open(output_path, "rb") as f:
            generated_bytes = f.read()

        doc_id = str(uuid.uuid4())
        db_doc = GeneratedDocument(template_id=template.id, file_content=generated_bytes, filename=f"{doc_id}.docx")
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
            raise HTTPException(status_code=404, detail="Generated file not found")

        return FileResponse(
            path_or_file=doc.file_content,
            media_type="application/vnd.openxmlformats-officedocument.wordprocessingml.document",
            filename=doc.filename,
            headers={"Content-Disposition": f"attachment; filename={doc.filename}"}
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

