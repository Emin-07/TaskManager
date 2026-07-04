<!--
  English version: jump to [#english](#english)
  Русская версия: перейти к [#russian](#russian)
-->

<div id="english"></div>

# 📋 TaskManager

REST API + CLI for task management with JWT auth and role-based access (admin/user).

---

## Tech Stack

- **Go** 1.25.6
- **Gin** — HTTP router
- **PostgreSQL** + **sqlx** — database
- **JWT** (RSA) — authentication
- **Swagger** — API docs (`/swagger`)
- **Cobra** — CLI interface
- **godotenv** — configuration

---

## Project Structure (Hexagonal)

```
├── cmd/
│   ├── web/          # HTTP server entrypoint
│   └── cli/          # CLI entrypoint (Cobra)
└── internal/
    ├── core/
    │   ├── domain/   # Entities (Task, User)
    │   ├── port/     # Interfaces (in/out ports)
    │   └── service/  # Business logic
    ├── adapter/
    │   ├── handler/  # HTTP handlers (inbound adapter)
    │   └── repo/     # PostgreSQL repo (outbound adapter)
    └── app/          # Wiring & server setup
```

---

## Setup

```bash
make env          # cp .env.example → .env, edit with your DB credentials
make certs        # generate RSA key pair for JWT
make swag         # regenerate Swagger docs

go mod tidy
go run ./cmd/web/           # start web server
go run ./cmd/cli/           # CLI version (Cobra-based)
```

Swagger UI will be available at `http://localhost:8080/swagger` after startup.

---

<div id="russian"></div>

# 📋 TaskManager

REST API + CLI для управления задачами с JWT-аутентификацией и ролевой моделью (admin/user).

---

## Стек

- **Go** 1.25.6
- **Gin** — HTTP роутер
- **PostgreSQL** + **sqlx** — база данных
- **JWT** (RSA) — аутентификация
- **Swagger** — документация (`/swagger`)
- **Cobra** — CLI интерфейс
- **godotenv** — конфигурация

---

## Структура проекта (Hexagonal)

```
├── cmd/
│   ├── web/          # точка входа HTTP-сервера
│   └── cli/          # точка входа CLI (Cobra)
└── internal/
    ├── core/
    │   ├── domain/   # сущности (Task, User)
    │   ├── port/     # интерфейсы (входящие/исходящие порты)
    │   └── service/  # бизнес-логика
    ├── adapter/
    │   ├── handler/  # HTTP-обработчики (входящий адаптер)
    │   └── repo/     # PostgreSQL репозиторий (исходящий адаптер)
    └── app/          # сборка и запуск сервера
```

---

## Запуск

```bash
make env          # cp .env.example → .env, отредактировать под свою БД
make certs        # сгенерировать RSA-ключи для JWT
make swag         # перегенерировать Swagger-документацию

go mod tidy
go run ./cmd/web/           # запуск веб-сервера
go run ./cmd/cli/           # CLI-версия (на Cobra)
```

После запуска Swagger UI будет доступен по адресу `http://localhost:8080/swagger`.
