FROM golang:1.18-alpine AS build
COPY . /app
WORKDIR /app
RUN go mod download
RUN go build -o /goapp app/main.go

FROM alpine:latest
WORKDIR /
COPY --from=build /goapp /goapp
EXPOSE 8080
ENTRYPOINT ["/goapp"]