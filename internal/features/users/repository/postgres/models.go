package users_postgres_repository

import (
	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
	core_postgres_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/postgres/pool"
)

// UserModel — структура для маппинга строки таблицы `todoapp.users` в Go-тип.
// Порядок полей совпадает с порядком столбцов в SELECT-запросах репозитория.
type UserModel struct {
	ID          uuid.UUID
	Version     int
	FullName    string
	PhoneNumber *string
}

// Scan заполняет поля модели из результата запроса к БД.
func (m *UserModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.FullName,
		&m.PhoneNumber,
	)
}

// modelToDomain конвертирует модель БД в доменный объект.
func modelToDomain(model UserModel) domain.User {
	return domain.NewUser(
		model.ID,
		model.Version,
		model.FullName,
		model.PhoneNumber,
	)
}

// modelsToDomains конвертирует список моделей БД в список доменных объектов.
func modelsToDomains(models []UserModel) []domain.User {
	domains := make([]domain.User, len(models))

	for i, model := range models {
		domains[i] = modelToDomain(model)
	}

	return domains
}
