package adapters_out_repository_postgres

import (
	"context"
	"fmt"

	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// GetTasks выбирает список задач с пагинацией и опциональной фильтрацией по автору.
// Если UserID задан — добавляется WHERE author_user_id=$3.
// Limit и Offset передаются как $1/$2; nil допустим — Postgres интерпретирует
// NULL LIMIT/OFFSET как отсутствие ограничения.
func (r *TasksRepository) GetTasks(
	ctx context.Context,
	params ports_out_repository.GetTasksParams,
) (ports_out_repository.GetTasksResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	%s
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2;
	`

	args := []any{params.Limit, params.Offset}

	if params.UserID != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id=$3")
		args = append(args, params.UserID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return ports_out_repository.GetTasksResult{}, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		var taskModel TaskModel

		if err := taskModel.Scan(rows); err != nil {
			return ports_out_repository.GetTasksResult{}, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return ports_out_repository.GetTasksResult{}, fmt.Errorf("next rows: %w", err)
	}

	taskDomains := modelsToDomains(taskModels)

	return ports_out_repository.NewGetTasksResult(
		taskDomains,
	), nil
}
