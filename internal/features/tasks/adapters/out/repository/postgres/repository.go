package adapters_out_repository_postgres

import core_postgres_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/postgres/pool"

// TasksRepository — исходящий адаптер (adapters/out) для PostgreSQL.
//
// В гексагональной архитектуре исходящий адаптер — это «правая сторона» гексагона.
// Его задача: реализовать контракт исходящего порта (ports/out: ports_out_repository.TasksRepository),
// переводя вызовы уровня сервиса в конкретные операции с инфраструктурой.
//
// Чего адаптер НЕ делает:
//   - не содержит бизнес-логики
//   - не знает о существовании HTTP-адаптера или других исходящих адаптеров
type TasksRepository struct {
	pool core_postgres_pool.Pool
}

func NewTasksRepository(
	pool core_postgres_pool.Pool,
) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}
