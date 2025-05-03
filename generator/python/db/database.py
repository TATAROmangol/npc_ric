from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from python.config import DB_URL

engine = create_engine(DB_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def init_db():
    from db.models import Base
    Base.metadata.create_all(bind=engine)
