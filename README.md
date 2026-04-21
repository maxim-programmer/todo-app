# Golang Todo App

REST API приложение на Go. Реализует управление пользователями, задачами и статистикой.

```bash
# Управление проектом
make help
```

## Технологический стек

| Компонент         | Технология                                       |
|-------------------|--------------------------------------------------|
| Язык              | Go 1.22+                                         |
| HTTP-фреймворк    | Стандартный `net/http` (без внешних фреймворков) |
| База данных       | PostgreSQL (`jackc/pgx/v5`)                      |
| Кеш               | Redis (`redis/go-redis/v9`)                      |
| Web-сервер        | Caddy 2                                          |
| Логгер            | `go.uber.org/zap`                                |
| Конфигурация      | `kelseyhightower/envconfig`                      |
| Валидация         | `go-playground/validator/v10`                    |
| Документация API  | Swagger (`swaggo/swag`)                          |
| Миграции БД       | golang-migrate                                   |
| Деплой            | Docker                                           |

---

## Архитектура

Разные фичи используют разные архитектурные подходы — намеренно, для демонстрации и сравнения.

### `users`, `statistics` — Чистая архитектура (Clean Architecture)

Классическое трёхслойное разделение с инверсией зависимостей:

```
Transport (HTTP Handler)
      │   Декодирует запрос, вызывает сервис, формирует ответ
      ↓
Service (Business Logic)          ← интерфейс TasksService живёт здесь
      │   Валидация, доменная логика
      ↓
Repository (Data Access)          ← интерфейс TasksRepository живёт здесь
      └─  SQL-запросы к PostgreSQL, маппинг моделей

Domain (Core)
          Сущности, инварианты, бизнес-правила — без зависимостей
```

Интерфейсы определяются в **потребляющем** слое, а не в реализующем — это и есть DIP.

---

### `tasks` — Гексагональная архитектура (Ports & Adapters)

> Подробная документация: [`internal/features/tasks/README.md`](internal/features/tasks/README.md)

Явное разделение на порты (контракты, принадлежащие ядру) и адаптеры (реализации инфраструктуры):

```
         ┌──────────────────────────────────────────────┐
         │                   ЯДРО                       │
  HTTP   │  ports/in         Service        ports/out   │  Postgres
 ──────► │  TasksService ──► tasks_service ──► TasksRepo│ ──────────►
         │  (интерфейс)      (логика)        (интерфейс)│  Redis
         └──────────────────────────────────────────────┘
    ▲                                                         ▲
адаптер/in                                             адаптер/out
(транслирует HTTP → порт)                  (реализует порт → SQL/кеш)
```

**Ключевые принципы:**

- **Ядро не зависит от адаптеров** — только от интерфейсов портов, которые само же определяет
- **Входящий порт** (`ports/in`) — контракт, который ядро предоставляет внешнему миру (HTTP, gRPC...)
- **Исходящий порт** (`ports/out`) — контракт, который ядро требует от инфраструктуры (БД, кеш...)
- **`CachedRepository`** — исходящий адаптер-декоратор: добавляет Redis-кеш поверх Postgres,
  не изменяя интерфейс порта и не затрагивая ядро

**Dependency Injection** реализован вручную в `cmd/todoapp/main.go`:

```
tasks_postgres.Repository
        ↓ обёрнут в
tasks_cached.CachedRepository  →  TasksService  →  TasksHTTPHandler
```

---

## Структура проекта

```
.
├── cmd/
│   └── todoapp/
│       └── main.go                    # Точка входа: DI, инициализация, запуск
├── internal/
│   ├── core/                          # Общие компоненты, не зависящие от фич
│   │   ├── config/                    # Общая конфигурация приложения
│   │   ├── domain/                    # Доменные сущности: Task, User, Statistics...
│   │   ├── errors/                    # Sentinel-ошибки: ErrNotFound, ErrConflict...
│   │   ├── logger/                    # Структурированный логгер (zap) + logger-in-context
│   │   └── repository/
│   │       ├── postgres/pool/         # Интерфейс пула + реализация на pgx
│   │       └── redis/pool/            # Интерфейс пула + реализация на go-redis
│   │   └── transport/http/
│   │       ├── middleware/            # CORS, RequestID, Logger, Trace, Panic
│   │       ├── request/               # Хелперы: декодирование тела, path/query параметры
│   │       ├── response/              # HTTPResponseHandler, ErrorResponse
│   │       ├── server/                # HTTPServer, APIVersionRouter, Route
│   │       └── types/                 # Nullable[T] с UnmarshalJSON для PATCH-запросов
│   └── features/
│       ├── tasks/                     # CRUD задач — гексагональная архитектура (см. tasks/README.md)
│       │   ├── ports/
│       │   │   ├── in/                # Входящий порт: TasksService интерфейс + DTO
│       │   │   └── out/repository/    # Исходящий порт: TasksRepository интерфейс + DTO
│       │   ├── adapters/
│       │   │   ├── in/transport/http/ # Driving adapter: HTTP → порт
│       │   │   └── out/repository/
│       │   │       ├── postgres/      # Driven adapter: порт → PostgreSQL
│       │   │       └── cached/        # Driven adapter: декоратор с Redis-кешем
│       │   └── *.go                   # Ядро: сервис + use cases
│       ├── users/                     # CRUD пользователей — чистая архитектура
│       ├── statistics/                # Статистика по задачам — чистая архитектура
│       └── web/                       # Отдача статических страниц
├── web/
│   ├── public/                        # Статические файлы (HTML, CSS, JS)
│   └── Caddyfile                      # Конфигурация Caddy (продакшн)
├── migrations/                        # SQL-миграции (golang-migrate)
├── docs/                              # Автогенерированная Swagger-документация
├── docker-compose.yaml                # Инфраструктура: PostgreSQL, Redis, Caddy, Swagger
├── postman_collection.json            # Коллекция запросов для Postman
└── Makefile                           # Команды для разработки и деплоя
```

