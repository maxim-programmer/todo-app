package core_http_middleware

import (
	"fmt"
	"net/http"

	core_logger "github.com/nilchan-social/golang-todoapp/internal/core/logger"
)

// Dummy — учебный middleware-заглушка для демонстрации принципа работы middleware.
// Логирует вход и выход из middleware с переданной строкой s.
//
// Пример использования: добавить в APIVersionRouter для отладки цепочки middleware.
// Не используется в production-сборке.
func Dummy(s string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			// Этот код выполняется ДО обработчика
			log.Debug(fmt.Sprintf("-> before: %s", s))

			next.ServeHTTP(w, r)

			// Этот код выполняется ПОСЛЕ обработчика — классический паттерн «вокруг» (around).
			log.Debug(fmt.Sprintf("<- after: %s", s))
		})
	}
}
