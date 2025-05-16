from fastapi import FastAPI
from routes.template_routes import router as template_router
from routes.document_routes import router as document_router

app = FastAPI()

app.include_router(template_router, prefix="/templates", tags=["Templates"])
app.include_router(document_router, prefix="/documents", tags=["Documents"])
