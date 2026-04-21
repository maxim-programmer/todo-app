package adapters_in_transport_http

import (
	"net/http"

	core_logger "github.com/nilchan-social/golang-todoapp/internal/core/logger"
	core_http_request "github.com/nilchan-social/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/nilchan-social/golang-todoapp/internal/core/transport/http/response"
	ports_in "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/in"
)

// DeleteTask    godoc
// @Summary      Удаление задачи
// @Description  Удалить существующую в системе задачу по её ID
// @Tags         tasks
// @Param        id   path  string  true "ID удаляемой задачи" Format(uuid)
// @Success      204  "Успешное удаление задачи"
// @Failure      400  {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      404  {object}  core_http_response.ErrorResponse "Task not found"
// @Failure      500  {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	serviceParams := ports_in.NewDeleteTaskParams(taskID)
	if _, err := h.tasksService.DeleteTask(ctx, serviceParams); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)

		return
	}

	responseHandler.NoContentResponse()
}
