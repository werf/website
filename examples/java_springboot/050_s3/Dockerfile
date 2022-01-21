FROM adoptopenjdk/openjdk11:jdk-11.0.11_9-alpine AS builder
ADD . /src
WORKDIR /src
RUN ./mvnw clean package -DskipTests

FROM adoptopenjdk/openjdk11:jdk-11.0.11_9-alpine as backend
RUN apk add -U mysql-client curl
WORKDIR /app
COPY --from=builder /src/target/*.jar ./app.jar
EXPOSE 8080

FROM nginx:stable-alpine as frontend
WORKDIR /www
ADD src/main/resources/static /www/static/
COPY .werf/nginx.conf /etc/nginx/nginx.conf
