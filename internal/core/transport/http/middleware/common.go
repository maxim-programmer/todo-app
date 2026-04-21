// Package core_http_middleware содержит HTTP middleware (промежуточные обработчики).
//
// Middleware — это функция-обёртка над http.Handler.
// Цепочка middleware выполняется последовательно: первый зарегистрированный
// middleware выполняется первым (внешний слой «луковицы»).
//
// Архитектурная схема:
//
//	Request → CORS → RequestID → Logger → Trace → Panic → Handler → Response
package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/nilchan-social/golang-todoapp/internal/core/logger"
	core_http_response "github.com/nilchan-social/golang-todoapp/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const (
	// requestIDHeader — заголовок для передачи уникального идентификатора запроса.
	// Позволяет отслеживать один запрос через все логи и системы.
	requestIDHeader = "X-Request-ID"
)

// CORS — middleware для обработки Cross-Origin Resource Sharing.
// Добавляет заголовки Access-Control-Allow-* только для разрешённых origins.
//
// Используем map для O(1) поиска вместо O(n) перебора списка.
// Preflight-запросы (OPTIONS) обрабатываются сразу без передачи дальше.
func CORS(allowedOriginsList []string) Middleware {
	allowedOrigins := make(map[string]struct{})
	for _, origin := range allowedOriginsList {
		allowedOrigins[origin] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			// Preflight OPTIONS — браузер проверяет разрешение CORS перед основным запросом.
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequestID — middleware, обеспечивающий каждый запрос уникальным идентификатором.
// Если клиент передаёт X-Request-ID — используем его (полезно для распределённой трассировки).
// Иначе генерируем новый X-Request-ID.
// Идентификатор добавляется и в заголовок ответа, чтобы клиент мог его использовать.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// Logger — middleware, кладущий логгер в контекст запроса.
// Обогащает логгер полями request_id и url, чтобы все последующие
// обработчики автоматически логировали эти поля.
//
// Важно: этот middleware должен идти ПОСЛЕ RequestID, чтобы request_id уже был доступен.
func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			// Создаём дочерний логгер с дополнительными полями.
			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Trace — middleware для логирования входящих запросов и времени их обработки.
// Использует ResponseWriter-обёртку, чтобы перехватить статус-код ответа,
// который иначе недоступен после вызова WriteHeader.
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			// Оборачиваем ResponseWriter, чтобы иметь возможность прочитать статус-код.
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Now().Sub(before)),
			)
		})
	}
}

// Panic — middleware для перехвата паник и возврата HTTP 500.
// Без этого middleware паника в обработчике уронила бы всю горутину,
// а стандартная библиотека Go вернула бы пустой ответ клиенту.
//
// Использует defer + recover — стандартный паттерн обработки паник в Go.
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
