package tasks_service

import (
	"context"
	"fmt"

	ports_in "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/in"
	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// DeleteTask — use case удаления задачи.
// Транслирует вызов через исходящий порт; адаптер хранилища вернёт ErrNotFound,
// если задача не существует, — ядро (уровень сервиса) пробросит его в ответ входящему адаптеру.
func (s *TasksService) DeleteTask(
	ctx context.Context,
	params ports_in.DeleteTaskParams,
) (ports_in.DeleteTaskResult, error) {
	repoParams := ports_out_repository.NewDeleteTaskParams(params.ID)
	_, err := s.tasksRepository.DeleteTask(ctx, repoParams)
	if err != nil {
		return ports_in.DeleteTaskResult{}, fmt.Errorf("delete task from repository: %w", err)
	}

	return ports_in.NewDeleteTaskResult(), nil
}