---

## Локальный запуск

#### Предварительные требования
- `Docker и Docker Compose`
- `Go 1.22+`
- `make`

#### Шаги

```bash
# 1. Создать .env по примеру
cp .env.example .env

# 2. Выставить недостающие переменные окружения
code .env

# 3. Поднять окружение (PostgreSQL + Redis)
make env-up

# 4. Применить миграции БД
make migrate-up

# 5. Открыть порты сервисов окружения (PostgreSQL :5432, Redis :6379)
make env-port-forward

# 6. Запустить приложение локально
make todoapp-run
```

## Нагрузочное тестирование
```bash
# 1. Предварительно очистим окружение
make env-cleanup

# 2. Запустим окружение
make env-up

# 3. Накатим миграции
make migrate-up

# 4. Откроем порты окружения
make env-port-forward

# 5. Запуск скрипта нагрузочного тестирования (далее по инструкции в самом скрипте)
make load-test
```

После запуска:
- Главная страница: Недоступна. Раздаётся только с Web-Сервера при деплое.
- Swagger UI: `http://127.0.0.1:5050/swagger/`
- API: `http://127.0.0.1:5050/api/v1/`

## Деплой

```bash
# Запустить PostgreSQL + Redis
make env-up

# Применить миграции
make migrate-up

# Запустить Go-приложение в Docker
make todoapp-deploy

# Запустить Caddy web-сервер (опционально, для продакшн)
make web-deploy
```

После запуска:
- Главная страница: http(s)/{CADDY_HOST}.
- Swagger UI: Недоступен. Раздаётся только с HTTP-Сервера Golang приложения
- API: `http(s)://{CADDY_HOST}/api/v1/`

---

## Переменные окружения

| Переменная               | Описание                                      | Пример / По умолчанию                   |
|--------------------------|-----------------------------------------------|-----------------------------------------|
| `TIME_ZONE`              | Часовой пояс (IANA)                           | `Europe/Moscow`                         |
| `LOGGER_LEVEL`           | Уровень логирования                           | `DEBUG`                                 |
| `LOGGER_FOLDER`          | Директория для лог-файлов                     | `out/logs`                              |
| `POSTGRES_HOST`          | Хост PostgreSQL                               | `localhost`                             |
| `POSTGRES_PORT`          | Порт PostgreSQL                               | `5432`                                  |
| `POSTGRES_USER`          | Пользователь БД                               | `todoapp-test-user`                     |
| `POSTGRES_PASSWORD`      | Пароль БД                                     | `todoapp-test-password`                 |
| `POSTGRES_DB`            | Имя базы данных                               | `todoapp`                               |
| `POSTGRES_TIMEOUT`       | Таймаут запроса к БД                          | `5s`                                    |
| `REDIS_HOST`             | Хост Redis                                    | `localhost`                             |
| `REDIS_PORT`             | Порт Redis                                    | `6379`                                  |
| `REDIS_PASSWORD`         | Пароль Redis                                  | —                                       |
| `REDIS_DB`               | Номер базы данных Redis (0–15)                | `0`                                     |
| `REDIS_TTL`              | TTL кешированных записей                      | `5m`                                    |
| `HTTP_ADDR`              | Адрес и порт HTTP-сервера                     | `:5050`                                 |
| `HTTP_SHUTDOWN_TIMEOUT`  | Таймаут graceful shutdown                     | `30s`                                   |
| `HTTP_ALLOWED_ORIGINS`   | Разрешённые CORS origins (через запятую)      | `http://localhost:3000,null`            |
| `CADDY_HOST`             | Хост Caddy (домен или `localhost`)            | `localhost`                             |
| `PROJECT_ROOT`           | Корень проекта (для Docker volumes)           | `/Users/MyUser/projects/golang-todoapp` |

---

## API

### Пользователи `/api/v1/users`

| Метод    | Путь          | Описание                         |
|----------|---------------|----------------------------------|
| `POST`   | `/users`      | Создать пользователя             |
| `GET`    | `/users`      | Список пользователей (пагинация) |
| `GET`    | `/users/{id}` | Получить пользователя по ID      |
| `PATCH`  | `/users/{id}` | Частично обновить пользователя   |
| `DELETE` | `/users/{id}` | Удалить пользователя             |

### Задачи `/api/v1/tasks`

| Метод    | Путь          | Описание                                        |
|----------|---------------|-------------------------------------------------|
| `POST`   | `/tasks`      | Создать задачу                                  |
| `GET`    | `/tasks`      | Список задач (пагинация + фильтр по `user_id`)  |
| `GET`    | `/tasks/{id}` | Получить задачу по ID                           |
| `PATCH`  | `/tasks/{id}` | Частично обновить задачу                        |
| `DELETE` | `/tasks/{id}` | Удалить задачу                                  |

### Статистика `/api/v1/statistics`

| Метод  | Путь           | Описание                                                                 |
|--------|----------------|--------------------------------------------------------------------------|
| `GET`  | `/statistics`  | Статистика задач (фильтры: `user_id`, `from`, `to` в формате YYYY-MM-DD) |

Полная интерактивная документация доступна в **Swagger UI** по адресу `/swagger/`.
