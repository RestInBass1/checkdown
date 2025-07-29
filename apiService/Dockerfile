# ---------- build stage ----------
FROM golang:1.24-alpine AS builder
WORKDIR /src

# 1) копируем только манифест (go.mod/go.sum) модуля apiService
COPY go.mod go.sum ./
WORKDIR /src/apiService
RUN go mod download

# 2) копируем исходники именно apiService
COPY . .

# 3) собираем бинарь
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/apiService ./cmd

# ---------- runtime stage ----------
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /bin/apiService .

EXPOSE 8080
ENTRYPOINT ["/app/apiService"]
