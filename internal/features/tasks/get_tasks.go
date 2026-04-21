package tasks_service

import (
	"context"
	"fmt"

	core_errors "github.com/nilchan-social/golang-todoapp/internal/core/errors"

	ports_in "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/in"
	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// GetTasks — use case получения списка задач с пагинацией.
//
// Валидация limit/offset живёт здесь, в ядре (сервисе), а не в HTTP-адаптере или репозитории.
// Это принципиально: правило «limit не может быть отрицательным» — бизнес-правило,
// оно не должно зависеть от того, пришёл запрос по HTTP или из очереди сообщений.
// Любой входящий адаптер (adapters/in) получит одинаковую проверку автоматически.
func (s *TasksService) GetTasks(
	ctx context.Context,
	params ports_in.GetTasksParams,
) (ports_in.GetTasksResult, error) {
	if params.Limit != nil && *params.Limit < 0 {
		return ports_in.GetTasksResult{}, fmt.Errorf(
			"limit must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if params.Offset != nil && *params.Offset < 0 {
		return ports_in.GetTasksResult{}, fmt.Errorf(
			"offset must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	repoParams := ports_out_repository.NewGetTasksParams(
		params.UserID,
		params.Limit,
		params.Offset,
	)
	repoResult, err := s.tasksRepository.GetTasks(ctx, repoParams)
	if err != nil {
		return ports_in.GetTasksResult{}, fmt.Errorf("get tasks from repository: %w", err)
	}

	return ports_in.NewGetTasksResult(
		repoResult.Tasks,
	), nil
}
