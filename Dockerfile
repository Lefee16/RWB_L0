# ===================================
# Стадия 1: Кэширование модулей
# ===================================
FROM golang:1.25-alpine AS modules

# Копируем только go.mod и go.sum для кэширования
COPY go.mod go.sum /modules/

WORKDIR /modules

# Скачиваем зависимости (этот слой будет закэширован)
RUN go mod download

# ===================================
# Стадия 2: Сборка приложения
# ===================================
FROM golang:1.25-alpine AS builder

# Копируем закэшированные модули
COPY --from=modules /go/pkg /go/pkg

# Копируем весь проект
COPY . /app

WORKDIR /app

# Собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/main.go

# ===================================
# Стадия 3: Финальный образ (минимальный)
# ===================================
FROM alpine:latest AS final

# Устанавливаем сертификаты для HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /bin/app .

# Копируем web файлы (templates и static)
COPY --from=builder /app/web ./web

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./app"]
