# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /src

# 1) Копируем манифесты и вендор
COPY go.mod go.sum vendor/ ./

# 2) Копируем весь исходный код
COPY . .

# 3) Собираем с флагом -mod=vendor (не нужна сеть и внешний доступ)
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -mod=vendor -o apiService ./cmd

# ── Runtime Stage ────────────────────────────────────────────────────────────
FROM alpine:3.19
WORKDIR /app

# Копируем бинарь
COPY --from=builder /src/apiService ./

# Копируем локальный .env (если нужно)
COPY .env ./.env

EXPOSE 8080
ENTRYPOINT ["./apiService"]


