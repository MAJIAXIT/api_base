<div align="center">
  <a href="#ru">🇷🇺 Русский</a> &nbsp;|&nbsp; <a href="#en">🇬🇧 English</a>
</div>

---

<a name="ru"></a>

# Go REST API — production-ready шаблон

Готовый к продакшену шаблон REST API на Go с чистой архитектурой, JWT-аутентификацией, управлением сессиями и транзакционным middleware.

Используется как основа для быстрого старта новых проектов — вместо того чтобы каждый раз реализовывать одно и то же с нуля.

## Архитектура

Проект строго разделён по слоям — каждый слой зависит только от нижележащего через интерфейс:

```
cmd/app/          → точка входа
internal/
  handlers/       → HTTP-слой: парсинг запросов, формирование ответов
  service/        → бизнес-логика, независима от транспорта
  models/         → доменные модели, GORM-маппинг
  dto/            → структуры запросов и ответов
  middleware/     → транзакции, аутентификация
pkg/
  database/       → инициализация PostgreSQL
  logger/         → структурированное логирование
  password/       → bcrypt-шифрование
  utils/          → обработка ошибок
config/           → конфигурация из env-переменных
```

Все сервисы и хендлеры реализованы через интерфейсы — легко мокировать в тестах и заменять реализацию без изменения зависимых слоёв.

## Ключевые паттерны

**Транзакционный middleware**

Каждый POST/PUT/DELETE запрос автоматически оборачивается в транзакцию PostgreSQL. При успешном ответе — коммит, при ошибке или статусе ≥ 400 — автоматический откат. Хендлеры получают транзакцию через контекст Gin.

```go
// middleware сам управляет транзакцией
tx := h.transactionsMiddleware.GetTx(c)
user, err := h.usersService.GetUserByID(tx, userID)
```

**JWT + сессии**

- Access token (15 мин) — стандартный JWT
- Refresh token (7 дней) — JWT + запись в БД (сессия)
- Logout/LogoutAll — физически удаляет сессии из БД, исключая replay-атаки
- WebSocket auth — токен передаётся через query-параметр

**Graceful shutdown**

Сервер перехватывает `SIGINT`/`SIGTERM`, отменяет фоновые задачи и дожидается завершения активных запросов с таймаутом 5 секунд.

## Стек

| Компонент | Технология |
|-----------|-----------|
| Язык | Go 1.24 |
| HTTP-фреймворк | Gin |
| ORM | GORM |
| База данных | PostgreSQL |
| Аутентификация | JWT (golang-jwt/jwt v5) |
| Контейнеризация | Docker, Docker Compose |
| Прокси | Nginx |

## API

Все эндпоинты доступны по базовому пути `/api/v1`.

**Аутентификация**

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/auth/signup` | Регистрация |
| `POST` | `/auth/login` | Вход, возвращает access + refresh токены |
| `POST` | `/auth/refresh` | Обновление токенов по refresh token |
| `POST` | `/auth/logout` | Выход (требует Bearer token) |
| `GET` | `/auth/me` | Текущий пользователь (требует Bearer token) |

**Пользователи**

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/users/current` | Профиль текущего пользователя |
| `POST` | `/users/current/update` | Обновление профиля |

## Запуск

**1. Скопируй `.env.example` в `.env` и заполни переменные**

```bash
cp .env.example .env
```

```env
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=mydb

SERVER_HTTP_PORT=8080
JWT_SECRET=your-secret-key
PASSWORD_KEY=your-password-key
```

**2. Запусти через Make**

```bash
make postgres   # поднять PostgreSQL
make migrate    # применить миграции (с автоматическим дампом БД перед запуском)
make api        # собрать и запустить API
make client     # поднять Nginx-прокси (опционально)
```

**Другие команды:**

```bash
make postgres_it    # подключиться к psql интерактивно
make postgres_dump  # дамп базы данных
make test           # запустить тесты в Docker
```

## Структура Docker Compose

Проект разделён на три независимых `compose.yml`:

- `postgres/compose.yml` — база данных в отдельной сети `db-network`
- `api/compose.yml` — Go-сервис, подключается к `db-network`
- `client/compose.yml` — Nginx-прокси для статики или reverse-proxy

## Расширение шаблона

