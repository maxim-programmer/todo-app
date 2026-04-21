package ports_out_repository

import (
	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// GetTaskParams / GetTaskResult — параметры и результат исходящего порта для чтения задачи.
// Ядро запрашивает задачу по UUID — в терминах домена, не SQL.
// Адаптер репозитория транслирует это в SELECT, GET или чтение из кеша.
type GetTaskParams struct {
	ID uuid.UUID
}

func NewGetTaskParams(
	id uuid.UUID,
) GetTaskParams {
	return GetTaskParams{
		ID: id,
	}
}

type GetTaskResult struct {
	Task domain.Task
}

func NewGetTaskResult(
	task domain.Task,
) GetTaskResult {
	return GetTaskResult{
		Task: task,
	}
}
