package core_http_middleware

import "net/http"

// Middleware — тип, описывающий функцию, которая принимает http.Handler
// и возвращает http.Handler (оборачивает его).
type Middleware func(http.Handler) http.Handler

// ChainMiddleware применяет цепочку middleware к обработчику h.
// Middleware применяются в порядке объявления: первый в списке будет выполнен первым.
//
// Реализация: обходим массив с конца и оборачиваем обработчик снаружи внутрь,
// чтобы первый middleware оказался самым внешним слоем.
func ChainMiddleware(
	h http.Handler,
	m ...Middleware,
) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}
