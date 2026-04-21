// Package statistics_service содержит бизнес-логику расчёта статистики по задачам.
package statistics_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// StatisticsService — сервис статистики.
// Читает задачи из репозитория и делегирует подсчёт в domain.CreateStatistics.
type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

// StatisticsRepository — интерфейс репозитория для сервиса статистики.
// Сервис требует только метод получения задач — остальные операции ему не нужны.
// Это хороший пример принципа Interface Segregation (ISP из SOLID):
// интерфейс содержит только то, что действительно нужно потребителю.
type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userID *uuid.UUID,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

// NewStatisticsService создаёт сервис статистики с внедрённым репозиторием.
func NewStatisticsService(
	statisticsRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
