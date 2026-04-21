package statistics_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
	core_postgres_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/postgres/pool"
)

// TaskModel — модель задачи для репозитория статистики.
// Идентичен TaskModel из tasks/repository, но дублируется намеренно:
// пакеты независимы, и связывать их через общую структуру создало бы лишнюю зависимость.
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

// Scan заполняет поля модели из результата запроса к БД.
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

// modelToDomain конвертирует модель БД в доменный объект Task.
func modelToDomain(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}

// modelsToDomains конвертирует список моделей в список доменных объектов.
func modelsToDomains(taskModels []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(taskModels))

	for i, model := range taskModels {
		domains[i] = modelToDomain(model)
	}

	return domains
}
