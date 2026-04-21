package core_goredis_pool

import (
	"errors"

	core_redis_pool "github.com/nilchan-social/golang-todoapp/internal/core/repository/redis/pool"
	"github.com/redis/go-redis/v9"
)

// goredisStringCmd оборачивает *redis.StringCmd и реализует core_redis_pool.StringCmd.
// Переопределяет Bytes(), чтобы транслировать go-redis-специфичные ошибки
// в доменные через mapError.
type goredisStringCmd struct {
	*redis.StringCmd
}

func (c goredisStringCmd) Bytes() ([]byte, error) {
	data, err := c.StringCmd.Bytes()
	if err != nil {
		return nil, mapError(err)
	}

	return data, nil
}

// goredisStatusCmd оборачивает *redis.StatusCmd и реализует core_redis_pool.StatusCmd.
type goredisStatusCmd struct {
	*redis.StatusCmd
}

// goredisIntCmd оборачивает *redis.IntCmd и реализует core_redis_pool.IntCmd.
type goredisIntCmd struct {
	*redis.IntCmd
}

// mapError переводит go-redis-специфичную ошибку redis.Nil
// в доменную core_redis_pool.NotFound.
// Это позволяет коду приложения не импортировать и не проверять go-redis напрямую.
func mapError(err error) error {
	if errors.Is(err, redis.Nil) {
		return core_redis_pool.NotFound
	}

	return err
}
