package core_postgres_pool

import "errors"

// Sentinel-ошибки уровня пула соединений.
// Адаптер pgx преобразует специфические ошибки pgx в эти типизированные ошибки,
// позволяя репозиториям работать с ними без зависимости от pgx.
var (
	// ErrNoRows — запрос не вернул строк.
	// Репозиторий преобразует это в core_errors.ErrNotFound.
	ErrNoRows = errors.New("no rows")

	// ErrViolatesForeignKey — нарушение ограничения внешнего ключа (PostgreSQL код 23503).
	// Например: попытка создать задачу с несуществующим author_user_id.
	ErrViolatesForeignKey = errors.New("violates foreign key")

	// ErrUnknown — любая другая ошибка базы данных, не обработанная явно.
	ErrUnknown = errors.New("unknown")
)
