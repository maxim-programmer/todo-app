package ports_out_repository

import (
	"context"
)

// TasksRepository — исходящий порт (ports/out) хранилища задач.
//
// В гексагональной архитектуре исходящий порт — это контракт, который
// УРОВЕНЬ СЕРВИСА ТРЕБУЕТ от инфраструктуры. Интерфейс определяется ядром (уровнем сервиса),
// а адаптер хранилища адаптируется под него — отсюда «инверсия зависимостей»:
//
//	Ядро → TasksRepository (интерфейс, владелец — ядро)
//	                ↑ реализует
//	    [Postgres-адаптер / Cached-адаптер / in-memory заглушка]
//
// Это позволяет подменять хранилище (Postgres → DynamoDB, добавить Redis-кеш)
// без единого изменения в бизнес-логике сервиса.
type TasksRepository interface {
	SaveTask(
		ctx context.Context,
		in SaveTaskParams,
	) (SaveTaskResult, error)

	GetTasks(
		ctx context.Context,
		in GetTasksParams,
	) (GetTasksResult, error)

	GetTask(
		ctx context.Context,
		in GetTaskParams,
	) (GetTaskResult, error)

	DeleteTask(
		ctx context.Context,
		in DeleteTaskParams,
	) (DeleteTaskResult, error)

	UpdateTask(
		ctx context.Context,
		in UpdateTaskParams,
	) (UpdateTaskResult, error)
}
