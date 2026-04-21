package tasks_service

import (
	"context"
	"fmt"

	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
	ports_in "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/in"
	ports_out_repository "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/out/repository"
)

// CreateTask — use case создания задачи. Это бизнес-логика ядра. Знакомый нам уровень сервиса.
//
// Поток данных через границы гексагона:
//
//	HTTP-адаптер: CreateTaskRequest → [ports_in.CreateTaskParams] → сервис
//	                                                                   ↓
//	                                               domain.Task (чистый доменный объект)
//	                                                                   ↓
//	                                                                сервис → [ports_out.SaveTaskParams] → репозиторий-адаптер → БД
//
// Квадратные скобки [] обозначают пересечение границы порта.
// Сервис оперирует только domain.Task и типами портов — никакого HTTP и SQL.
//
// Шаги:
//  1. domain.CreateTask — конструирует доменный объект (UUID, version=0, completed=false).
//  2. task.Validate — проверяет инварианты домена до обращения к хранилищу.
//  3. SaveTask — передаёт объект через исходящий порт в адаптер хранилища.
func (s *TasksService) CreateTask(
	ctx context.Context,
	params ports_in.CreateTaskParams,
) (ports_in.CreateTaskResult, error) {
	task := domain.CreateTask(
		params.Title,
		params.Description,
		params.AuthorUserID,
	)

	if err := task.Validate(); err != nil {
		return ports_in.CreateTaskResult{}, fmt.Errorf("validate task domain: %w", err)
	}

	repoParams := ports_out_repository.NewSaveTaskParams(task)
	repoResult, err := s.tasksRepository.SaveTask(ctx, repoParams)
	if err != nil {
		return ports_in.CreateTaskResult{}, fmt.Errorf("save task in repository: %w", err)
	}

	return ports_in.NewCreateTaskResult(
		repoResult.Task,
	), nil
}
