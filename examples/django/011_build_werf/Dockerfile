FROM python:3.9.1-alpine3.12
WORKDIR /app
COPY project /app

RUN pip3 install -r requirements.txt

ENTRYPOINT ["python3", "manage.py"]
CMD ["runserver"]
