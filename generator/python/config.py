from dotenv import load_dotenv
import os

load_dotenv()

DB_URL = f"postgresql://{os.getenv('PG_USER')}:{os.getenv('PG_PASSWORD')}@postgres:{os.getenv('PG_PORT')}/{os.getenv('PG_DB_NAME')}"
