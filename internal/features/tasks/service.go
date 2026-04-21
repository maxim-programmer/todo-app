package tasks_service

import (
	repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// TasksService — ядро гексагона фичи tasks. Уровень сервиса.
//
// В гексагональной архитектуре сервис является центральным элементом:
// он не зависит ни от одного адаптера (adapters: HTTP, Postgres, Redis и т.д.) —
// только от интерфейсов портов (ports), которые сам же и определяет.
//
// Зависимости TasksService:
//
//	[входящий адаптер] → ports_in.TasksService (определяет ядро)
//	                                ↓ реализует
//	                              TasksService (это ядро)
//	                                ↓ использует
//	                     ports_out.TasksRepository (определяет ядро)
//	                                ↓ реализует
//	[исходящий адаптер: Postgres, Cached...]
//
// Благодаря этому сервис можно тестировать в изоляции, подменив
// оба адаптера заглушками — без поднятия HTTP-сервера или базы данных.
type TasksService struct {
	tasksRepository repository.TasksRepository
}

func NewTasksService(
	tasksRepository repository.TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
