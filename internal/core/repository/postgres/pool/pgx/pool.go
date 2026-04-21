package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	core_postgres_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/postgres/pool"
)

// Pool — конкретная реализация интерфейса core_postgres_pool.Pool
// на базе pgxpool.Pool (пул соединений pgx).
//
// Встраивание *pgxpool.Pool даёт доступ ко всем методам pgx,
// но мы переопределяем Query/QueryRow/Exec, чтобы:
//  1. Оборачивать результаты в наши интерфейсы (Rows, Row, CommandTag)
//  2. Скрывать зависимость от pgx от слоя репозиториев
type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

// NewPool создаёт и проверяет пул соединений с PostgreSQL.
// Ping() при инициализации гарантирует, что БД доступна до начала работы сервера.
func NewPool(
	ctx context.Context,
	config Config,
) (*Pool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	// Проверяем доступность БД сразу при старте.
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}

	return &Pool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

// Query выполняет SELECT-запрос и возвращает интерфейс Rows.
// Оборачивает pgx.Rows в pgxRows для соответствия интерфейсу пула.
func (p *Pool) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (core_postgres_pool.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxRows{rows}, nil
}

// QueryRow выполняет запрос, возвращающий одну строку.
// Оборачивает pgx.Row в pgxRow (с маппингом ошибок).
func (p *Pool) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) core_postgres_pool.Row {
	row := p.Pool.QueryRow(ctx, sql, args...)

	return pgxRow{row}
}

// Exec выполняет DML-запрос (INSERT, UPDATE, DELETE).
// Оборачивает pgconn.CommandTag в pgxCommandTag.
func (p *Pool) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (core_postgres_pool.CommandTag, error) {
	tag, err := p.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}

	return pgxCommandTag{tag}, nil
}

// OpTimeout возвращает максимальное время выполнения одного запроса.
func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}
