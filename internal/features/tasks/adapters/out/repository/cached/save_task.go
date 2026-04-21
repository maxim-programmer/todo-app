package adapters_out_repository_cached

import (
	"context"

	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// SaveTask сохраняет задачу в основное хранилище, затем обновляет кеш:
//   - cacheTask: кеширует новую задачу по её ID (write-through для одиночной задачи)
//   - invalidateTasks: инвалидирует все СПИСКИ задач автора, т.к. они устарели.
func (r *CachedRepository) SaveTask(
	ctx context.Context,
	params ports_out_repository.SaveTaskParams,
) (ports_out_repository.SaveTaskResult, error) {
	repoSaveTaskResult, err := r.mainRepository.SaveTask(
		ctx,
		ports_out_repository.NewSaveTaskParams(params.Task),
	)
	if err != nil {
		return ports_out_repository.SaveTaskResult{}, err
	}

	task := repoSaveTaskResult.Task

	r.cacheTask(ctx, task)
	//     taskID == nil означает «не удалять ключ конкретной задачи» — она только что создана.
	r.invalidateTasks(ctx, task.AuthorUserID, nil)

	return repoSaveTaskResult, nil
}
