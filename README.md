# Dealer Development Platform

Платформа для анализа и управления дилерами с интеграцией frontend и backend через Docker.

## Архитектура

- **Frontend**: React + TypeScript + Vite + Tailwind CSS
- **Backend**: Go + Echo + PostgreSQL
- **Database**: PostgreSQL
- **Containerization**: Docker + Docker Compose

## Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Go 1.21+ (для локальной разработки)
- Node.js 18+ (для локальной разработки)

### Запуск через Docker Compose

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd dealer_dev_platform
```

2. Запустите все сервисы:
```bash
docker-compose up --build
```

3. Откройте приложение:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Database: localhost:5432

### Структура сервисов

- **postgres**: База данных PostgreSQL
- **backend**: Go API сервер
- **frontend**: React приложение с Nginx

### Переменные окружения

#### Backend
- `DATABASE_URL`: Строка подключения к PostgreSQL
- `JWT_SECRET`: Секретный ключ для JWT токенов
- `SERVER_PORT`: Порт сервера (по умолчанию 8080)

#### Frontend
- `VITE_API_BASE_URL`: URL API бэкенда (по умолчанию http://backend:8080/api)

### API Endpoints

#### Пользователи
- `GET /api/users` - Получить список пользователей
- `GET /api/users/:id` - Получить пользователя по ID
- `POST /api/users` - Создать пользователя
- `PUT /api/users/:id` - Обновить пользователя
- `DELETE /api/users/:id` - Удалить пользователя
- `GET /api/users/stats` - Статистика пользователей

#### Дилеры
- `GET /api/dealers` - Получить список дилеров
- `GET /api/dealers/:id` - Получить дилера по ID
- `GET /api/dealers/:id/card` - Получить карточку дилера

#### Производительность
- `GET /api/performance` - Данные производительности дилеров

#### After Sales
- `GET /api/aftersales` - Данные After Sales

#### Команда продаж
- `GET /api/sales` - Данные команды продаж

#### Сравнение кварталов
- `GET /api/quarter-comparison` - Сравнение кварталов

#### Все данные
- `GET /api/all-data` - Комплексные данные всех таблиц

### Разработка

#### Локальная разработка Frontend

1. Перейдите в папку frontend:
```bash
cd frontend
```

2. Установите зависимости:
```bash
npm install
# или
yarn install
```

3. Создайте файл `.env`:
```bash
VITE_API_BASE_URL=http://localhost:8080/api
```

4. Запустите dev сервер:
```bash
npm run dev
# или
yarn dev
```

#### Локальная разработка Backend

1. Перейдите в папку backend:
```bash
cd backend
```

2. Установите зависимости:
```bash
go mod download
```

3. Запустите PostgreSQL локально или через Docker:
```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dealer_platform -p 5432:5432 -d postgres:14
```

4. Запустите миграции:
```bash
make migrate-up
```

5. Запустите сервер:
```bash
make run
```

### Структура проекта

```
dealer_dev_platform/
├── backend/                 # Go backend
│   ├── cmd/app/            # Точка входа
│   ├── internal/           # Внутренние пакеты
│   │   ├── app/           # Инициализация приложения
│   │   ├── config/       # Конфигурация
│   │   ├── database/     # Подключение к БД
│   │   ├── delivery/     # HTTP handlers
│   │   ├── model/        # Модели данных
│   │   ├── repository/    # Слой данных
│   │   ├── service/      # Бизнес-логика
│   │   └── utils/        # Утилиты
│   ├── migrations/        # SQL миграции
│   └── Dockerfile        # Docker образ для backend
├── frontend/              # React frontend
│   ├── src/              # Исходный код
│   │   ├── api/          # API клиенты
│   │   ├── components/   # React компоненты
│   │   ├── hooks/        # React hooks
│   │   ├── pages/        # Страницы
│   │   └── types/        # TypeScript типы
│   ├── Dockerfile        # Docker образ для frontend
│   └── nginx.conf        # Конфигурация Nginx
├── docker-compose.yml    # Docker Compose конфигурация
└── README.md            # Документация
```

### Команды Docker Compose

```bash
# Запуск всех сервисов
docker-compose up

# Запуск в фоновом режиме
docker-compose up -d

# Пересборка и запуск
docker-compose up --build

# Остановка сервисов
docker-compose down

# Просмотр логов
docker-compose logs -f

# Просмотр логов конкретного сервиса
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
```

### Тестирование

#### Backend тесты
```bash
cd backend
make test
```

#### Frontend тесты
```bash
cd frontend
npm test
# или
yarn test
```

### Мониторинг

- Health check endpoints:
  - Backend: http://localhost:8080/health
  - Frontend: http://localhost:3000/health

### Troubleshooting

#### Проблемы с подключением к базе данных
1. Убедитесь, что PostgreSQL запущен
2. Проверьте переменные окружения DATABASE_URL
3. Проверьте логи: `docker-compose logs postgres`

#### Проблемы с CORS
1. Проверьте настройки CORS в backend
2. Убедитесь, что frontend обращается к правильному URL API

#### Проблемы с сетью Docker
1. Проверьте, что все сервисы в одной сети `dealer_network`
2. Убедитесь, что сервисы могут обращаться друг к другу по именам

### Лицензия

MIT License