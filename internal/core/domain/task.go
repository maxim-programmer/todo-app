package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/nilchan-social/golang-todoapp/internal/core/errors"
)

// Task — доменная сущность задачи. Представляет задачу в системе.
//
// Version — счётчик для оптимистичной блокировки (optimistic locking):
// при обновлении БД проверяет, что версия не изменилась с момента чтения.
// Если версия изменилась (другой запрос успел сохранить запись) — вернётся ErrConflict.
//
// Description и CompletedAt — nil означает отсутствие значения (NULL в базе данных).
type Task struct {
	ID      uuid.UUID
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserID uuid.UUID
}

// NewTask — конструктор для восстановления задачи по имеющемуся набору данных
func NewTask(
	id uuid.UUID,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserID uuid.UUID,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: authorUserID,
	}
}

// CreateTask создаёт новую задачу с автоматически сгенерированными полями:
// - ID: уникальный UUID v4
// - Version: 1 (начальная версия для optimistic locking)
// - Completed: false (новая задача всегда не выполнена)
// - CreatedAt: текущее время
// - CompletedAt: nil (не выполнена → нет времени выполнения)
//
// Используется при создании задачи в системе.
func CreateTask(
	title string,
	description *string,
	authorUserID uuid.UUID,
) Task {
	var (
		id                     = uuid.New()
		version                = 1
		completed              = false
		createdAt              = time.Now()
		completedAt *time.Time = nil
	)

	return NewTask(
		id,
		version,
		title,
		description,
		completed,
		createdAt,
		completedAt,
		authorUserID,
	)
}

// CompletionDuration возвращает время, затраченное на выполнение задачи.
// Возвращает nil если задача ещё не выполнена.
// Используется при подсчёте статистики.
func (t *Task) CompletionDuration() *time.Duration {
	if !t.Completed {
		return nil
	}

	if t.CompletedAt == nil {
		return nil
	}

	duration := t.CompletedAt.Sub(t.CreatedAt)

	return &duration
}

// Validate проверяет инварианты доменной модели — правила, которые всегда должны
// выполняться для корректной задачи. Инварианты зеркалируют CHECK-ограничения в БД.
//
// Используем len([]rune(...)) вместо len(string) для корректного подсчёта
// символов в Unicode-строках (кириллица занимает 2 байта, но 1 руну).
func (t *Task) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf(
			"invalid `Title` len: %d: %w",
			titleLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf(
				"invalid `Description` len: %d: %w",
				descriptionLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	// Инвариант: Completed и CompletedAt должны быть согласованы.
	// Если задача выполнена, CompletedAt обязателен и не может быть раньше CreatedAt.
	// Если задача не выполнена, CompletedAt должен быть nil.
	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"`CompletedAt` can't be `nil` if `Completed`==`true`: %w",
				core_errors.ErrInvalidArgument,
			)
		}

		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"`CompletedAt` can't be before `CreatedAt`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"`CompletedAt` must be `nil` if `Completed`==`false`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

// TaskPatch содержит изменения для частичного обновления задачи (PATCH).
// Каждое поле обёрнуто в Nullable, чтобы различать «не передано» и «передано null».
// Подробнее о Nullable: см. internal/core/domain/nullable.go.
type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

// NewTaskPatch — конструктор TaskPatch.
func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
	completed Nullable[bool],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
}

// Validate проверяет корректность патча до его применения.
// Title и Completed не могут быть явно выставлены в null — это обязательные поля.
func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf(
			"`Title` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Completed.Set && p.Completed.Value == nil {
		return fmt.Errorf(
			"`Completed` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

// ApplyPatch применяет изменения из патча к задаче (патчит задачу).
//
// Важная деталь реализации: изменения применяются к копии (tmp := *t),
// а не к самому объекту. Это позволяет проверить результат через Validate()
// до финальной замены — транзакционно: либо всё применяется, либо ничего.
//
// При смене Completed → true автоматически устанавливается CompletedAt = now().
// При смене Completed → false CompletedAt сбрасывается в nil.
func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}

	// Работаем с копией, чтобы не мутировать оригинал до успешной валидации.
	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.Completed.Set {
		tmp.Completed = *patch.Completed.Value

		if tmp.Completed {
			completedAt := time.Now()
			tmp.CompletedAt = &completedAt
		} else {
			tmp.CompletedAt = nil
		}
	}

	// Проверяем, что итоговое состояние задачи валидно.
	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}

	// Только после успешной валидации заменяем оригинал.
	*t = tmp

	return nil
}
