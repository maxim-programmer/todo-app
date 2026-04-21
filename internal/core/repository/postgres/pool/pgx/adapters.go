// Package core_pgx_pool содержит конкретную реализацию интерфейса Pool
// на базе библиотеки jackc/pgx/v5 — одного из самых популярных
// и производительных драйверов PostgreSQL для Go.
package core_pgx_pool

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/postgres/pool"
)

// pgxRows оборачивает pgx.Rows для реализации интерфейса core_postgres_pool.Rows.
// Встраивание pgx.Rows даёт все методы (Next, Close, Err, Scan) «бесплатно».
type pgxRows struct {
	pgx.Rows
}

// pgxRow оборачивает pgx.Row для реализации интерфейса core_postgres_pool.Row.
// Переопределяем Scan, чтобы преобразовывать pgx-ошибки в наши типизированные ошибки.
type pgxRow struct {
	pgx.Row
}

// Scan вызывает оригинальный pgx Scan и преобразует ошибки через mapErrors.
func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}

	return nil
}

// pgxCommandTag оборачивает pgconn.CommandTag для реализации интерфейса CommandTag.
type pgxCommandTag struct {
	pgconn.CommandTag
}

// mapErrors преобразует специфические ошибки pgx в типизированные ошибки пакета pool.
// Это «антикоррупционный слой» (Anti-Corruption Layer) — изолирует детали pgx
// от остального кода приложения.
func mapErrors(err error) error {
	const (
		// Код PostgreSQL для ошибки нарушения внешнего ключа.
		// Полный список кодов: https://www.postgresql.org/docs/current/errcodes-appendix.html
		pgxViolatesForeignKeyErrorCode = "23503"
	)

	// pgx.ErrNoRows → наш ErrNoRows (запись не найдена)
	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	// Проверяем, является ли ошибка структурированной PostgreSQL-ошибкой.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgxViolatesForeignKeyErrorCode {
			return fmt.Errorf(
				"%v: %w",
				err,
				core_postgres_pool.ErrViolatesForeignKey,
			)
		}
	}

	// Все остальные ошибки оборачиваем в ErrUnknown.
	return fmt.Errorf(
		"%v: %w",
		err,
		core_postgres_pool.ErrUnknown,
	)
}
