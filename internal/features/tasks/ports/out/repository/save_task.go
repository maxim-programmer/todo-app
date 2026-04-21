package ports_out_repository

import "github.com/nilchan-social/golang-todoapp/internal/core/domain"

// SaveTaskParams / SaveTaskResult — параметры и результат исходящего порта для сохранения задачи.
//
// Ядро (сервис) передаёт адаптеру готовый доменный объект domain.Task.
// Адаптер репозитория (Postgres/Cached) сам решает, как его сохранить: INSERT в Postgres, PUT в DynamoDB, Кеш в Redis и т.д.
// Ядро (уровень сервиса) ничего не знает о схеме хранения — это намеренная изоляция.
type SaveTaskParams struct {
	Task domain.Task
}

func NewSaveTaskParams(
	task domain.Task,
) SaveTaskParams {
	return SaveTaskParams{
		Task: task,
	}
}

type SaveTaskResult struct {
	Task domain.Task
}

func NewSaveTaskResult(
	task domain.Task,
) SaveTaskResult {
	return SaveTaskResult{
		Task: task,
	}
}
