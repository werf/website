FROM golang:1.15-alpine as builder
RUN apk add --update gcc musl-dev
WORKDIR /app
COPY . .
RUN go build ./cmd/demoapp

FROM alpine:3.13
COPY --from=builder /app/demoapp /app/demoapp
COPY /app/app.db /app/app.db

CMD ["/app/demoapp"]
