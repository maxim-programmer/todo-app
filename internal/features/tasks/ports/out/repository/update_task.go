package ports_out_repository

import "github.com/nilchan-social/golang-todoapp/internal/core/domain"

// UpdateTaskParams передаёт полный снимок задачи для обновления.
// Хранилище использует task.Version для оптимистичной блокировки:
//
//	UPDATE ... WHERE id=$N AND version=$M
type UpdateTaskParams struct {
	Task domain.Task
}

func NewUpdateTaskParams(
	task domain.Task,
) UpdateTaskParams {
	return UpdateTaskParams{
		Task: task,
	}
}

type UpdateTaskResult struct {
	Task domain.Task
}

func NewUpdateTaskResult(
	task domain.Task,
) UpdateTaskResult {
	return UpdateTaskResult{
		Task: task,
	}
}
