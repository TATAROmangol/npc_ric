from fastapi import FastAPI
from routes.template_routes import router as template_router
from routes.document_routes import router as document_router
from db.database import init_db
from fastapi.middleware.cors import CORSMiddleware

from logger import setup_logger
setup_logger()


app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.on_event("startup")
def on_startup():
    init_db()


app.include_router(template_router, prefix="/templates", tags=["Templates"])
app.include_router(document_router, prefix="/documents", tags=["Documents"])
