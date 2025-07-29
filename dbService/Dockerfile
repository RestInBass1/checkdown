FROM golang:1.24-alpine AS builder
WORKDIR /src

# в контексте есть go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# весь код сервиса
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/dbService ./cmd

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /bin/dbService .
EXPOSE 50051
ENTRYPOINT ["/app/dbService"]
