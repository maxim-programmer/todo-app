package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/nilchan-social/golang-todoapp/internal/core/transport/http/middleware"
)

// APIVersion — тип для идентификатора версии API.
// Использование отдельного типа (не string) делает код явным и защищает от опечаток.
type APIVersion string

var (
	ApiVersion1 = APIVersion("v1")
	ApiVersion2 = APIVersion("v2")
	ApiVersion3 = APIVersion("v3")
)

// APIVersionRouter группирует маршруты под единым версионным префиксом /api/v1, /api/v2 и т.д.
// Поддерживает собственные middleware, применяемые только к маршрутам (Route) этой версии API.
// Это позволяет, например, добавить аутентификацию только для /api/v2, не трогая /api/v1.
type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVersion
	routes     []Route
	middleware []core_http_middleware.Middleware
}

// NewAPIVersionRouter создаёт роутер для заданной версии API.
// Необязательные middleware будут применяться ко всем маршрутам этой версии.
func NewAPIVersionRouter(
	apiVersion APIVersion,
	middleware ...core_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

// AddRoutes добавляет маршруты в роутер.
func (r *APIVersionRouter) AddRoutes(routes ...Route) {
	r.routes = append(r.routes, routes...)
}

// Handlers формирует мапу «паттерн маршрута → обработчик» для регистрации в http.ServeMux.
// Паттерн строится как: "METHOD /api/v1/path".
// Middleware роутера (APIVersionRouter) оборачивают middleware маршрута (Route) снаружи.
func (r *APIVersionRouter) Handlers() map[string]http.Handler {
	handlers := make(map[string]http.Handler, len(r.routes))

	for _, route := range r.routes {
		// Формируем полный паттерн: "GET /api/v1/tasks", "POST /api/v1/users" и т.д.
		pattern := route.Method + " /api/" + string(r.apiVersion) + route.Path
		handler := core_http_middleware.ChainMiddleware(
			route.WithMiddleware(),
			r.middleware...,
		)

		handlers[pattern] = handler
	}

	return handlers
}
