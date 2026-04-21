package ports_in

import (
	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// GetTasksParams / GetTasksResult — параметры и результат входящего порта для списка задач.
// Принадлежат уровню сервиса; именно сервис решает, какие параметры фильтрации и пагинации
// он готово принять.
//
// Семантика опциональных полей (nil = «не задано»):
//   - UserID nil — задачи всех пользователей
//   - Limit nil — без ограничения на количество
//   - Offset nil — с нулевым смещением
type GetTasksParams struct {
	UserID *uuid.UUID
	Limit  *int
	Offset *int
}

func NewGetTasksParams(
	userID *uuid.UUID,
	limit *int,
	offset *int,
) GetTasksParams {
	return GetTasksParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}
}

type GetTasksResult struct {
	Tasks []domain.Task
}

func NewGetTasksResult(
	tasks []domain.Task,
) GetTasksResult {
	return GetTasksResult{
		Tasks: tasks,
	}
}
