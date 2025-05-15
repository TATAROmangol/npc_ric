from sqlalchemy import Column, Integer, String, LargeBinary
from sqlalchemy.orm import declarative_base

Base = declarative_base()


class Template(Base):
    __tablename__ = "templates"
    id = Column(Integer, primary_key=True)
    name = Column(String, unique=True)
    content = Column(LargeBinary)


class GeneratedDocument(Base):
    __tablename__ = "generated_docs"
    id = Column(Integer, primary_key=True)
    template_id = Column(Integer)
    file_content = Column(LargeBinary)
    filename = Column(String)
