package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/nilchan-social/golang-todoapp/internal/core/transport/http/middleware"
)

// Route описывает один HTTP-маршрут: метод, путь, обработчик и middleware.
// Middleware в Route применяются только к этому конкретному маршруту,
// в отличие от middleware сервера (применяются ко всем маршрутам).
type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []core_http_middleware.Middleware
}

// WithMiddleware применяет middleware маршрута к обработчику и возвращает готовый http.Handler.
func (r *Route) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r.Handler,
		r.Middleware...,
	)
}
