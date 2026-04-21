package core_redis_pool

import "errors"

var (
	// NotFound — сигнальная ошибка: ключ или поле отсутствует в Redis.
	// Адаптеры конкретных клиентов (goredis и др.) должны отображать свои
	// внутренние ошибки "not found" (например, redis.Nil) в эту ошибку,
	// чтобы код приложения не зависел от деталей библиотеки.
	NotFound = errors.New("not found")
)
