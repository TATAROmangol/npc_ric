import os
from db.database import init_db, SessionLocal
from db.models import Template
from services.docx_service import DocxGeneratorService


def insert_template(template_path: str, name: str):
    with open(template_path, "rb") as f:
        content = f.read()

    db = SessionLocal()
    existing = db.query(Template).filter_by(name=name).first()
    if not existing:
        db.add(Template(name=name, content=content))
        db.commit()
        print(f"Шаблон '{name}' добавлен в базу.")
    else:
        print(f"Шаблон '{name}' уже есть в базе.")
    db.close()


if __name__ == "__main__":
    init_db()

    TEMPLATE_NAME = "contractct"
    TEMPLATE_PATH = "python/contract.docx"

    if not os.path.exists(TEMPLATE_PATH):
        print(f"Шаблон не найден: {TEMPLATE_PATH}")
        exit(1)

    insert_template(TEMPLATE_PATH, TEMPLATE_NAME)

    context = {
        "org_name": "ООО",
        "org_phone": "+7 (495) 123-45-67",
        "org_address": "г. Москва, ул. Примерная, д. 1",
        "org_inn": "7701234567",
        "students": [
            "Иванов Иван Иванович",
            "Петров Пётр Петрович",
            "Сидоров Сидор Сидорович"
        ]
    }

    result = DocxGeneratorService.generate_docx(TEMPLATE_NAME,
                                                context,
                                                save_to_db=True)
    if result:
        with open("output_test.docx", "wb") as f:
            f.write(result)
        print("Файл сгенерирован: output_test.docx")
    else:
        print("Не удалось найти шаблон в базе.")
