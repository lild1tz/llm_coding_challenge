FROM apache/superset

COPY requirements.txt /app/requirements.txt

RUN pip install -r /app/requirements.txt

EXPOSE 8088

RUN superset db upgrade
RUN superset fab create-admin --username ${SUPERSET_ADMIN_USERNAME} --firstname ${SUPERSET_ADMIN_FIRSTNAME} --lastname ${SUPERSET_ADMIN_LASTNAME} --email ${SUPERSET_ADMIN_EMAIL} --password ${SUPERSET_ADMIN_PASSWORD}
RUN superset init

CMD ["superset", "run", "-h", "0.0.0.0", "-p", "8088"]
