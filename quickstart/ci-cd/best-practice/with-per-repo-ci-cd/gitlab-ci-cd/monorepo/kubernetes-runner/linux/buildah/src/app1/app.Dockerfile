FROM alpine:3.17 as builder
WORKDIR /app
RUN apk add go-task
COPY . .
RUN go-task build

FROM alpine:3.17
WORKDIR /app
RUN apk add nmap-ncat go-task
COPY --from=builder /app/app.sh /app/Taskfile.yaml ./
CMD ["go-task", "run"]
