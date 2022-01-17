# [<ru>] Используем многоступенчатую сборку образа (multi-stage build)
# [<ru>] Образ, в котором будет собираться проект
FROM eclipse-temurin:11-alpine AS builder
# [<ru>] Копируем исходники приложения
COPY . /src
WORKDIR /src
# [<ru>] Запускаем Maven для сборки проекта
RUN ./mvnw clean package

# [<ru>] Образ, который будет выкатываться в кластер
FROM eclipse-temurin:11-alpine
WORKDIR /app
# [<ru>] Копируем из сборочного образа исполняемый файл проекта
COPY --from=builder /src/target/*.jar ./app.jar
EXPOSE 8080
ENTRYPOINT ["java","-jar","./app.jar"]
