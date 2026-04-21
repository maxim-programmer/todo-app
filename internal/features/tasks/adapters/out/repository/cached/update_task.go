package adapters_out_repository_cached

import (
	"context"

	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// UpdateTask обновляет задачу в основном хранилище, затем:
//   - cacheTask: обновляет кеш задачи по ID (write-through, чтобы следующий
//     GetTask не вернул устаревшие данные)
//   - invalidateTasks: инвалидирует списки задач автора, т.к. данные в них устарели.
func (r *CachedRepository) UpdateTask(
	ctx context.Context,
	params ports_out_repository.UpdateTaskParams,
) (ports_out_repository.UpdateTaskResult, error) {
	repoUpdateTaskResult, err := r.mainRepository.UpdateTask(
		ctx,
		ports_out_repository.NewUpdateTaskParams(params.Task),
	)
	if err != nil {
		return ports_out_repository.UpdateTaskResult{}, err
	}

	task := repoUpdateTaskResult.Task

	r.cacheTask(ctx, task)
	r.invalidateTasks(ctx, task.AuthorUserID, nil)

	return repoUpdateTaskResult, nil
}
