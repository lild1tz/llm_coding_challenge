import os

# ваш уже заданный секретный ключ
SECRET_KEY = "z7FjK2pQ9xY0e4WbR8vNl3sD5mU1tA6y"

# здесь указываем адрес вашей БД Postgres для метаданных Superset
SQLALCHEMY_DATABASE_URI = (
    os.environ.get("SUPERSET_SQLALCHEMY_DATABASE_URI")
    or "postgresql://hermes:password@postgres:5432/hermes?sslmode=disable"
)
