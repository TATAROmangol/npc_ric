from dotenv import load_dotenv
import os

load_dotenv()

PG_USER = os.getenv('PG_USER')
PG_PASS = os.getenv('PG_PASSWORD')
PG_HOST = 'postgres'
PG_PORT = os.getenv('PG_PORT')
DB_NAME = os.getenv('PG_DB_NAME')

DB_URL = f"postgresql://{PG_USER}:{PG_PASS}@{PG_HOST}:{PG_PORT}/{DB_NAME}"
