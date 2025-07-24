FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники auth-сервиса
COPY auth ./auth
COPY pkg ./pkg
COPY proto ./proto

# Собираем бинарник
RUN go build -o /app/auth-server ./auth/cmd/auth/main.go

FROM alpine:latest

WORKDIR /app

# Копируем бинарник из builder-слоя
COPY --from=builder /app/auth-server .

ENV JWT_SECRET=jwt_secret_key_SoftTradeIT

# Открываем порт (gRPC)
EXPOSE 8080

# Запускаем сервис
CMD ["./auth-server"]