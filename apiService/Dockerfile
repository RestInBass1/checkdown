# ---------- build stage ----------
FROM golang:1.21-alpine AS builder
WORKDIR /src

# 1) кешируем зависимости — только go.mod/go.sum
COPY go.mod go.sum ./
RUN go mod download

# 2) копируем весь проект
COPY . .

# 3) собираем конкретно api‑сервис
RUN CGO_ENABLED=0 go build -o /bin/apiService ./apiService/cmd

# ---------- runtime stage ----------
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /bin/apiService .

EXPOSE 8080                               # внутренний порт (для инфы)
ENTRYPOINT ["/app/apiService"]
