FROM adoptopenjdk/openjdk11:jdk-11.0.11_9-alpine AS builder
ADD . /src
WORKDIR /src
RUN ./mvnw clean package

FROM adoptopenjdk/openjdk11:jdk-11.0.11_9-alpine
WORKDIR /app
COPY --from=builder /src/target/*.jar ./app.jar
EXPOSE 8080
ENTRYPOINT ["java","-jar","./app.jar"]
