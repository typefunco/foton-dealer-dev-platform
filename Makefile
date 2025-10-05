.PHONY: help
help: ## Показать это сообщение помощи
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: up
up: ## Запустить все сервисы
	docker-compose up -d

.PHONY: down
down: ## Остановить все сервисы
	docker-compose down

.PHONY: down-v
down-v: ## Остановить все сервисы и удалить volumes
	docker-compose down -v

.PHONY: build
build: ## Собрать все образы
	docker-compose build

.PHONY: logs
logs: ## Показать логи всех сервисов
	docker-compose logs -f

.PHONY: logs-backend
logs-backend: ## Показать логи backend
	docker-compose logs -f backend

.PHONY: logs-db
logs-db: ## Показать логи database
	docker-compose logs -f postgres

.PHONY: restart
restart: down up ## Перезапустить все сервисы

.PHONY: restart-backend
restart-backend: ## Перезапустить только backend
	docker-compose restart backend

.PHONY: ps
ps: ## Показать статус сервисов
	docker-compose ps

# Backend команды
.PHONY: backend-test
backend-test: ## Запустить тесты backend
	cd backend && go test -v -race ./...

.PHONY: backend-lint
backend-lint: ## Запустить линтер backend
	cd backend && golangci-lint run

.PHONY: backend-fmt
backend-fmt: ## Форматировать код backend
	cd backend && go fmt ./...

.PHONY: backend-vet
backend-vet: ## Проверить код backend
	cd backend && go vet ./...

.PHONY: backend-build
backend-build: ## Собрать бинарник backend
	cd backend && go build -o bin/dealer-platform ./cmd/app/main.go

.PHONY: backend-run
backend-run: ## Запустить backend локально
	cd backend && go run ./cmd/app/main.go

.PHONY: backend-install-tools
backend-install-tools: ## Установить инструменты для backend
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

# Frontend команды
.PHONY: frontend-install
frontend-install: ## Установить зависимости frontend
	cd frontend && npm ci

.PHONY: frontend-dev
frontend-dev: ## Запустить frontend в dev режиме
	cd frontend && npm run dev

.PHONY: frontend-build
frontend-build: ## Собрать frontend для production
	cd frontend && npm run build

.PHONY: frontend-lint
frontend-lint: ## Запустить линтер frontend
	cd frontend && npm run lint

.PHONY: frontend-preview
frontend-preview: ## Предпросмотр production сборки frontend
	cd frontend && npm run preview

# Database команды
.PHONY: db-migrate
db-migrate: ## Применить миграции к базе данных
	docker-compose exec postgres psql -U postgres -d dealer_platform -f /docker-entrypoint-initdb.d/001_create_initial_tables.sql

.PHONY: db-shell
db-shell: ## Открыть psql shell
	docker-compose exec postgres psql -U postgres -d dealer_platform

.PHONY: db-reset
db-reset: down-v up ## Сбросить базу данных (удалить и пересоздать)

# CI/CD команды
.PHONY: ci-local
ci-local: backend-lint backend-test frontend-lint frontend-build ## Запустить все проверки локально

.PHONY: ci-backend
ci-backend: backend-fmt backend-vet backend-lint backend-test ## Запустить все backend проверки

.PHONY: ci-frontend
ci-frontend: frontend-lint frontend-build ## Запустить все frontend проверки

# Docker команды
.PHONY: docker-prune
docker-prune: ## Очистить неиспользуемые Docker ресурсы
	docker system prune -af --volumes

.PHONY: docker-backend-shell
docker-backend-shell: ## Открыть shell в backend контейнере
	docker-compose exec backend sh

# Git команды
.PHONY: git-clean
git-clean: ## Очистить неотслеживаемые файлы
	git clean -fd

# Полная проверка перед commit
.PHONY: pre-commit
pre-commit: backend-fmt backend-vet backend-lint backend-test frontend-lint ## Полная проверка перед коммитом
	@echo "\n✅ Все проверки пройдены! Можно делать commit."

# Инициализация проекта
.PHONY: init
init: frontend-install backend-install-tools up ## Инициализировать проект (первый запуск)
	@echo "\n✅ Проект инициализирован!"
	@echo "Backend доступен по адресу: http://localhost:8080"
	@echo "Frontend: cd frontend && npm run dev"

.DEFAULT_GOAL := help

