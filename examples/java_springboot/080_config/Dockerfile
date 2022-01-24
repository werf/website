FROM eclipse-temurin:11-alpine AS builder
COPY . /src
WORKDIR /src
RUN ./mvnw clean package -DskipTests

FROM eclipse-temurin:11-alpine as backend
RUN apk add -U mysql-client curl
WORKDIR /app
COPY --from=builder /src/target/*.jar ./app.jar
EXPOSE 8080

FROM nginx:stable-alpine as frontend
WORKDIR /www
ADD src/main/resources/static /www/static/
