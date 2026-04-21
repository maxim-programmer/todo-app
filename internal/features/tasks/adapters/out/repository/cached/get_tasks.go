package adapters_out_repository_cached

import (
	"context"
	"errors"

	core_logger "github.com/nilchan-social/golang-todoapp/internal/core/logger"
	core_redis_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/redis/pool"
	"go.uber.org/zap"

	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// GetTasks реализует cache-aside для списка задач.
//
// Кеш организован как Redis hash:
//   - key   = "tasks:<userID>" или "tasks:all" (если userID == nil)
//   - field = "<limit>:<offset>", например "10:0", "nil:nil"
//
// Такая структура позволяет хранить несколько вариантов пагинации под одним ключом для конкретного пользователя,
// и инвалидировать все варианты одним DEL по ключу hash'а.
//
// Ошибки сериализации/кеша логируются, но не прерывают запрос —
// при любом сбое данные будут получены из основного хранилища (graceful degradation).
func (r *CachedRepository) GetTasks(
	ctx context.Context,
	params ports_out_repository.GetTasksParams,
) (ports_out_repository.GetTasksResult, error) {
	log := core_logger.FromContext(ctx)

	key := tasksListKey(params.UserID)
	field := tasksListField(params.Limit, params.Offset)

	bytes, err := r.pool.HGet(ctx, key, field).Bytes()
	if err != nil {
		if !errors.Is(err, core_redis_pool.NotFound) {
			log.Error("hget task list", zap.Error(err))
		}
	} else {
		var taskListModel TaskListModel
		if err := taskListModel.Deserialize(bytes); err != nil {
			log.Error("deserialize cached task list", zap.Error(err))
		} else {
			tasks := modelToDomains(taskListModel)

			return ports_out_repository.NewGetTasksResult(
				tasks,
			), nil
		}
	}

	repoGetTasksResult, err := r.mainRepository.GetTasks(
		ctx,
		ports_out_repository.NewGetTasksParams(
			params.UserID,
			params.Limit,
			params.Offset,
		),
	)
	if err != nil {
		return ports_out_repository.GetTasksResult{}, err
	}

	taskListModel := domainsToModel(repoGetTasksResult.Tasks)
	bytes, err = taskListModel.Serialize()
	if err != nil {
		log.Error("serialize task list", zap.Error(err))
	} else {
		if err := r.pool.HSet(ctx, key, field, bytes).Err(); err != nil {
			log.Error("hset task list in cache", zap.Error(err))
		}
	}

	return repoGetTasksResult, nil
}
