FROM python:3-alpine
WORKDIR /app

RUN apk add python3-dev gcc libc-dev linux-headers

# [<en>] Install the Django web framework.
# [<ru>] Устанавливаем Django веб-фреймворк.
RUN python -m pip install Django uwsgi

# [<en>] Copy application files into the image.
# [<ru>] Копируем в образ файлы приложения.
COPY app .
