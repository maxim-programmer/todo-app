// Package statistics_transport_http содержит HTTP-обработчики для фичи статистики.
package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
	core_http_server "github.com/nilchan-social/golang-todoapp/internal/core/transport/http/server"
)

// StatisticsHTTPHandler — HTTP-обработчик для операций со статистикой.
type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

// StatisticsService — интерфейс сервиса статистики.
type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *uuid.UUID,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

// NewStatisticsHTTPHandler создаёт обработчик статистики с внедрённым сервисом.
func NewStatisticsHTTPHandler(
	statisticsService StatisticsService,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statisticsService,
	}
}

// Routes возвращает маршруты REST API для статистики.
func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
