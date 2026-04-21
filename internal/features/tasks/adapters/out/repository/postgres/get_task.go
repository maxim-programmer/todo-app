package adapters_out_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/nilchan-social/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/postgres/pool"

	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// GetTask выбирает задачу по ID.
// ErrNoRows → ErrNotFound: задача с таким ID не существует.
func (r *TasksRepository) GetTask(
	ctx context.Context,
	params ports_out_repository.GetTaskParams,
) (ports_out_repository.GetTaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, params.ID)

	var taskModel TaskModel
	if err := taskModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return ports_out_repository.GetTaskResult{}, fmt.Errorf(
				"task with id='%s': %w",
				params.ID,
				core_errors.ErrNotFound,
			)
		}

		return ports_out_repository.GetTaskResult{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := modelToDomain(taskModel)

	return ports_out_repository.NewGetTaskResult(
		taskDomain,
	), nil
}
