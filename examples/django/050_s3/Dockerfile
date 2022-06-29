# [<en>] Here, a multi-stage build is used to assemble the image.
# [<ru>] Используем многоступенчатую сборку образа (multi-stage build).
# [<en>] The image to deploy to the cluster.
# [<ru>] Образ, который будет выкатываться в кластер.
FROM python:3-alpine as backend
WORKDIR /app

RUN apk add python3-dev gcc libc-dev linux-headers musl-dev mariadb-connector-c-dev

# [<en>] Install the Django web framework.
# [<ru>] Устанавливаем Django веб-фреймворк.
RUN python -m pip install Django uwsgi mysqlclient boto3

# [<en>] Copy application files into the image.
# [<ru>] Копируем в образ файлы приложения.
COPY app .

# [<en>] The NGINX image that contains the static files.
# [<ru>] NGINX-образ со статическими файлами.
FROM nginx:stable-alpine as frontend
WORKDIR /app
# [<en>] Copy static files.
# [<ru>] Копируем статические файлы.
COPY app/static /app/static
# [<en>] Copy NGINX configuration.
# [<ru>] Копируем конфигурацию NGINX.
COPY .werf/nginx.conf /etc/nginx/nginx.conf
