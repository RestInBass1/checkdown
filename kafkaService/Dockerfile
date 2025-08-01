# 1) Build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o worker ./cmd/worker

# 2) Runtime
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/worker ./worker

# точку монтирования делаем только для логов
VOLUME ["/app/logs"]

CMD ["./worker"]
