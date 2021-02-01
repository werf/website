FROM golang:1.15-alpine as builder
RUN apk add --update gcc musl-dev
WORKDIR /app
COPY go.* /app/
RUN go mod download
COPY /cmd /app/cmd
RUN go build ./cmd/demoapp

FROM alpine:3.13
COPY --from=builder /app/demoapp /app/demoapp
COPY /app.db /app/app.db

CMD ["/app/demoapp"]
