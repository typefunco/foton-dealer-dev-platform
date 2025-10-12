# Makefile для управления проектом

.PHONY: help build test test-integration test-unit clean run migrate-up migrate-down docker-build docker-run

# Переменные
BINARY_NAME=dealer_dev_platform
DOCKER_IMAGE=dealer-dev-platform
DOCKER_TAG=latest

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  build           - Собрать приложение"
	@echo "  test            - Запустить все тесты"
	@echo "  test-integration - Запустить интеграционные тесты"
	@echo "  test-unit       - Запустить unit тесты"
	@echo "  clean           - Очистить артефакты сборки"
	@echo "  run             - Запустить приложение"
	@echo "  migrate-up      - Применить миграции"
	@echo "  migrate-down    - Откатить миграции"
	@echo "  docker-build    - Собрать Docker образ"
	@echo "  docker-run      - Запустить в Docker"

# Сборка приложения
build:
	@echo "Сборка приложения..."
	go build -o bin/$(BINARY_NAME) ./cmd/app

# Запуск всех тестов
test:
	@echo "Запуск всех тестов..."
	go test -v ./...

# Интеграционные тесты
test-integration:
	@echo "Запуск интеграционных тестов..."
	go test -v -tags=integration ./internal/service/...

# Unit тесты
test-unit:
	@echo "Запуск unit тестов..."
	go test -v -short ./...

# Очистка
clean:
	@echo "Очистка артефактов..."
	rm -rf bin/
	go clean

# Запуск приложения
run:
	@echo "Запуск приложения..."
	go run ./cmd/app

# Применение миграций
migrate-up:
	@echo "Применение миграций..."
	# Здесь можно добавить команды для миграций

# Откат миграций
migrate-down:
	@echo "Откат миграций..."
	# Здесь можно добавить команды для отката миграций

# Сборка Docker образа
docker-build:
	@echo "Сборка Docker образа..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Запуск в Docker
docker-run:
	@echo "Запуск в Docker..."
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Установка зависимостей
deps:
	@echo "Установка зависимостей..."
	go mod download
	go mod tidy

# Проверка кода
lint:
	@echo "Проверка кода..."
	golangci-lint run

# Форматирование кода
fmt:
	@echo "Форматирование кода..."
	go fmt ./...

# Проверка безопасности
security:
	@echo "Проверка безопасности..."
	gosec ./...

# Генерация документации
docs:
	@echo "Генерация документации..."
	godoc -http=:6060

# Покрытие тестами
coverage:
	@echo "Анализ покрытия тестами..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Отчет о покрытии сохранен в coverage.html"

# Бенчмарки
bench:
	@echo "Запуск бенчмарков..."
	go test -bench=. -benchmem ./...

# Профилирование
profile:
	@echo "Профилирование..."
	go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...

# Проверка зависимостей
check-deps:
	@echo "Проверка зависимостей..."
	go mod verify
	go list -m all

# Обновление зависимостей
update-deps:
	@echo "Обновление зависимостей..."
	go get -u ./...
	go mod tidy
