package adapters_out_repository_cached

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
	core_logger "github.com/nilchan-social/golang-todoapp/internal/core/logger"
	core_redis_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/redis/pool"
	"go.uber.org/zap"
)

// getTaskFromCache читает задачу из Redis по ключу "task:<id>".
// Возвращает (task, true) при cache hit, (zero, false) при промахе или ошибке.
// Ошибки логируются; при промахе вызывающий код обращается к основному хранилищу.
func (r *CachedRepository) getTaskFromCache(
	ctx context.Context,
	id uuid.UUID,
) (domain.Task, bool) {
	log := core_logger.FromContext(ctx)

	key := taskKey(id)

	bytes, err := r.pool.Get(ctx, key).Bytes()
	if err != nil {
		if !errors.Is(err, core_redis_pool.NotFound) {
			log.Error("read from cache", zap.Error(err))
		}

		return domain.Task{}, false
	}

	var taskModel TaskModel
	if err := taskModel.Deserialize(bytes); err != nil {
		log.Error("deserialize cached task", zap.Error(err))

		return domain.Task{}, false
	}

	taskDomain := modelToDomain(taskModel)

	return taskDomain, true
}

// cacheTask сериализует задачу в JSON и сохраняет в Redis с TTL пула.
// Ошибки логируются и не прерывают работу (best-effort caching).
func (r *CachedRepository) cacheTask(
	ctx context.Context,
	task domain.Task,
) {
	log := core_logger.FromContext(ctx)

	taskModel := domainToModel(task)
	bytes, err := taskModel.Serialize()
	if err != nil {
		log.Error("serialize task", zap.Error(err))
	} else {
		if err := r.pool.Set(
			ctx,
			taskKey(taskModel.ID),
			bytes,
			r.pool.TTL(),
		).Err(); err != nil {
			log.Error("set task in cache", zap.Error(err))
		}
	}
}

// invalidateTasks удаляет из Redis ключи, связанные со списками задач:
//   - "tasks:all"        — список всех задач
//   - "tasks:<userID>"   — список задач конкретного пользователя
//   - "task:<id>"        — ключ одиночной задачи, только если taskID != nil
func (r *CachedRepository) invalidateTasks(
	ctx context.Context,
	userID uuid.UUID,
	taskID *uuid.UUID,
) {
	log := core_logger.FromContext(ctx)

	invalidateKeys := []string{
		tasksListKey(nil),
		tasksListKey(&userID),
	}
	if taskID != nil {
		invalidateKeys = append(invalidateKeys, taskKey(*taskID))
	}

	if err := r.pool.Del(ctx, invalidateKeys...).Err(); err != nil {
		log.Error("invalidate cached tasks lists", zap.Error(err))
	}
}
