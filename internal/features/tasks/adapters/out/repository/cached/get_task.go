package adapters_out_repository_cached

import (
	"context"

	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// GetTask реализует паттерн cache-aside для одиночной задачи:
//  1. Cache hit: возвращаем задачу напрямую из Redis (ключ "task:<id>").
//  2. Cache miss: читаем из основного хранилища, кладём результат в кеш, возвращаем.
func (r *CachedRepository) GetTask(
	ctx context.Context,
	params ports_out_repository.GetTaskParams,
) (ports_out_repository.GetTaskResult, error) {
	if task, ok := r.getTaskFromCache(ctx, params.ID); ok {
		return ports_out_repository.NewGetTaskResult(
			task,
		), nil
	}

	repoGetTaskResult, err := r.mainRepository.GetTask(
		ctx,
		ports_out_repository.NewGetTaskParams(params.ID),
	)
	if err != nil {
		return ports_out_repository.GetTaskResult{}, err
	}

	r.cacheTask(ctx, repoGetTaskResult.Task)

	return repoGetTaskResult, nil
}
