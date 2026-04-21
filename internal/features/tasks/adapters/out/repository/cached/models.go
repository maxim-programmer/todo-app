package adapters_out_repository_cached

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// TaskModel — Redis-представление задачи для JSON-сериализации.
// Хранится по ключу "task:<uuid>" с TTL, заданным пулом.
type TaskModel struct {
	ID           uuid.UUID  `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserID uuid.UUID  `json:"author_user_id"`
}

func domainToModel(task domain.Task) TaskModel {
	return TaskModel{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func modelToDomain(model TaskModel) domain.Task {
	return domain.NewTask(
		model.ID,
		model.Version,
		model.Title,
		model.Description,
		model.Completed,
		model.CreatedAt,
		model.CompletedAt,
		model.AuthorUserID,
	)
}

// taskKey возвращает Redis-ключ для одиночной задачи.
// Формат: "task:<uuid>"
func taskKey(id uuid.UUID) string {
	return fmt.Sprintf("task:%s", id)
}

func (m *TaskModel) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("serialize task: %w", err)
	}

	return bytes, nil
}

func (m *TaskModel) Deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, m); err != nil {
		return fmt.Errorf("deserialize task: %w", err)
	}

	return nil
}

// TaskListModel — Redis-представление списка задач для JSON-сериализации.
// Хранится в Redis hash по ключу "tasks:<userID>" или "tasks:all".
// Поле (field) хеша задаётся функцией tasksListField и кодирует пагинацию.
type TaskListModel []TaskModel

func domainsToModel(tasks []domain.Task) TaskListModel {
	tasksModels := make([]TaskModel, len(tasks))

	for i, task := range tasks {
		tasksModels[i] = domainToModel(task)
	}

	return tasksModels
}

func modelToDomains(list TaskListModel) []domain.Task {
	tasks := make([]domain.Task, len(list))

	for i, model := range list {
		tasks[i] = modelToDomain(model)
	}

	return tasks
}

// tasksListKey возвращает Redis-ключ для hash'а со списками задач.
// Формат: "tasks:<uuid>" — задачи конкретного пользователя,
//
//	"tasks:all"    — задачи всех пользователей (userID == nil).
func tasksListKey(userID *uuid.UUID) string {
	if userID == nil {
		return "tasks:all"
	}

	return fmt.Sprintf("tasks:%s", *userID)
}

// tasksListField возвращает имя поля в Redis hash для конкретного варианта пагинации.
// Формат: "<limit>:<offset>", например "10:0", "nil:nil", "20:nil".
func tasksListField(limit *int, offset *int) string {
	ptrStr := func(v *int) string {
		if v == nil {
			return "nil"
		}

		return strconv.Itoa(*v)
	}

	return fmt.Sprintf("%s:%s", ptrStr(limit), ptrStr(offset))
}

func (m *TaskListModel) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("serialize task list: %w", err)
	}

	return bytes, err
}

func (m *TaskListModel) Deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, m); err != nil {
		return fmt.Errorf("deserialize task list: %w", err)
	}

	return nil
}
