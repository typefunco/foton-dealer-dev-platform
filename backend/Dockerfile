# Многоступенчатая сборка для минимизации размера образа

# Этап 1: Сборка приложения
FROM golang:1.25-alpine AS builder

# Установка необходимых зависимостей
RUN apk add --no-cache git

WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app/main.go

# Этап 2: Создание минимального образа
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем скомпилированное приложение из builder
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]

