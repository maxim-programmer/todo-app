package ports_in

import (
	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// PatchTaskParams / PatchTaskResult — параметры и результат входящего порта для частичного обновления.
//
// Patch — доменный тип (domain.TaskPatch), а не HTTP-специфичная структура.
// HTTP-адаптер сам отвечает за перевод своего PatchTaskRequest → domain.TaskPatch,
// чтобы уровень сервиса не знал ничего про JSON, nullable-поля или HTTP-заголовки.
type PatchTaskParams struct {
	ID    uuid.UUID
	Patch domain.TaskPatch
}

func NewPatchTaskParams(
	id uuid.UUID,
	patch domain.TaskPatch,
) PatchTaskParams {
	return PatchTaskParams{
		ID:    id,
		Patch: patch,
	}
}

type PatchTaskResult struct {
	Task domain.Task
}

func NewPatchTaskResult(
	task domain.Task,
) PatchTaskResult {
	return PatchTaskResult{
		Task: task,
	}
}
