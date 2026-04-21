package adapters_out_repository_cached

import (
	core_redis_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/redis/pool"
	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// CachedRepository — исходящий адаптер (adapters/out) с кешированием, реализующий Decorator pattern.
//
// В гексагональной архитектуре это ещё один исходящий адаптер, реализующий тот же
// ports_out_repository.TasksRepository. Ядро (уровень сервиса) не будет знать, что перед ним кешированная
// обёртка — сервис работает с интерфейсом, который одинаков для Postgres и Cached.
//
// Decorator позволяет добавить Redis-кеш без изменения уровня сервиса и Postgres-адаптера.
// В main.go достаточно обернуть один адаптер в другой:
//
//	tasks_cached.NewCachedRepository(redisPool, tasks_postgres.NewTasksRepository(pool))
//
// Для отключения кеша — убрать обёртку, ядро (уровень сервиса) не заметит разницы.
//
// Стратегии кеширования:
//   - Одиночная задача (ключ "task:<id>"): cache-aside с GET/SET
//   - Список задач (ключ "tasks:<userID>" / "tasks:all"): cache-aside через Redis hash,
//     где поле (fields) хеша кодирует пагинацию (<limit>:<offset>)
//
// Инвалидация происходит при любой мутации: SaveTask, UpdateTask, DeleteTask.
// Ошибки кеша логируются и не прерывают запрос — graceful degradation к основному хранилищу.
type CachedRepository struct {
	pool           core_redis_pool.Pool
	mainRepository ports_out_repository.TasksRepository
}

func NewCachedRepository(
	pool core_redis_pool.Pool,
	mainRepository ports_out_repository.TasksRepository,
) *CachedRepository {
	return &CachedRepository{
		pool:           pool,
		mainRepository: mainRepository,
	}
}
