package adapters_out_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/nilchan-social/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/postgres/pool"

	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// SaveTask вставляет задачу в БД и возвращает сохранённую строку через RETURNING.
// При нарушении внешнего ключа (author_user_id → users) возвращает ErrNotFound,
// который HTTP-слой транслирует в 404.
func (r *TasksRepository) SaveTask(
	ctx context.Context,
	params ports_out_repository.SaveTaskParams,
) (ports_out_repository.SaveTaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.tasks (id, version, title, description, completed, created_at, completed_at, author_user_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
	`

	task := params.Task
	row := r.pool.QueryRow(
		ctx,
		query,
		task.ID,
		task.Version,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)
	var taskModel TaskModel
	if err := taskModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return ports_out_repository.SaveTaskResult{}, fmt.Errorf(
				"%v: user with id='%s': %w",
				err,
				task.AuthorUserID,
				core_errors.ErrNotFound,
			)
		}

		return ports_out_repository.SaveTaskResult{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := modelToDomain(taskModel)

	return ports_out_repository.NewSaveTaskResult(
		taskDomain,
	), nil
}
