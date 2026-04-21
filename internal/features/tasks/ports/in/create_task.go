package ports_in

import (
	"github.com/google/uuid"
	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// CreateTaskParams — параметры входящего порта TasksService.CreateTask.
//
// Этот тип принадлежит ЯДРУ (уровню сервиса), а не HTTP-адаптеру.
// HTTP-адаптер переводит свои типы (PatchTaskRequest) в этот тип,
// чтобы «заговорить» на языке ядра.
// Если появится gRPC-адаптер/kafka — они будут делать то же самое независимо.
type CreateTaskParams struct {
	Title        string
	Description  *string
	AuthorUserID uuid.UUID
}

func NewCreateTaskParams(
	title string,
	description *string,
	authorUserID uuid.UUID,
) CreateTaskParams {
	return CreateTaskParams{
		Title:        title,
		Description:  description,
		AuthorUserID: authorUserID,
	}
}

// CreateTaskResult — результирующее значение входящего порта TasksService.CreateTask.
// Возвращает доменный объект domain.Task — универсальный язык ядра,
// который каждый адаптер переведёт в свой формат (JSON, protobuf, []byte и т.д.).
type CreateTaskResult struct {
	Task domain.Task
}

func NewCreateTaskResult(
	task domain.Task,
) CreateTaskResult {
	return CreateTaskResult{
		Task: task,
	}
}
