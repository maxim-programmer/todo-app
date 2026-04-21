package adapters_out_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/nilchan-social/golang-todoapp/internal/core/errors"
	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// DeleteTask удаляет задачу по ID.
// Проверяет RowsAffected: если 0 — задача не существовала → ErrNotFound.
func (r *TasksRepository) DeleteTask(
	ctx context.Context,
	params ports_out_repository.DeleteTaskParams,
) (ports_out_repository.DeleteTaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.tasks
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, params.ID)
	if err != nil {
		return ports_out_repository.DeleteTaskResult{}, fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ports_out_repository.DeleteTaskResult{}, fmt.Errorf(
			"task with id='%s': %w",
			params.ID,
			core_errors.ErrNotFound,
		)
	}

	return ports_out_repository.NewDeleteTaskResult(), nil
}
