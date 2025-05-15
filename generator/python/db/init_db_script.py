# db/init_db_script.py
import sys
import os
from db.database import init_db
# Добавляем путь к корню проекта (где config.py)
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

if __name__ == "__main__":
    init_db()
    print("✅ База данных инициализирована.")
