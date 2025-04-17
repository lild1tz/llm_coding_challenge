# Используем официальный образ Superset
FROM apache/superset:latest

# Переходим на root, чтобы скопировать свои файлы и установить зависимости
USER root

# Копируем свой конфиг и список зависимостей
COPY superset_config.py /app/pythonpath/superset_config.py
COPY requirements.txt  /app/requirements.txt

# Создаём HOME-директорию и устанавливаем права, затем ставим pip‑пакеты
RUN mkdir -p /app/superset_home \
    && chown superset:superset /app/superset_home \
    && pip install --no-cache-dir -r /app/requirements.txt

# Возвращаемся к непривилегированному пользователю
USER superset

# Указываем Superset, где искать конфиг, и где хранить HOME
ENV SUPERSET_CONFIG_PATH=/app/pythonpath/superset_config.py \
    SUPERSET_HOME=/app/superset_home

EXPOSE 8088

# Shell‑форма CMD: выполнит миграции, создаст админа (если нужно), и запустит сервер
CMD superset db upgrade && \
    superset fab create-admin \
      --username "$SUPERSET_ADMIN_USERNAME" \
      --firstname "$SUPERSET_ADMIN_FIRSTNAME" \
      --lastname "$SUPERSET_ADMIN_LASTNAME" \
      --email "$SUPERSET_ADMIN_EMAIL" \
      --password "$SUPERSET_ADMIN_PASSWORD" && \
    superset init && \
    exec superset run -h 0.0.0.0 -p 8088
