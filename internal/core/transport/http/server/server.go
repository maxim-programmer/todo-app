package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nilchan-social/golang-todoapp/docs"
	core_logger "github.com/nilchan-social/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/nilchan-social/golang-todoapp/internal/core/transport/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

// HTTPServer — обёртка над стандартным net/http, добавляющая:
//   - Поддержку версионирования API через APIVersionRouter
//   - Цепочку middleware для всех маршрутов
//   - Swagger UI
//   - Graceful shutdown
type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	// middleware применяются ко всем маршрутам сервера (глобальные middleware).
	middleware []core_http_middleware.Middleware
}

// NewHTTPServer создаёт HTTP-сервер с заданными глобальными middleware.
func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

// RegisterAPIRouters регистрирует версионированные роутеры API в ServeMux.
// Каждый роутер добавляет маршруты вида "{METHOD} /api/v{N}/{path}".
func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		handlers := router.Handlers()

		for path, handler := range handlers {
			s.mux.Handle(path, handler)
		}
	}
}

// RegisterRoutes регистрирует маршруты без версионного префикса (например, главная страница "/").
func (s *HTTPServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		path := route.Method + " " + route.Path
		handler := route.WithMiddleware()

		s.mux.Handle(path, handler)
	}
}

// RegisterSwagger регистрирует два маршрута:
//   - GET /swagger/      — Swagger UI (интерактивная документация)
//   - GET /swagger/doc.json — спецификация OpenAPI в формате JSON
func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle(
		"GET /swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DefaultModelsExpandDepth(-1),
		),
	)

	s.mux.HandleFunc(
		"GET /swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

// Run запускает HTTP-сервер и блокирует выполнение до получения сигнала завершения.
//
// Graceful shutdown:
//  1. При отмене ctx (SIGINT/SIGTERM) вызывается server.Shutdown()
//  2. Shutdown ждёт завершения активных HTTP обработчиков до ShutdownTimeout
//  3. По истечении таймаута принудительно закрывает соединения через server.Close()
//
// Канал ch используется для передачи ошибки из горутины в основной поток.
func (s *HTTPServer) Run(ctx context.Context) error {
	// Применяем глобальные middleware к ServeMux.
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	// Буферизированный канал (размер 1), чтобы горутина не заблокировалась
	// при отправке ошибки, если основной поток уже ушёл в select.
	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		// http.ErrServerClosed — нормальное завершение после Shutdown(), не ошибка.
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		// HTTP сервер завершился из-за ошибки (например, порт занят).
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		// Получен сигнал завершения — выполняем graceful shutdown.
		s.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			// Если graceful shutdown не успел — принудительно закрываем.
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}
