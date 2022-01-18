FROM eclipse-temurin:11-alpine AS builder
ADD . /src
WORKDIR /src
RUN ./mvnw clean package

FROM eclipse-temurin:11-alpine
WORKDIR /app
COPY --from=builder /src/target/*.jar ./app.jar
EXPOSE 8080
ENTRYPOINT ["java","-jar","./app.jar"]