Чтобы добавить новый домен (например, `products`):

1. Создай модель в `internal/models/products/`
2. Создай DTO в `internal/dto/products/`
3. Реализуй сервис в `internal/service/products/` за интерфейсом
4. Реализуй хендлер в `internal/handlers/products/` за интерфейсом
5. Подключи в `internal/app/app.go` по аналогии с `auth` и `users`

---

<div align="center">
  <a href="#ru">🇷🇺 Русский</a> &nbsp;|&nbsp; <a href="#en">🇬🇧 English</a>
</div>

---

<a name="en"></a>

# Go REST API — production-ready boilerplate

A production-ready REST API template in Go with clean architecture, JWT authentication, session management, and transactional middleware.

Used as a foundation for new projects — eliminating the need to implement the same core infrastructure from scratch each time.

## Architecture

The project is strictly layered — each layer depends only on the one below it, through an interface:

```
cmd/app/          → entry point
internal/
  handlers/       → HTTP layer: request parsing, response formatting
  service/        → business logic, transport-agnostic
  models/         → domain models, GORM mapping
  dto/            → request/response structs
  middleware/     → transactions, authentication
pkg/
  database/       → PostgreSQL initialization
  logger/         → structured logging
  password/       → bcrypt encryption
  utils/          → error handling
config/           → configuration from env variables
```

All services and handlers are implemented behind interfaces — easy to mock in tests and swap implementations without touching dependent layers.

## Key patterns

**Transaction middleware**

Every POST/PUT/DELETE request is automatically wrapped in a PostgreSQL transaction. On success — commit; on error or status ≥ 400 — automatic rollback. Handlers receive the transaction through the Gin context.

```go
// middleware manages the transaction automatically
tx := h.transactionsMiddleware.GetTx(c)
user, err := h.usersService.GetUserByID(tx, userID)
```

**JWT + sessions**

- Access token (15 min) — standard JWT
- Refresh token (7 days) — JWT + database session record
- Logout/LogoutAll — physically deletes sessions from DB, preventing replay attacks
- WebSocket auth — token passed via query parameter

**Graceful shutdown**

Server intercepts `SIGINT`/`SIGTERM`, cancels background jobs, and waits for active requests to complete with a 5-second timeout.

## Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.24 |
| HTTP framework | Gin |
| ORM | GORM |
| Database | PostgreSQL |
| Authentication | JWT (golang-jwt/jwt v5) |
| Containerization | Docker, Docker Compose |
| Proxy | Nginx |

## API

All endpoints are available under the base path `/api/v1`.

**Authentication**

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/auth/signup` | Registration |
| `POST` | `/auth/login` | Login, returns access + refresh tokens |
| `POST` | `/auth/refresh` | Refresh tokens using refresh token |
| `POST` | `/auth/logout` | Logout (requires Bearer token) |
| `GET` | `/auth/me` | Current user (requires Bearer token) |

**Users**

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/users/current` | Current user profile |
| `POST` | `/users/current/update` | Update profile |

## Setup

**1. Copy `.env.example` to `.env` and fill in the variables**

```bash
cp .env.example .env
```

```env
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=mydb

SERVER_HTTP_PORT=8080
JWT_SECRET=your-secret-key
PASSWORD_KEY=your-password-key
```

**2. Run via Make**

```bash
make postgres   # start PostgreSQL
make migrate    # apply migrations (auto-dumps DB before running)
make api        # build and start the API
make client     # start Nginx proxy (optional)
```

**Other commands:**

```bash
make postgres_it    # connect to psql interactively
make postgres_dump  # dump the database
make test           # run tests in Docker
```

## Docker Compose structure

The project uses three independent `compose.yml` files:

- `postgres/compose.yml` — database on a dedicated `db-network`
- `api/compose.yml` — Go service, connects to `db-network`
- `client/compose.yml` — Nginx proxy for static files or reverse proxy

## Extending the template

To add a new domain (e.g. `products`):

1. Create a model in `internal/models/products/`
2. Create DTOs in `internal/dto/products/`
3. Implement a service in `internal/service/products/` behind an interface
4. Implement a handler in `internal/handlers/products/` behind an interface
5. Wire it up in `internal/app/app.go` following the pattern of `auth` and `users`
