English version: jump to [#english](#english)
Русская версия: перейти к [#russian](#russian)

<div id="english"></div>

# TaskManager

REST API + CLI for task management with JWT auth, role-based access, Kafka event streaming, Redis caching, and rate limiting.

---

## Tech Stack

| Category | Technology | Purpose |
|----------|-----------|---------|
| **Language** | Go 1.25+ | Core language |
| **HTTP** | Gin | HTTP router & middleware |
| **Database** | PostgreSQL + sqlx | Persistent storage |
| **Migrations** | Goose | SQL schema versioning |
| **Auth** | JWT (RSA 2048-bit) | Token-based authentication |
| **Cache** | Redis | Caching & rate limiting |
| **Messaging** | Kafka (Confluent 7.5) | Async event streaming |
| **Docs** | Swagger (swaggo) | API documentation |
| **CLI** | Cobra | Command-line interface |
| **Config** | godotenv | Environment management |
| **Container** | Docker + Compose | Multi-service orchestration |
| **Testing** | testify | Unit & integration tests |
| **Linting** | golangci-lint | Static analysis |

---

## Project Structure (Hexagonal Architecture)

```
├── cmd/
│   ├── web/              # HTTP server entrypoint + Swagger docs
│   └── cli/              # CLI entrypoint (Cobra)
├── internal/
│   ├── core/
│   │   ├── domain/       # Entities (Task, User), constants, errors
│   │   ├── port/         # Interfaces (inbound & outbound ports)
│   │   └── service/      # Business logic (user, task, token, validator, cache, broker)
│   ├── adapter/
│   │   ├── handler/      # HTTP handlers + middleware (inbound adapter)
│   │   ├── repo/
│   │   │   ├── postgres/ # PostgreSQL repositories (outbound adapter)
│   │   │   └── redis/    # Redis client — caching & rate limiting (outbound adapter)
│   │   └── kafka/
│   │       ├── producer/ # Kafka producer (outbound adapter)
│   │       ├── consumer/ # Kafka consumer with worker pool (inbound adapter)
│   │       └── shared/   # Kafka config & shared message types
│   ├── app/              # Wiring, config, server setup
│   └── testutil/         # Mocks & test helpers
├── migrations/           # Goose SQL migrations
├── certs/                # RSA key pair for JWT
├── compose.yaml          # Docker Compose (app, postgres, redis, kafka, zookeeper)
├── Dockerfile            # Multi-stage Go build
├── Makefile              # Dev commands
└── .env.example          # Environment template
```

---

## Database Migrations

Powered by **Goose** — a database migration tool supporting SQL and Go migrations.

### Available Commands

```bash
make migrate-up              # Apply all pending migrations
make migrate-down            # Rollback the last migration
make migrate-status          # Show migration status
make migrate-reset           # Drop all tables and re-run migrations
make migrate-create name=X   # Create a new migration file
```

### Migration Files

| File | Description |
|------|-------------|
| `20260712144233_create_users.sql` | Creates `users` table with index on `username` |
| `20260712144251_create_tasks.sql` | Creates `tasks` table with indexes on `user_id` and `title` |
| `20260712150111_insert_default_users_data.sql` | Seeds 3 users: admin, test, john_doe |
| `20260712150115_insert_default_tasks_data.sql` | Seeds 7 tasks across all users |

### Schema Overview

```
users                           tasks
┌──────────────┐               ┌──────────────┐
│ id       PK │               │ id       PK │
│ username    │               │ title        │
│ role        │               │ text         │
│ email    UQ │               │ priority     │
│ password_hash│              │ created      │
│ created_at  │               │ expires      │
└──────┬───────┘               │ user_id  FK │──→ users.id (CASCADE)
       │                       └──────────────┘
       └──────────────────────────────────────────
```

---

## Kafka Integration

Event-driven architecture using **Apache Kafka** (Confluent Platform 7.5) for asynchronous message processing.

### Topics

| Topic | Events | Description |
|-------|--------|-------------|
| `users` | create, patch, delete | User lifecycle events |
| `tasks` | create, patch, delete | Task lifecycle events |

### Architecture

```
┌─────────────┐    publish     ┌─────────┐    consume     ┌──────────────┐
│ HTTP Handler │───────────────→│  Kafka  │───────────────→│   Consumer   │
│  (Producer)  │               │  Broker │               │ (4 workers)  │
└─────────────┘               └─────────┘               └──────┬───────┘
                                                               │
                                                               ▼
                                                         ┌──────────┐
                                                         │ Postgres │
                                                         └──────────┘
```

- **Producer**: Sends JSON messages with operation type (`create`/`patch`/`delete`)
- **Consumer**: Worker pool (4 goroutines) processes messages from topics
- **Auto-retry**: 3 retries with 250ms backoff on leader unavailable
- **Auto topic creation**: Topics are created automatically if missing

### Message Format

```json
{
  "operation": "create",
  "title": "New Task",
  "text": "Description",
  "priority": 2,
  "expire_days": 7,
  "user_id": "1"
}
```

---

## Redis

Two primary use cases:

1. **Caching** — Task data cached with TTL to reduce DB load
2. **Rate Limiting** — Per-user request counter with 1-minute sliding window

---

## Docker Compose

Single command to start all services:

```bash
make compose-up       # Start app, postgres, redis, kafka, zookeeper
make compose-down     # Stop all services
```

### Services

| Service | Image | Port | Purpose |
|---------|-------|------|---------|
| `app` | Custom (Dockerfile) | 8080 | REST API server |
| `db` | postgres:latest | 5432 | PostgreSQL database |
| `redis` | redis:latest | 6379 | Cache & rate limiter |
| `zookeeper` | confluentinc/cp-zookeeper:7.5.0 | 2181 | Kafka coordination |
| `kafka` | confluentinc/cp-kafka:7.5.0 | 9092, 29092 | Message broker |

All services include health checks and dependency ordering.

---

## Setup

```bash
# 1. Clone and enter the project
git clone https://github.com/Emin-07/TaskManager.git
cd TaskManager

# 2. Generate RSA keys for JWT
make certs

# 3. Configure environment
make env              # copies .env.example → .env, edit with your credentials

# 4. Start infrastructure (or run services locally)
make compose-up

# 5. Apply database migrations
make migrate-up

# 6. Run the app
go run ./cmd/web/     # HTTP server on :8080
go run ./cmd/cli/     # CLI interface (Cobra)
```

Swagger UI: `http://localhost:8080/swagger`

---

## API Endpoints

### Auth

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/signup` | Register a new user |
| POST | `/auth/login` | Get JWT token |

### Tasks (requires JWT)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks` | List all user tasks |
| GET | `/tasks/:id` | Get task by ID |
| POST | `/tasks` | Create a new task |
| PATCH | `/tasks/:id` | Update a task |
| DELETE | `/tasks/:id` | Delete a task |

### Users (admin only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users` | List all users |
| GET | `/users/:id` | Get user by ID |
| PATCH | `/users/:id` | Update a user |
| DELETE | `/users/:id` | Delete a user |

---

## Testing

```bash
make test-unit           # Run unit tests (core/)
make test-integration    # Run integration tests (adapter/)
make test-all            # Run all tests
```

---

<div id="russian"></div>

# TaskManager

REST API + CLI для управления задачами с JWT-аутентификацией, ролевой моделью, Kafka для событийного стриминга, Redis для кэширования и rate limiting.

---

## Стек технологий

| Категория | Технология | Назначение |
|-----------|-----------|------------|
| **Язык** | Go 1.25+ | Основной язык |
| **HTTP** | Gin | HTTP-роутер и middleware |
| **БД** | PostgreSQL + sqlx | Персистентное хранилище |
| **Миграции** | Goose | Версионирование SQL-схемы |
| **Авторизация** | JWT (RSA 2048-bit) | Токенная аутентификация |
| **Кэш** | Redis | Кэширование и rate limiting |
| **Мессенджер** | Kafka (Confluent 7.5) | Асинхронная обработка событий |
| **Документация** | Swagger (swaggo) | Документация API |
| **CLI** | Cobra | Командный интерфейс |
| **Конфигурация** | godotenv | Управление переменными окружения |
| **Контейнеры** | Docker + Compose | Оркестрация мультисервисного стека |
| **Тестирование** | testify | Unit и интеграционные тесты |
| **Линтер** | golangci-lint | Статический анализ |

---

## Структура проекта (Гексагональная архитектура)

```
├── cmd/
│   ├── web/              # Точка входа HTTP-сервера + Swagger
│   └── cli/              # Точка входа CLI (Cobra)
├── internal/
│   ├── core/
│   │   ├── domain/       # Сущности (Task, User), константы, ошибки
│   │   ├── port/         # Интерфейсы (входящие и исходящие порты)
│   │   └── service/      # Бизнес-логика (user, task, token, validator, cache, broker)
│   ├── adapter/
│   │   ├── handler/      # HTTP-обработчики + middleware (входящий адаптер)
│   │   ├── repo/
│   │   │   ├── postgres/ # PostgreSQL репозитории (исходящий адаптер)
│   │   │   └── redis/    # Redis клиент — кэш и rate limiting (исходящий адаптер)
│   │   └── kafka/
│   │       ├── producer/ # Kafka producer (исходящий адаптер)
│   │       ├── consumer/ # Kafka consumer с пулом воркеров (входящий адаптер)
│   │       └── shared/   # Конфигурация Kafka и общие типы сообщений
│   ├── app/              # Сборка, конфигурация, запуск сервера
│   └── testutil/         # Моки и вспомогательные функции для тестов
├── migrations/           # SQL-миграции Goose
├── certs/                # RSA-ключи для JWT
├── compose.yaml          # Docker Compose (app, postgres, redis, kafka, zookeeper)
├── Dockerfile            # Multi-stage сборка Go
├── Makefile              # Команды для разработки
└── .env.example          # Шаблон переменных окружения
```

---

## Миграции базы данных

Используется **Goose — инструмент для миграций, поддерживающий SQL и Go-миграции.

### Доступные команды

```bash
make migrate-up              # Применить все pending миграции
make migrate-down            # Откатить последнюю миграцию
make migrate-status          # Показать статус миграций
make migrate-reset           # Удалить все таблицы и перезапустить миграции
make migrate-create name=X   # Создать новый файл миграции
```

### Файлы миграций

| Файл | Описание |
|------|----------|
| `20260712144233_create_users.sql` | Создаёт таблицу `users` с индексом по `username` |
| `20260712144251_create_tasks.sql` | Создаёт таблицу `tasks` с индексами по `user_id` и `title` |
| `20260712150111_insert_default_users_data.sql` | Заполняет 3 пользователей: admin, test, john_doe |
| `20260712150115_insert_default_tasks_data.sql` | Заполняет 7 задач для всех пользователей |

### Схема базы данных

```
users                           tasks
┌──────────────┐               ┌──────────────┐
│ id       PK │               │ id       PK │
│ username    │               │ title        │
│ role        │               │ text         │
│ email    UQ │               │ priority     │
│ password_hash│              │ created      │
│ created_at  │               │ expires      │
└──────┬───────┘               │ user_id  FK │──→ users.id (CASCADE)
       │                       └──────────────┘
       └──────────────────────────────────────────
```

---

## Интеграция с Kafka

Event-driven архитектура на базе **Apache Kafka** (Confluent Platform 7.5) для асинхронной обработки сообщений.

### Топики

| Топик | События | Описание |
|-------|---------|----------|
| `users` | create, patch, delete | События жизненного цикла пользователей |
| `tasks` | create, patch, delete | События жизненного цикла задач |

### Архитектура

```
┌─────────────┐    publish     ┌─────────┐    consume     ┌──────────────┐
│ HTTP Handler │───────────────→│  Kafka  │───────────────→│   Consumer   │
│  (Producer)  │               │  Broker │               │ (4 workers)  │
└─────────────┘               └─────────┘               └──────┬───────┘
                                                               │
                                                               ▼
                                                         ┌──────────┐
                                                         │ Postgres │
                                                         └──────────┘
```

- **Producer**: Отправляет JSON-сообщения с типом операции (`create`/`patch`/`delete`)
- **Consumer**: Пул из 4 горутин обрабатывает сообщения из топиков
- **Автоперезапуск**: 3 попытки с 250ms паузой при недоступности лидера
- **Автосоздание топиков**: Топики создаются автоматически при отсутствии

### Формат сообщения

```json
{
  "operation": "create",
  "title": "Новая задача",
  "text": "Описание",
  "priority": 2,
  "expire_days": 7,
  "user_id": "1"
}
```

---

## Redis

Два основных сценария использования:

1. **Кэширование** — Данные задач кэшируются с TTL для снижения нагрузки на БД
2. **Rate Limiting** — Счётчик запросов на пользователя с 1-минутным скользящим окном

---

## Docker Compose

Одна команда для запуска всех сервисов:

```bash
make compose-up       # Запустить app, postgres, redis, kafka, zookeeper
make compose-down     # Остановить все сервисы
```

### Сервисы

| Сервис | Образ | Порт | Назначение |
|--------|-------|------|------------|
| `app` | Свой (Dockerfile) | 8080 | REST API сервер |
| `db` | postgres:latest | 5432 | PostgreSQL база данных |
| `redis` | redis:latest | 6379 | Кэш и rate limiter |
| `zookeeper` | confluentinc/cp-zookeeper:7.5.0 | 2181 | Координация Kafka |
| `kafka` | confluentinc/cp-kafka:7.5.0 | 9092, 29092 | Брокер сообщений |

Все сервисы имеют health checks и правильный порядок зависимостей.

---

## Запуск

```bash
# 1. Клонировать и перейти в проект
git clone https://github.com/Emin-07/TaskManager.git
cd TaskManager

# 2. Сгенерировать RSA-ключи для JWT
make certs

# 3. Настроить окружение
make env              # копирует .env.example → .env, отредактировать под свои данные

# 4. Запустить инфраструктуру (или запускать сервисы локально)
make compose-up

# 5. Применить миграции базы данных
make migrate-up

# 6. Запустить приложение
go run ./cmd/web/     # HTTP-сервер на :8080
go run ./cmd/cli/     # CLI-интерфейс (Cobra)
```

Swagger UI: `http://localhost:8080/swagger`

---

## API Эндпоинты

### Авторизация

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | `/auth/signup` | Регистрация нового пользователя |
| POST | `/auth/login` | Получение JWT-токена |

### Задачи (требуется JWT)

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| GET | `/tasks` | Список задач пользователя |
| GET | `/tasks/:id` | Получить задачу по ID |
| POST | `/tasks` | Создать новую задачу |
| PATCH | `/tasks/:id` | Обновить задачу |
| DELETE | `/tasks/:id` | Удалить задачу |

### Пользователи (только admin)

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| GET | `/users` | Список всех пользователей |
| GET | `/users/:id` | Получить пользователя по ID |
| PATCH | `/users/:id` | Обновить пользователя |
| DELETE | `/users/:id` | Удалить пользователя |

---

## Тестирование

```bash
make test-unit           # Unit-тесты (core/)
make test-integration    # Интеграционные тесты (adapter/)
make test-all            # Все тесты
```
