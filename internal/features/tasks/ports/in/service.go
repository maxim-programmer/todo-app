package ports_in

import (
	"context"
)

// TasksService — входящий порт (ports/in) фичи tasks.
//
// В гексагональной архитектуре входящий порт — это контракт, который
// УРОВЕНЬ СЕРВИСА предоставляет внешнему миру. Интерфейс определяется ядром (сервисом),
// а не каким-либо адаптером: HTTP-хендлер адаптируется под этот интерфейс, а не наоборот.
//
// Это обеспечивает инверсию зависимости:
//
//	[HTTP-адаптер] → TasksService (интерфейс, владелец — ядро)
//	                      ↑ реализует
//	                 tasks_service.TasksService (конкретная реализация)
//
// Добавить gRPC или CLI-адаптер можно без единого изменения в ядре (уровне сервиса).
type TasksService interface {
	CreateTask(
		ctx context.Context,
		in CreateTaskParams,
	) (CreateTaskResult, error)

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

	PatchTask(
		ctx context.Context,
		in PatchTaskParams,
	) (PatchTaskResult, error)
}
