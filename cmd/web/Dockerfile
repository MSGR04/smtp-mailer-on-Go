# Используем официальный образ Go для сборки
FROM golang:1.21 AS builder

# Задаём рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum, затем скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл для веб‑приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o smtp-mailer ./cmd/web/main.go

# Финальный образ: минимальный образ с бинарником
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/smtp-mailer .
COPY templates ./templates
EXPOSE 8080
CMD ["/app/smtp-mailer"]
