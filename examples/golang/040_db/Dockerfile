# [<en>] Here, a multi-stage build is used to assemble the image
# [<ru>] Используем многоступенчатую сборку образа (multi-stage build)
# [<en>] The image to build the project
# [<ru>] Образ, в котором будет собираться проект
FROM golang:1.18-alpine AS build
# [<en>] Устанавливаем curl и tar.
# [<ru>] Устанавливаем curl и tar.
RUN apk add curl tar
# [<en>] Copy the source code of the application.
# [<ru>] Копируем исходники приложения
COPY . /app
WORKDIR /app
# [<en>] Скачиваем утилиту migrate и распаковываем полученный архив.
# [<ru>] Скачиваем утилиту migrate и распаковываем полученный архив.
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
# [<en>] Запускаем загрузку нужных пакетов.
# [<ru>] Запускаем загрузку нужных пакетов.
RUN go mod download
# [<en>] Запускаем сборку приложения.
# [<ru>] Запускаем сборку приложения.
RUN go build -o /goapp cmd/main.go

# [<en>] The image to deploy to the cluster.
# [<ru>] Образ, который будет разворачиваться в кластере.
FROM alpine:latest as backend
WORKDIR /
# [<en>] Copy the project executable from the build image.
# [<ru>] Копируем из сборочного образа исполняемый файл проекта.
COPY --from=build /goapp /goapp
# [<en>] Копируем из сборочного образа распакованный файл утилиты migrate и схемы миграции.
# [<ru>] Копируем из сборочного образа распакованный файл утилиты migrate и схемы миграции.
COPY --from=build /app/migrate /migrations/migrate
COPY db/migrations /migrations/schemes
# [<en>] Копируем файлы ассетов и шаблоны.
# [<ru>] Копируем файлы ассетов и шаблоны.
COPY ./templates /templates
COPY ./static /static
EXPOSE 8080
ENTRYPOINT ["/goapp"]

# [<en>] The NGINX image that contains the static files.
# [<ru>] NGINX-образ со статическими файлами.
FROM nginx:stable-alpine as frontend
WORKDIR /www
# [<en>] Copy static files.
# [<ru>] Копируем статические файлы.
COPY static /www/static/
# [<en>] Copy NGINX configuration.
# [<ru>] Копируем конфигурацию NGINX.
COPY .werf/nginx.conf /etc/nginx/nginx.conf