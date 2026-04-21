package ports_in

import "github.com/google/uuid"

// DeleteTaskParams / DeleteTaskResult — параметры и результат входящего порта TaskService.DeleteTask.
// DeleteTaskResult намеренно пустой: операция удаления не возвращает данных.
// Явный пустой тип выбран вместо error-only сигнатуры для единообразия интерфейса.
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
