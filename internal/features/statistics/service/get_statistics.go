package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
	core_errors "github.com/nilchan-social/golang-todoapp/internal/core/errors"
)

// GetStatistics возвращает агрегированную статистику по задачам.
//
// Параметры фильтрации (все опциональны):
//   - userID — считать статистику только для задач конкретного пользователя
//   - from   — начало временного диапазона (включительно)
//   - to     — конец временного диапазона (не включительно)
//
// Проверяем, что to > from до обращения к БД — бессмысленный диапазон
// невозможно обработать корректно.
func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *uuid.UUID,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"`to` must be after `from`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository: %w", err)
	}

	statistics := domain.CreateStatistics(tasks)

	return statistics, nil
}
