package core_redis_pool

import (
	"context"
	"time"
)

// Pool — абстракция над клиентом Redis.
// Инкапсулирует только команды, используемые приложением, и изолирует
// доменный код от конкретной библиотеки (go-redis, redigo и т.д.).
type Pool interface {
	Get(ctx context.Context, key string) StringCmd
	Set(ctx context.Context, key string, value any, ttl time.Duration) StatusCmd
	Del(ctx context.Context, keys ...string) IntCmd
	HGet(ctx context.Context, key string, field string) StringCmd
	HSet(ctx context.Context, key string, values ...any) IntCmd
	Close() error

	TTL() time.Duration
}

// StringCmd — результат команды, возвращающей строку (GET, HGET).
type StringCmd interface {
	Bytes() ([]byte, error)
}

// StatusCmd — результат команды, возвращающей статус операции (SET).
type StatusCmd interface {
	Err() error
}

// IntCmd — результат команды, возвращающей целое число (DEL, HSET).
type IntCmd interface {
	Err() error
	Val() int64
}
