package adapters_in_transport_http

import (
	"net/http"

	core_http_server "github.com/nilchan-social/golang-todoapp/internal/core/transport/http/server"
	ports_in "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/in"
)

// TasksHTTPHandler — входящий адаптер (adapters/in) для адаптирования HTTP транспорта под входящий порт (ports/in) уровня сервиса.
//
// В гексагональной архитектуре входящий адаптер — это «левая сторона» гексагона.
// Его единственная ответственность: принять внешний вызов (HTTP-запрос) и
// перевести его в вызов входящего порта.
//
// Обязанности адаптера:
//   - декодировать и валидировать HTTP-запрос
//   - перевести HTTP-типы в DTO порта (например, PatchTaskRequest → ports_in.PatchTaskParams)
//   - вызвать метод tasksService (через интерфейс порта, не напрямую)
//   - перевести результат / ошибку обратно в HTTP-ответ
//
// Чего адаптер НЕ делает:
//   - не содержит бизнес-логики
//   - не обращается к хранилищу напрямую
//   - не знает о существовании Postgres, Redis или других адаптеров
//
// Зависит от ports_in.TasksService (интерфейс), а не от конкретного tasks_service.TasksService.
type TasksHTTPHandler struct {
	tasksService ports_in.TasksService
}

func NewTasksHTTPHandler(
	tasksService ports_in.TasksService,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: h.DeleteTask,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{id}",
			Handler: h.PatchTask,
		},
	}
}
