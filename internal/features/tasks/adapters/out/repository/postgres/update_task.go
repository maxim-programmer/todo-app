package adapters_out_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/nilchan-social/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/postgres/pool"

	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// UpdateTask обновляет задачу в БД с оптимистичной блокировкой через version.
//
// Условие WHERE id=$5 AND version=$6 гарантирует атомарность:
// если другой запрос успел изменить запись между GetTask и UpdateTask,
// version в БД будет отличаться — RETURNING вернёт 0 строк → ErrNoRows → ErrConflict.
func (r *TasksRepository) UpdateTask(
	ctx context.Context,
	params ports_out_repository.UpdateTaskParams,
) (ports_out_repository.UpdateTaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.tasks
	SET
		title=$1,
		description=$2,
		completed=$3,
		completed_at=$4,
		version=version + 1

	WHERE id=$5 AND version=$6

	RETURNING
		id,
		version,
		title,
		description,
		completed,
		created_at,
		completed_at,
		author_user_id;
	`

	task := params.Task
	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		task.ID,
		task.Version,
	)

	var taskModel TaskModel
	if err := taskModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return ports_out_repository.UpdateTaskResult{}, fmt.Errorf(
				"task with id='%s' concurrently accessed: %w",
				task.ID,
				core_errors.ErrConflict,
			)
		}

		return ports_out_repository.UpdateTaskResult{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := modelToDomain(taskModel)

	return ports_out_repository.NewUpdateTaskResult(
		taskDomain,
	), nil
}
