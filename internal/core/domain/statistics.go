package domain

import (
	"time"
)

// Statistics — доменная модель агрегированной статистики по задачам.
// Поля-указатели (*float64, *time.Duration) означают, что значение может
// отсутствовать: например, если нет ни одной завершённой задачи,
// среднее время выполнения не имеет смысла и будет nil.
type Statistics struct {
	TasksCreated               int
	TasksCompleted             int
	TasksCompletedRate         *float64       // процент выполненных задач (nil если задач нет)
	TasksAverageCompletionTime *time.Duration // среднее время выполнения (nil если нет выполненных)
}

// CreateStatistics вычисляет статистику из набора задач.
func CreateStatistics(tasks []Task) Statistics {
	if len(tasks) == 0 {
		return NewStatistics(0, 0, nil, nil)
	}

	tasksCreated := len(tasks)

	tasksCompleted := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		// CompletionDuration() возвращает nil для незавершённых задач,
		// поэтому суммируем только реальные значения.
		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	// Среднее время вычисляем только если есть хотя бы одна завершённая задача.
	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)

		tasksAverageCompletionTime = &avg
	}

	return NewStatistics(
		tasksCreated,
		tasksCompleted,
		&tasksCompletedRate,
		tasksAverageCompletionTime,
	)
}

// NewStatistics — конструктор для создания объекта Statistics из готовых данных
func NewStatistics(
	tasksCreated int,
	tasksCompleted int,
	tasksCompletedRate *float64,
	tasksAverageCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:               tasksCreated,
		TasksCompleted:             tasksCompleted,
		TasksCompletedRate:         tasksCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}
