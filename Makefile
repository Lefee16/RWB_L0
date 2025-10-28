.PHONY: help docker-up docker-down migrate-up migrate-down migrate-create run test clean

# Показать все доступные команды
help:
	@echo "Доступные команды:"
	@echo "  make docker-up       - Запустить Docker контейнеры (Postgres + NATS)"
	@echo "  make docker-down     - Остановить Docker контейнеры"
	@echo "  make migrate-up      - Применить миграции БД"
	@echo "  make migrate-down    - Откатить миграции БД"
	@echo "  make run             - Запустить сервис"
	@echo "  make test            - Запустить тесты"
	@echo "  make clean           - Очистить Docker volumes"

# Запустить Docker контейнеры
docker-up:
	docker-compose up -d

# Остановить Docker контейнеры
docker-down:
	docker-compose down

# Применить миграции
migrate-up:
	migrate -path migrations -database "postgresql://orderuser:orderpass@localhost:5432/orders_db?sslmode=disable" up

# Откатить миграции
migrate-down:
	migrate -path migrations -database "postgresql://orderuser:orderpass@localhost:5432/orders_db?sslmode=disable" down

# Создать новую миграцию (использование: make migrate-create NAME=add_index)
migrate-create:
	migrate create -ext sql -dir migrations -seq $(NAME)

# Запустить сервис
run:
	go run cmd/main.go

# Запустить тесты
test:
	go test -v ./...

# Очистить Docker volumes
clean:
	docker-compose down -v
