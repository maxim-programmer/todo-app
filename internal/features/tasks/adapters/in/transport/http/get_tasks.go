package adapters_in_transport_http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	core_logger "github.com/nilchan-social/golang-todoapp/internal/core/logger"
	core_http_request "github.com/nilchan-social/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/nilchan-social/golang-todoapp/internal/core/transport/http/response"
	ports_in "github.com/nilchan-social/golang-todoapp/internal/features/tasks/ports/in"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks      godoc
// @Summary      Список задач
// @Description  Просмотр списка задач с опциональной пагинацией и/или фильтрацией по ID автора задачи
// @Tags         tasks
// @Produce      json
// @Param        user_id  query     string  false  "Фильтрация задач по ID автора" Format(uuid)
// @Param        limit    query     int  false  "Размер страницы с задачами"
// @Param        offset   query     int  false  "Смещение страницы с задачами"
// @Success      200      {object}  GetTasksResponse                 "Список задач"
// @Failure      400      {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      500      {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/limit/offset query params",
		)

		return
	}

	serviceParams := ports_in.NewGetTasksParams(
		userID,
		limit,
		offset,
	)
	serviceResult, err := h.tasksService.GetTasks(ctx, serviceParams)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)

		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(serviceResult.Tasks))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*uuid.UUID, *int, *int, error) {
	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetUUIDQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil
}
