FROM ubuntu:20.04
WORKDIR /app

# [<en>] Install system dependencies
# [<ru>] Уставливаем системные зависимости
RUN apt-get -o Acquire::Check-Valid-Until=false -o Acquire::Check-Date=false update -qq && \
    apt-get install -y curl ncat

# [<en>] Add all the files from our repo to the image, including app source code
# [<ru>] Добавляем в образ все остальные файлы из нашего репозитория, включая исходный код приложения
COPY . .
