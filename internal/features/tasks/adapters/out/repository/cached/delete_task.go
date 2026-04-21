package adapters_out_repository_cached

import (
	"context"
	"fmt"

	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// DeleteTask удаляет задачу из основного хранилища и инвалидирует кеш.
//
// Для инвалидации нужен AuthorUserID (чтобы удалить ключ "tasks:<userID>").
// Если задача есть в кеше — берём оттуда без лишнего запроса к БД.
// Если нет — делаем GetTask к основному хранилищу перед удалением.
func (r *CachedRepository) DeleteTask(
	ctx context.Context,
	params ports_out_repository.DeleteTaskParams,
) (ports_out_repository.DeleteTaskResult, error) {
	task, ok := r.getTaskFromCache(ctx, params.ID)
	if !ok {
		var (
			err error
		)

		repoGetTaskResult, err := r.mainRepository.GetTask(
			ctx,
			ports_out_repository.NewGetTaskParams(params.ID),
		)
		if err != nil {
			return ports_out_repository.DeleteTaskResult{}, fmt.Errorf("get task info: %w", err)
		}

		task = repoGetTaskResult.Task
	}

	r.invalidateTasks(ctx, task.AuthorUserID, &params.ID)

	return r.mainRepository.DeleteTask(
		ctx,
		ports_out_repository.NewDeleteTaskParams(params.ID),
	)
}
