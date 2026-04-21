package ports_out_repository

import (
	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// GetTasksParams описывает параметры выборки задач из хранилища.
// nil-значения интерпретируются как «без ограничения»:
//   - UserID nil — выбрать задачи всех пользователей
//   - Limit nil — без LIMIT
//   - Offset nil — без OFFSET
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
