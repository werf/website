FROM python:3-alpine as backend
WORKDIR /app

RUN apk add python3-dev gcc libc-dev linux-headers musl-dev mariadb-connector-c-dev mariadb-dev

RUN python -m pip install Django uwsgi mysqlclient boto3

COPY app .

# [<snippet frontend>]
FROM nginx:stable-alpine as frontend
WORKDIR /app
COPY app/static /app/static
# [<endsnippet frontend>]
