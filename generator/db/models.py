from sqlalchemy import Column, Integer, String, LargeBinary
from sqlalchemy.orm import declarative_base

Base = declarative_base()


# Модель шаблона документа  
class Template(Base):
    __tablename__ = "templates"
    id = Column(Integer, primary_key=True)
    name = Column(String)
    content = Column(LargeBinary)
    institution_id = Column(Integer, unique=True)


# Модель сгенерированного документа
class GeneratedDocument(Base):
    __tablename__ = "generated_docs"
    id = Column(Integer, primary_key=True)
    template_id = Column(Integer)
    file_content = Column(LargeBinary)
    filename = Column(String)
