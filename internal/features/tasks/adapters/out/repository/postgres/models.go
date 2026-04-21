package adapters_out_repository_postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
	core_postgres_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/postgres/pool"
)

// TaskModel — промежуточное представление задачи для сканирования строк pgx.
// Поля соответствуют столбцам таблицы todoapp.tasks.
// Description и CompletedAt — указатели, т.к. могут быть NULL в БД.
type TaskModel struct {
	ID           uuid.UUID
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID uuid.UUID
}

func (m *TaskModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.Title,
		&m.Description,
		&m.Completed,
		&m.CreatedAt,
		&m.CompletedAt,
		&m.AuthorUserID,
	)
}

func modelToDomain(model TaskModel) domain.Task {
	return domain.NewTask(
		model.ID,
		model.Version,
		model.Title,
		model.Description,
		model.Completed,
		model.CreatedAt,
		model.CompletedAt,
		model.AuthorUserID,
	)
}

func modelsToDomains(models []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(models))

	for i, model := range models {
		domains[i] = modelToDomain(model)
	}

	return domains
}
