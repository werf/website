# [<en>] Here, a multi-stage build is used to assemble the image.
# [<ru>] Используем многоступенчатую сборку образа (multi-stage build).
# [<en>] The image for the project assembly.
# [<ru>] Образ, в котором будет собираться проект.
FROM eclipse-temurin:11-alpine AS builder
# [<en>] Copy the source code of the application.
# [<ru>] Копируем исходники приложения.
COPY . /src
WORKDIR /src
# [<en>] Run Maven to build the project.
# [<ru>] Запускаем Maven для сборки проекта.
RUN ./mvnw clean package

# [<en>] The image to deploy to the cluster.
# [<ru>] Образ, который будет выкатываться в кластер.
FROM eclipse-temurin:11-alpine as backend
WORKDIR /app
# [<en>] Copy the project executable from the build image.
# [<ru>] Копируем из сборочного образа исполняемый файл проекта.
COPY --from=builder /src/target/*.jar ./app.jar
EXPOSE 8080
ENTRYPOINT ["java","-jar","./app.jar"]

# [<en>] The NGINX image that contains the static files.
# [<ru>] NGINX-образ со статическими файлами.
FROM nginx:stable-alpine as frontend
WORKDIR /www
# [<en>] Copy static files.
# [<ru>] Копируем статические файлы.
COPY src/main/resources/static /www/static/
# [<en>] Copy NGINX configuration.
# [<ru>] Копируем конфигурацию NGINX.
COPY .werf/nginx.conf /etc/nginx/nginx.conf
