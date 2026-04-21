package tasks_service

import (
	"context"
	"fmt"

	ports_in "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/in"
	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// GetTask — use case получения задачи по ID.
//
// Бизнес-логика здесь отсутствует: сервис выступает транслятором между
// входящим портом (params от адаптера) и исходящим портом (запрос к хранилищу).
// Это нормально — не каждый use case содержит сложную логику.
// Ценность сервисного слоя здесь в том, что адаптеры не говорят друг с другом напрямую.
func (s *TasksService) GetTask(
	ctx context.Context,
	params ports_in.GetTaskParams,
) (ports_in.GetTaskResult, error) {
	repoParams := ports_out_repository.NewGetTaskParams(params.ID)
	repoResult, err := s.tasksRepository.GetTask(ctx, repoParams)
	if err != nil {
		return ports_in.GetTaskResult{}, fmt.Errorf("get task from repository: %w", err)
	}

	return ports_in.NewGetTaskResult(
		repoResult.Task,
	), nil
}
