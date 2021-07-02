FROM alpine:3.14
WORKDIR /app

# [<en>] Install app dependencies
# [<ru>] Уставливаем зависимости приложения
RUN apk add --no-cache --update nmap-ncat

# [<en>] Add to the image a script to run the echo server and set the permission to execution
# [<ru>] Добавляем в образ скрипт для запуска эхо-сервера и устанавливаем разрешение на выполнение
COPY start.sh .
RUN chmod +x start.sh
