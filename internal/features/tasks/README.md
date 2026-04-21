# Feature: Tasks — Гексагональная архитектура

Фича `tasks` реализует CRUD задач и служит образцом **гексагональной архитектуры** (Ports & Adapters) в проекте.

---

## Содержание

1. [Что такое гексагональная архитектура](#1-что-такое-гексагональная-архитектура)
2. [Правило зависимостей](#2-правило-зависимостей)
3. [Структура фичи](#3-структура-фичи-tasks)
4. [Стратегия кеширования](#4-стратегия-кеширования)
5. [Как расширить фичу](#5-как-расширить-фичу)

---

## 1. Что такое гексагональная архитектура

Гексагональная архитектура (Hexagonal Architecture, Ports & Adapters) — подход, предложенный Алистером Кокберном (2005). Основная идея: **изолировать бизнес-логику** от любой инфраструктуры (HTTP, БД, кеш, очереди сообщений).

Приложение делится на три зоны:

```
┌─────────────────────────────────────────────────────┐
│                                                     │
│   [Входящий адаптер]     [Исходящий адаптер]        │
│   HTTP, gRPC, CLI...     Postgres, Redis, S3...     │
│          │                        ▲                 │
│          ▼                        │                 │
│   ┌─────────────────────────────────────────┐       │
│   │                Ядро (Hexagon)           │       │
│   │                                         │       │
│   │  [Входящий порт]    [Исходящий порт]    │       │
│   │  "Что сервис        "Что сервис требует │       │
│   │  предоставляет"     от инфраструктуры"  │       │
│   │                                         │       │
│   │            Уровень сервиса              │       │
│   │                                         │       │
│   └─────────────────────────────────────────┘       │
│                                                     │
└─────────────────────────────────────────────────────┘
```

**Порт** — это интерфейс (контракт), принадлежащий уровню сервиса. Уровень сервиса определяет, как с ним говорить.

**Адаптер** — это конкретная реализация, которая переводит «язык» внешней системы (например, базы данных PostgreSQL) в «язык» порта (контракта, определённого уровнем сервиса) и обратно.

**Почему «гексагон»?** Название символическое: шестиугольник наглядно показывает, что у уровня сервиса может быть много сторон — каждая сторона это порт, к которому можно подключить разные адаптеры.

---

## 2. Правило зависимостей

Главное правило: **зависимости направлены только внутрь, к ядру (сервису)**.

```
Внешний мир (Например, HTTP) → Адаптер → Порт → Сервис
                                                  │
                                                Порт → Адаптер → Внешний мир (Например, PostgreSQL)
```

Уровень сервиса (`TasksService`) не импортирует ни одного пакета HTTP или БД. Он оперирует только доменными типами и интерфейсами портов, которые сам же определяет.

---

## 3. Структура фичи "tasks"

```
internal/features/tasks/
│
│   # ─── УРОВЕНЬ СЕРВИСА ─────────────────────────────────────
│
├── service.go              # Структура TasksService + конструктор
├── create_task.go          # Use case: создание задачи
├── get_task.go             # Use case: получение задачи по ID
├── get_tasks.go            # Use case: список задач с пагинацией
├── patch_task.go           # Use case: частичное обновление (read-modify-write)
├── delete_task.go          # Use case: удаление задачи
│
│   # ─── ПОРТЫ ──────────────────────────────────────────────
│
├── ports/
│   ├── in/                 # Входящий порт (Что сервис из себя представляет для внешнего мира)
│   │   ├── service.go      #   Интерфейс TasksService
│   │   ├── create_task.go  #   CreateTaskParams / CreateTaskResult
│   │   ├── get_task.go     #   GetTaskParams / GetTaskResult
│   │   ├── get_tasks.go    #   GetTasksParams / GetTasksResult
│   │   ├── patch_task.go   #   PatchTaskParams / PatchTaskResult
│   │   └── delete_task.go  #   DeleteTaskParams / DeleteTaskResult
│   │
│   └── out/repository/     # Исходящий порт (Что сервис ожидает от внешнего мира)
│       ├── repository.go   #   Интерфейс TasksRepository
│       ├── save_task.go    #   SaveTaskParams / SaveTaskResult
│       ├── get_task.go     #   GetTaskParams / GetTaskResult
│       ├── get_tasks.go    #   GetTasksParams / GetTasksResult
│       ├── update_task.go  #   UpdateTaskParams / UpdateTaskResult
│       └── delete_task.go  #   DeleteTaskParams / DeleteTaskResult
│
│   # ─── АДАПТЕРЫ ────────────────────────────────────────────
│
└── adapters/
    ├── in/
    │   └── transport/http/ # Входящий адаптер (Адаптирует элемент внешнего мира HTTP под входящий порт сервиса)
    │       ├── transport.go    #   TasksHTTPHandler, маршруты
    │       ├── dto_common.go   #   Общий TaskDTOResponse
    │       ├── create_task.go  #   POST /tasks
    │       ├── get_task.go     #   GET /tasks/{id}
    │       ├── get_tasks.go    #   GET /tasks
    │       ├── patch_task.go   #   PATCH /tasks/{id}
    │       └── delete_task.go  #   DELETE /tasks/{id}
    │
    └── out/repository/     # Исходящие адаптеры (Адапитирует элементы внешнего мира PostgreSQL и Redis под исходящие порты сервиса)
        ├── postgres/        #   Реализация через PostgreSQL (pgx)
        │   ├── repository.go
        │   ├── models.go
        │   ├── save_task.go
        │   ├── get_task.go
        │   ├── get_tasks.go
        │   ├── update_task.go
        │   └── delete_task.go
        │
        └── cached/          #   Redis-кеш поверх основного репозитория
            ├── repository.go
            ├── models.go
            ├── tools.go
            ├── get_task.go
            ├── get_tasks.go
            ├── save_task.go
            ├── update_task.go
            └── delete_task.go
```

---

## 4. Стратегия кеширования

`CachedRepository` реализует паттерн **cache-aside** (ленивое заполнение).

### Одиночная задача (GET /tasks/{id})

```
GetTask(id)
    │
    ├─► Redis GET "task:<uuid>"
    │       │
    │   hit ├─► десериализовать → domain.Task → вернуть
    │       │
    │  miss └─► PostgreSQL SELECT WHERE id=...
    │               │
    │               └─► cacheTask: Redis SET "task:<uuid>" (TTL)
    │                       └─► вернуть
```

### Список задач (GET /tasks)

```
GetTasks(userID, limit, offset)
    │
    ├─► Redis HGET "tasks:<userID>" field="<limit>:<offset>"
    │       │
    │   hit ├─► десериализовать → []domain.Task → вернуть
    │       │
    │  miss └─► PostgreSQL SELECT ... LIMIT ... OFFSET ...
    │               │
    │               └─► Redis HSET "tasks:<userID>" field="<limit>:<offset>" value
    │                       └─► вернуть
```

Ключи Redis hash:
- `tasks:all` — все задачи (userID == nil)
- `tasks:<uuid>` — задачи конкретного пользователя

Поле (field) хеша кодирует пагинацию: `"10:0"`, `"nil:nil"`, `"20:40"`.
Такая структура позволяет инвалидировать все варианты одним `DEL` по ключу хеша.

### Инвалидация при мутациях

| Операция       | Что инвалидируется                                               |
|----------------|------------------------------------------------------------------|
| `SaveTask`     | `tasks:all`, `tasks:<authorUserID>` (списки устарели)            |
| `UpdateTask`   | `task:<id>` (обновлён), `tasks:all`, `tasks:<authorUserID>`      |
| `DeleteTask`   | `task:<id>` (удалён), `tasks:all`, `tasks:<authorUserID>`        |

При `DeleteTask` нужен `AuthorUserID` задачи для инвалидации — берётся из кеша или запросом к Postgres.

### Отказоустойчивость

Все ошибки кеша логируются, но **не прерывают запрос**. При недоступности Redis приложение продолжает работать, обращаясь напрямую к Postgres (graceful degradation).

---

## 5. Как расширить фичу

### Добавить новый входящий адаптер (например, gRPC)

1. Создать `adapters/in/transport/grpc/`
2. Реализовать gRPC-хендлер, зависящий от `ports_in.TasksService`
3. Зарегистрировать в `main.go`

Уровень сервиса (`service.go`, use cases, порты) — **не трогать**.

### Добавить новый исходящий адаптер (например, MongoDB)

1. Создать `adapters/out/repository/mongodb/`
2. Реализовать `ports_out_repository.TasksRepository`
3. В `main.go` заменить `tasks_postgres.NewTasksRepository` на новый адаптер

Уровень сервиса и HTTP-адаптер — **не трогать**.

### Добавить новый use case (например, `ArchiveTask`)

1. Добавить метод в `ports/in/service.go` (интерфейс)
2. Добавить Params/Result в `ports/in/archive_task.go`
3. Реализовать метод в `archive_task.go` (use case на уровне сервиса)
4. Если нужна новая операция с хранилищем — добавить метод в `ports/out/repository/repository.go`
   и реализовать новый метод порта в обоих адаптерах (`postgres/`, `cached/`)
5. Добавить HTTP/gRPC/Kafka хендлер в `adapters/in/transport/{http/grpc/kafka}/archive_task.go`
