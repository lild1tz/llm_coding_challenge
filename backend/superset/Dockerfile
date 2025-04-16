# backend/superset/Dockerfile
FROM apache/superset:latest

# 1) Копируем ваш конфиг
COPY superset_config.py /app/
ENV SUPERSET_CONFIG_PATH=/app/superset_config.py
# (опционально) чтобы том с данными тоже использовался как HOME
ENV SUPERSET_HOME=/app/superset_home

# 2) Устанавливаем зависимости
COPY requirements.txt /app/requirements.txt
RUN pip install -r /app/requirements.txt

EXPOSE 8088

# 3) При старте: миграции в Postgres, создание админа, инициализация и запуск сервера
CMD ["bash", "-lc", "\
    superset db upgrade && \
    superset fab create-admin \
      --username $SUPERSET_ADMIN_USERNAME \
      --firstname $SUPERSET_ADMIN_FIRSTNAME \
      --lastname $SUPERSET_ADMIN_LASTNAME \
      --email $SUPERSET_ADMIN_EMAIL \
      --password $SUPERSET_ADMIN_PASSWORD && \
    superset init && \
    exec superset run -h 0.0.0.0 -p 8088"]
