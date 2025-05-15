from fastapi import APIRouter, UploadFile, File, Form, HTTPException
from services import template_service

router = APIRouter()


@router.post("/upload")
async def upload_template(name: str = Form(...), file: UploadFile = File(...)):
    try:
        content = await file.read()
        template_service.save_template(name, content)
        return {"message": "Template uploaded to DB"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/")
def list_templates():
    return template_service.get_all_templates()
