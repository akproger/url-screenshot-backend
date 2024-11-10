# Базовый образ Go
FROM golang:1.18 AS builder

WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Установка зависимостей
RUN go mod download

# Копируем все файлы проекта
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Минимальный образ для продакшена
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
