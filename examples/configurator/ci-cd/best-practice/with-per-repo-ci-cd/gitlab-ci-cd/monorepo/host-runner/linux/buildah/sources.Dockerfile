FROM alpine:3.17
WORKDIR /src
RUN apk add go-task
COPY . .
