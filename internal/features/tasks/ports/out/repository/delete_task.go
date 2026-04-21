package ports_out_repository

import "github.com/google/uuid"

// DeleteTaskParams / DeleteTaskResult — параметры и результат исходящего порта для удаления задачи.
// Ядро формулирует намерение («удали задачу с этим ID») в терминах домена.
// Детали реализации (DELETE FROM ... / invalidate cache) остаются за адаптером.
type DeleteTaskParams struct {
	ID uuid.UUID
}

func NewDeleteTaskParams(
	id uuid.UUID,
) DeleteTaskParams {
	return DeleteTaskParams{
		ID: id,
	}
}

type DeleteTaskResult struct{}

func NewDeleteTaskResult() DeleteTaskResult {
	return DeleteTaskResult{}
}
