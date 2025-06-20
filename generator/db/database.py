from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from config import DB_URL

# Создание движка SQLAlchemy
engine = create_engine(DB_URL)

# Создание фабрики сессий для взаимодействия с БД
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def init_db():
    """
    Создаёт таблицы в базе данных согласно моделям, если они ещё не существуют.
    Вызывается при старте приложения.
    """
    from db.models import Base
    Base.metadata.create_all(bind=engine)
