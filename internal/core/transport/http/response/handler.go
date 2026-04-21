package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
	core_errors "github.com/nilchan-social/golang-todoapp/internal/core/errors"
	core_logger "github.com/nilchan-social/golang-todoapp/internal/core/logger"
	"go.uber.org/zap"
)

// HTTPResponseHandler инкапсулирует логику записи HTTP-ответов.
// Хранит логгер и ResponseWriter, чтобы обработчикам не нужно было
// передавать их каждый раз явно.
type HTTPResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

// NewHTTPResponseHandler создаёт обработчик ответов для конкретного запроса.
// Вызывается в начале каждого HTTP-обработчика.
func NewHTTPResponseHandler(
	log *core_logger.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

// JSONResponse сериализует responseBody в JSON и записывает в ответ с указанным статус-кодом.
// Content-Type автоматически определяется json.NewEncoder.
func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

// NoContentResponse отправляет HTTP 204 No Content — используется при успешном DELETE.
func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

// HTMLResponse отправляет HTML-страницу с Content-Type: text/html.
func (h *HTTPResponseHandler) HTMLResponse(htmlFile domain.File) {
	h.rw.WriteHeader(http.StatusOK)

	h.rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := h.rw.Write(htmlFile.Buffer()); err != nil {
		h.log.Error("write HTML HTTP response", zap.Error(err))
	}
}

// ErrorResponse транслирует core ошибку в HTTP-статус через errors.Is().
//
// Маппинг:
//   - ErrInvalidArgument → 400
//   - ErrNotFound        → 404
//   - ErrConflict        → 409
//   - остальное          → 500
//
// Каждый тип ошибки логируется на соответствующем уровне (Warn/Debug/Error).
func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		err,
		msg,
	)
}

// PanicResponse формирует HTTP 500 при перехвате паники.
// Вызывается из middleware Panic — см. internal/core/transport/http/middleware/common.go.
func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		err,
		msg,
	)
}

// errorResponse — внутренний метод: собирает ErrorResponse и вызывает JSONResponse.
func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	response := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}

	h.JSONResponse(
		response,
		statusCode,
	)
}
