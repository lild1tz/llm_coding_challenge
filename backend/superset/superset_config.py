# backend/superset/superset_config.py

import os

# Секретный ключ для Flask-сессий и CSRF
SECRET_KEY = os.environ.get("SUPERSET_SECRET_KEY")
if not SECRET_KEY:
    raise RuntimeError("Не задана переменная окружения SUPERSET_SECRET_KEY")

# Подключение к БД метаданных Superset
SQLALCHEMY_DATABASE_URI = os.environ.get("SUPERSET_SQLALCHEMY_DATABASE_URI")
if not SQLALCHEMY_DATABASE_URI:
    raise RuntimeError("Не задана переменная окружения SUPERSET_SQLALCHEMY_DATABASE_URI")

# (Опционально) дополнительные настройки:
# - Включение React CRUD-views
FEATURE_FLAGS = {
    "LISTVIEWS_DEFAULT_CARD_VIEW": True,
}

# Место хранения пользовательских артефактов
SUPERSET_HOME = os.environ.get("SUPERSET_HOME", "/app/superset_home")
