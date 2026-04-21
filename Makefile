include .env
export


export PROJECT_ROOT=$(shell pwd)


.DEFAULT_GOAL := help


env-up: ## env: Запустить окружение проекта
	@docker compose up -d todoapp-postgres todoapp-redis

env-down: ## env: Остановить окружение проекта
	@docker compose down todoapp-postgres todoapp-redis

env-cleanup: ## env: Очистить окружение проекта
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres todoapp-redis port-forwarder web-server && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		rm -rf ${PROJECT_ROOT}/out/redis_data && \
		rm -rf ${PROJECT_ROOT}/out/caddy_data && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

env-port-forward: ## env: Открыть порты сервисов окружения
	@docker compose up -d port-forwarder

env-port-close: ## env: Закрыть порты сервисов окружения
	@docker compose down port-forwarder

logs-cleanup: ## env: Очистить файлы логов из out/logs
	@read -p "Очистить все log файлы? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi

swagger-gen: ## env: Сгенерировать актуальную Swagger спецификацию
	@docker compose run --rm swagger \
		init \
		-g cmd/todoapp/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

ps: ## env: Посмотреть запущенные Docker Compose сервисы
	@docker compose ps

migrate-create: ## PostgreSQL: Создать новую версию схемы данных
	@if [ -z "$(seq)" ]; then \
		echo "Отсутсвует необходимый параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up: ## PostgreSQL: Накатить миграции
	@make migrate-action action=up

migrate-down: ## PostgreSQL: Откатить миграции
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутсвует необходимый параметр action. Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

todoapp-run: ## Golang приложение: Запустить локально на хост-системе (для локальной разработки)
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	export REDIS_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

todoapp-deploy: ## Golang приложение: Запустить в Docker Compose сервисе (для деплоя)
	@docker compose up -d --build todoapp

todoapp-undeploy: ## Golang приложение: Остановить Docker Compose сервис
	@docker compose down todoapp

web-deploy: ## Web-Сервер: Запустить Caddy в Docker Compose сервисе
	@docker compose up -d --build web-server

web-undeploy: ## Web-Сервер: Остановить Docker Compose сервис
	@docker compose down web-server

load-test: ## Тесты: нагрузочное тестирование
	@go run scripts/load_test/main.go \
		-users 10 \
		-tasks-per-user 1000 \
		-concurrency 100 \
		-phase-duration 30s \
		-read-burst 50 \
		-mixed-reads 10 \
		-mixed-writes 1 \
		-report ${PROJECT_ROOT}/out/load_test/result.txt

help: ## Показать справку по командам
	@echo "=== Центр управления проектом ==="
	@echo ""
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
